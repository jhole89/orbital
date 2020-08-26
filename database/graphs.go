package database

import (
	"log"
	"net/http"
	"strings"
)

type Graph interface {
	Connect()
	Query(queryString string) string
	CreateEntity(e Entity) string
	CreateRelationship(r Relationship) string
	Read(w http.ResponseWriter) []byte
}

type Entity struct {
	Context    string
	Name       string
	Properties []Property
}

type Property struct {
	Attribute string
	Value string
}

type Relationship struct {
	From    Entity
	To      Entity
	Context string
}

func GetGraph(graphName string, endpoint string) Graph {

	var supportedGraph = map[string]Graph {
		"awsneptune": &AwsNeptuneDB{Address: endpoint},
		"gremlin": &AwsNeptuneDB{Address: endpoint},
		"tinkerpop": &AwsNeptuneDB{Address: endpoint},
	}

	g, ok := supportedGraph[strings.ToLower(graphName)]

	if ok {
		return g
	} else {
		keys := make([]string, len(supportedGraph))
		for k := range supportedGraph {
			keys = append(keys, k)
		}
		log.Fatalf("DB: %s is not supported. Please specifiy a supported DB in your config.yaml.\nValid DB's: %s", graphName, keys)
		return nil
	}
}
