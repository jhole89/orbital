package database

import (
	"encoding/json"
	"fmt"
	"github.com/schwartzmx/gremtune"
	"log"
	"net/http"
)

var err error

type AwsNeptuneDB struct {
	Connection gremtune.Client
}

func (n *AwsNeptuneDB) Connect(address string)  {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatalf("Lost connection to the database: %s\n", err.Error())
	}(errs)

	dialer := gremtune.NewDialer(address)
	n.Connection, err = gremtune.Dial(dialer, errs)
	if err != nil {
		panic(err)
		return
	}
	n.clean()
}

func (n *AwsNeptuneDB) clean() string {
	return n.Query("g.V().drop().iterate()")
}

func (n *AwsNeptuneDB) CreateEntity(e Entity) string {
	queryString := fmt.Sprintf("g.addV('%s').property('name', '%s')", e.Context, e.Name)

	for _, property := range e.Properties {
		queryString += fmt.Sprintf(".property('%s', '%s')", property.Attribute, property.Value)
	}
	return n.Query(queryString)
}

func (n *AwsNeptuneDB) CreateRelationship(r Relationship) string {
	queryString := fmt.Sprintf("g.addE('%s').from(g.V().has('%s', 'name', '%s')).to(g.V().has('%s', 'name', '%s'))", r.Context, r.From.Context, r.From.Name, r.To.Context, r.To.Name)

	return n.Query(queryString)
}

func (n *AwsNeptuneDB) Query(queryString string) string {
	resp, err := n.runQuery(queryString)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", unmarshall(resp))
}

func (n *AwsNeptuneDB) runQuery(queryString string) ([]gremtune.Response, error) {
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

func (n *AwsNeptuneDB) Read(w http.ResponseWriter) []byte {
	return n.httpQuery(w, "g.V().has('field', 'name', 'foofield').in('has_field').values('name')")
}

func (n *AwsNeptuneDB) httpQuery(w http.ResponseWriter, queryString string) []byte {
	resp, err := n.runQuery(queryString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return unmarshall(resp)
}
