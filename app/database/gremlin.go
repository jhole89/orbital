package database

import (
	"encoding/json"
	"fmt"
	"github.com/schwartzmx/gremtune"
	"log"
	"net/http"
	"time"
)

type Gremlin struct {
	Client *gremtune.Client
}

func NewGremlin(dsn string) (Graph, error) {

	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatalf("Lost connection to the database: %s\n", err.Error())
	}(errs)

	dialer := gremtune.NewDialer(dsn)

	retryCount := 10
	for {
		log.Println("Attempting to connect to Gremlin server at: " + dsn)
		conn, err := gremtune.Dial(dialer, errs)
		if err != nil {
			if retryCount == 0 {
				log.Println("Unable to connect to Gremlin server: " + err.Error())
				return nil, err
			}

			log.Printf("Could not connect to Gremlin server. Wait 2 seconds. %d retries left...\n", retryCount)
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			return &Gremlin{Client: &conn}, nil
		}
	}
}

func (g *Gremlin) Clean() (string, error) {
	return g.Query("g.V().drop().iterate()")
}

func (g *Gremlin) CreateEntity(e Entity) (string, error) {
	queryString := fmt.Sprintf("g.addV('%s').property('name', '%s')", e.Context, e.Name)

	for _, property := range e.Properties {
		queryString += fmt.Sprintf(".property('%s', '%s')", property.Attribute, property.Value)
	}
	return g.Query(queryString)
}

func (g *Gremlin) CreateRelationship(r Relationship) (string, error) {
	queryString := fmt.Sprintf("g.addE('%s').from(g.V().has('%s', 'name', '%s')).to(g.V().has('%s', 'name', '%s'))", r.Context, r.From.Context, r.From.Name, r.To.Context, r.To.Name)

	return g.Query(queryString)
}

func (g *Gremlin) Query(queryString string) (string, error) {
	resp, err := g.runQuery(queryString)
	if err != nil {
		return "", err
	}

	s, e := unmarshall(resp)
	if e != nil {
		return "", e
	}
	return fmt.Sprintf("%s", s), nil
}

func (g *Gremlin) runQuery(queryString string) ([]gremtune.Response, error) {
	resp, err := g.Client.Execute(queryString)
	if err != nil {
		log.Printf("Unable to execute query: %s. Err: %s\n", queryString, err.Error())
		return nil, err
	}
	return resp, nil
}

func unmarshall(resp []gremtune.Response) ([]byte, error) {
	j, err := json.Marshal(resp[0].Result.Data)

	if err != nil {
		log.Printf("Unable to unpack result: %s\n", err.Error())
		return nil, err
	}
	return j, nil
}

func (g *Gremlin) Read(w http.ResponseWriter) ([]byte, error) {
	return g.httpQuery(w, "g.V().elementMap()")
}

func (g *Gremlin) httpQuery(w http.ResponseWriter, queryString string) ([]byte, error) {
	resp, err := g.runQuery(queryString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return unmarshall(resp)
}
