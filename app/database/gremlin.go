package database

import (
	"encoding/json"
	"fmt"
	"github.com/schwartzmx/gremtune"
	"log"
	"net/http"
	"time"
)

var err error

type Gremlin struct {
	Connection gremtune.Client
}

func (n *Gremlin) Connect(address string) {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatalf("Lost connection to the database: %s\n", err.Error())
	}(errs)

	dialer := gremtune.NewDialer(address)

	retryCount := 10
	for {
		log.Println("Attempting to connect to Gremlin server at: " + address)
		n.Connection, err = gremtune.Dial(dialer, errs)
		if err != nil {
			if retryCount == 0 {
				log.Fatalln("Unable to connect to Gremlin server: " + err.Error())
				return
			}

			log.Printf("Could not connect to Gremlin server. Wait 2 seconds. %d retries left...\n", retryCount)
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	n.clean()
}

func (n *Gremlin) clean() string {
	return n.Query("g.V().drop().iterate()")
}

func (n *Gremlin) CreateEntity(e Entity) string {
	queryString := fmt.Sprintf("g.addV('%s').property('name', '%s')", e.Context, e.Name)

	for _, property := range e.Properties {
		queryString += fmt.Sprintf(".property('%s', '%s')", property.Attribute, property.Value)
	}
	return n.Query(queryString)
}

func (n *Gremlin) CreateRelationship(r Relationship) string {
	queryString := fmt.Sprintf("g.addE('%s').from(g.V().has('%s', 'name', '%s')).to(g.V().has('%s', 'name', '%s'))", r.Context, r.From.Context, r.From.Name, r.To.Context, r.To.Name)

	return n.Query(queryString)
}

func (n *Gremlin) Query(queryString string) string {
	resp, err := n.runQuery(queryString)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", unmarshall(resp))
}

func (n *Gremlin) runQuery(queryString string) ([]gremtune.Response, error) {
	resp, err := n.Connection.Execute(queryString)
	if err != nil {
		log.Printf("Unable to execute query: %s. Err: %s\n", queryString, err.Error())
		return nil, err
	}
	return resp, nil
}

func unmarshall(resp []gremtune.Response) []byte {
	j, err := json.Marshal(resp[0].Result.Data)

	if err != nil {
		log.Printf("Unable to unpack result: %s\n", err.Error())
		return []byte{}
	}

	return j
}

func (n *Gremlin) Read(w http.ResponseWriter) []byte {
	return n.httpQuery(w, "g.V().elementMap()")
}

func (n *Gremlin) httpQuery(w http.ResponseWriter, queryString string) []byte {
	resp, err := n.runQuery(queryString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return unmarshall(resp)
}
