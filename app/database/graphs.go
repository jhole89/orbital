package database

import (
	"log"
	"net/http"
	"strings"
)

type Graph interface {
	Clean() (string, error)
	Query(queryString string) (string, error)
	CreateEntity(e Entity) (string, error)
	CreateRelationship(r Relationship) (string, error)
	Read(w http.ResponseWriter) ([]byte, error)
}

type Entity struct {
	Context    string
	Name       string
	Properties []Property
}

type Property struct {
	Attribute string
	Value     string
}

type Relationship struct {
	From    Entity
	To      Entity
	Context string
}

func GetGraph(graphName string, endpoint string) Graph {

	var supportedGraph = map[string]func(string) (Graph, error){
		"awsneptune": NewGremlin,
		"gremlin":    NewGremlin,
		"tinkerpop":  NewGremlin,
	}

	graphInitialiser, ok := supportedGraph[strings.ToLower(graphName)]

	if ok {
		conn, _ := graphInitialiser(endpoint)
		return conn
	} else {
		keys := make([]string, len(supportedGraph))
		for k := range supportedGraph {
			keys = append(keys, k)
		}
		log.Printf("DB: %s is not supported. Please specifiy a supported DB in your config.yaml.\nValid DB's: %s", graphName, keys)
		return nil
	}
}
