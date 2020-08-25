package main

import (
	"fmt"
	"github.com/jhole89/discovery-backend/connectors"
	"github.com/jhole89/discovery-backend/database"
	"log"
	"net/http"
)

var (
	graph database.Graph
)

func main() {

	var conf Config
	conf.getConf()

	graph = &database.AwsNeptuneDB{Address: conf.Database.Endpoint}
	graph.Connect()

	for _, lake := range conf.Lakes {
		driver := &connectors.AwsAthenaConnector{Address: lake.Address}
		conn := driver.Connect()
		dbTopology := driver.Index(conn)
		for _, node := range dbTopology {
			nodeToGraph(graph, node)
		}
		resp := graph.Query("g.V().elementMap()")
		fmt.Printf("Entities: %s\n", resp)

		resp = graph.Query("g.E().elementMap()")
		fmt.Printf("Relationships: %s\n", resp)
	}

	registerRoutes()
	fmt.Printf("Server running at http://127.0.0.1:%d\n", conf.Service.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Service.Port), nil)
	if err != nil {
		log.Fatal("Unable to serve")
	}
}
