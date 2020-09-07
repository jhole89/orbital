package main

import (
	"fmt"
	"github.com/jhole89/orbital/connectors"
	"github.com/jhole89/orbital/database"
	"log"
	"net/http"
)

var (
	graph database.Graph
)

func main() {

	var conf Config
	conf.getConf()

	graph := database.GetGraph(conf.Database.Type, conf.Database.Endpoint)
	_, err := graph.Clean()

	for _, lake := range conf.Lakes {
		driver := connectors.GetDriver(fmt.Sprintf("%s%s", lake.Provider, lake.Store), lake.Address)

		dbTopology, _ := driver.Index()

		for _, node := range dbTopology {
			nodeToGraph(graph, node)
		}
		resp, _ := graph.Query("g.V().elementMap()")
		fmt.Printf("Entities: %s\n", resp)

		resp, _ = graph.Query("g.E().elementMap()")
		fmt.Printf("Relationships: %s\n", resp)
	}

	registerRoutes()
	fmt.Printf("Server running at http://127.0.0.1:%d\n", conf.Service.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Service.Port), nil)
	if err != nil {
		log.Fatal("Unable to serve")
	}
}
