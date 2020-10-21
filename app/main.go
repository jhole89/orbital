package main

import (
	"context"
	"fmt"
	"github.com/jhole89/orbital/ent"
	"log"
	"net/http"
)

var (
	graph *ent.Client
	ctx   context.Context
)

func main() {
	ctx = context.Background()

	var conf Config
	conf.getConf()

	var err error
	graph, err = newGraph(conf.Database.Type, conf.Database.Endpoint)
	if err != nil {
		log.Println(err)
	}

	if err = deleteAll(ctx, graph); err != nil {
		log.Println(err)
	}

	for _, lake := range conf.Lakes {
		if err := loadGraph(graph, lake); err != nil {
			log.Println(err)
		}
	}

	registerRoutes()

	res, err := listDataEntities(ctx, graph)
	if err != nil {
		log.Println(err)
	}
	for _, d := range res {
		fmt.Printf("Ent: %s\n", d.String())
	}

	fmt.Printf("Starting server at http://127.0.0.1:%d\n", conf.Service.Port)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Service.Port), nil); err != nil {
		log.Fatalln("Unable to serve.")
	}

	if err = graph.Close(); err != nil {
		log.Fatalln("Could not close connection.")
	}
}
