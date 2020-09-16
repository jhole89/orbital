package database

import (
	"log"
	"strings"
)

type Graph interface {
	Clean() ([]byte, error)
	Query(queryString string) ([]byte, error)
	CreateEntity(e Entity) ([]byte, error)
	CreateRelationship(r Relationship) ([]byte, error)
}

type Entity struct {
	ID         int        `json:"id"`
	Context    string     `json:"context"`
	Name       string     `json:"name"`
	Properties []Property `json:"properties"`
}

type Property struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

type Relationship struct {
	From    *Entity `json:"from"`
	To      *Entity `json:"true"`
	Context string  `json:"context"`
}

// GetGraph establishes a new connection to a supported GraphDB passed by string name
func GetGraph(graphName string, endpoint string) *Graph {

	var supportedGraph = map[string]func(string) (Graph, error){
		"awsneptune": NewGremlin,
		"gremlin":    NewGremlin,
		"tinkerpop":  NewGremlin,
	}

	graphInitialiser, ok := supportedGraph[strings.ToLower(graphName)]

	if ok {
		conn, _ := graphInitialiser(endpoint)
		return &conn
	} else {
		keys := make([]string, len(supportedGraph))
		for k := range supportedGraph {
			keys = append(keys, k)
		}
		log.Printf("DB: %s is not supported. Please specifiy a supported DB in your config.yaml.\nValid DB's: %s", graphName, keys)
		return nil
	}
}
