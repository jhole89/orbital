package main

import (
	"encoding/json"
	"fmt"
	"github.com/jhole89/orbital/database"
	"log"
	"net/http"
)

var (
	conf  Config
	graph database.Graph
	err error
)

func main() {
	conf.getConf()

	graph, err = database.GetGraph(conf.Database.Type, conf.Database.Endpoint)
	if err != nil {
		log.Println(err)
	}

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if err = json.NewEncoder(w).Encode(graphqlQuery(r.URL.Query().Get("query"), dataSchema)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		if err = json.NewEncoder(w).Encode(graphqlQuery(r.URL.Query().Get("task"), adminSchema)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})

	log.Printf("Launching Orbital API at http://127.0.0.1:%d\n", conf.Service.Port)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Service.Port), nil); err != nil {
		log.Fatalln("Unable to serve.")
	}
}
