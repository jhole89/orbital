package main

import (
	"fmt"
	"github.com/jhole89/orbital/database"
	"log"
	"net/http"
)

var (
	graph *database.Graph
)

func main() {

	var conf Config
	conf.getConf()

	graph = database.GetGraph(conf.Database.Type, conf.Database.Endpoint)
	_, err := (*graph).Clean()
	if err != nil {
		fmt.Println(err)
	}

	for _, lake := range conf.Lakes {
		err = loadGraph(&lake)
		if err != nil {
			log.Println(err)
		}
	}

	registerRoutes()

	fmt.Printf("Server running at http://127.0.0.1:%d\n", conf.Service.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Service.Port), nil)
	if err != nil {
		log.Fatal("Unable to serve")
	}
}
