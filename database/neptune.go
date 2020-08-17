package database

import (
	"encoding/json"
	"fmt"
	"github.com/schwartzmx/gremtune"
	"log"
	"net/http"
)

var (
	db gremtune.Client
	err error
)

func Start(endpoint string)  {
	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Fatal("Lost connection to the database: " + err.Error())
	}(errs) // Example of connection error handling logic

	dialer := gremtune.NewDialer(endpoint) // Returns a WebSocket dialer to connect to Gremlin Server
	db, err = gremtune.Dial(dialer, errs)  // Returns a gremtune client to interact with
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Read(w http.ResponseWriter) []byte {

	drop()

	query(w,"g.addV('data').property('name', 'footable').property('location', 'some-uri')")
	query(w,"g.addV('data').property('name', 'bartable').property('location', 'some-uri2')")
	query(w,"g.addV('field').property('name', 'foofield')")
	query(w, "g.addE('has_field').from(g.V().has('data', 'name', 'footable')).to(g.V().has('field', 'name', 'foofield'))")
	query(w, "g.addE('has_field').from(g.V().has('data', 'name', 'bartable')).to(g.V().has('field', 'name', 'foofield'))")

	return query(w, "g.V().has('field', 'name', 'foofield').in('has_field').values('name')")
}

func query(w http.ResponseWriter, queryString string) []byte {
	res, err := db.Execute(queryString)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	j, err := json.Marshal(res[0].Result.Data) // res will return a list of resultsets,  where the data is a json.RawMessage

	if err != nil {
		fmt.Println("Unable to unpack result")
		panic(err)
	}

	return j
}

func drop() {
	_, err = db.Execute("g.V().drop().iterate()")
	if err != nil {
		fmt.Println("Unable to empty db")
		panic(err)
	}
}

