package database

import (
	"fmt"
	"strings"
)

type Graph interface {
	Clean() error
	CreateEntity(e *Entity) (*Entity, error)
	GetEntity(id interface{}) (*Entity, error)
	ListEntities() ([]*Entity, error)
	CreateRelationship(r *Relationship) (*Relationship, error)
	GetRelationships(id interface{}, context string) ([]*Entity, error)
}

type Entity struct {
	ID         interface{} `json:"id"`
	Context    string      `json:"context"`
	Name       string      `json:"name"`
	Properties []*Property `json:"properties"`
}

type Property struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

type Relationship struct {
	ID      interface{} `json:"id"`
	From    *Entity     `json:"from"`
	To      *Entity     `json:"true"`
	Context string      `json:"context"`
}

// GetGraph establishes a new connection to a supported GraphDB passed by string name
func GetGraph(graphName string, endpoint string) (Graph, error) {

	var supportedGraph = map[string]func(string) (Graph, error){
		"awsneptune": newGremlin,
		"gremlin":    newGremlin,
		"tinkerpop":  newGremlin,
	}

	graphInitializer, ok := supportedGraph[strings.ToLower(graphName)]

	if ok {
		conn, err := graphInitializer(endpoint)
		if err != nil {
			return nil, err
		}
		return conn, nil
	} else {
		keys := make([]string, len(supportedGraph))
		for k := range supportedGraph {
			keys = append(keys, k)
		}
		return nil, fmt.Errorf("DB: %s is not supported. Please specifiy a supported DB in your config.yaml.\nValid DB's: %s", graphName, keys)
	}
}
