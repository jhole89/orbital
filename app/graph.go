package main

import (
	"fmt"
	"github.com/jhole89/orbital/connectors"
	"github.com/jhole89/orbital/database"
)

func reIndex(graph database.Graph, lakes []*LakeConfig) error {
	if err := graph.Clean(); err != nil {
		return err
	}
	if err := index(graph, lakes); err != nil {
		return err
	}
	return nil
}

func index(graph database.Graph, lakes []*LakeConfig) error {
	for _, lake := range lakes {
		if err := indexLake(graph, lake); err != nil {
			return err
		}
	}
	return nil
}

func indexLake(graph database.Graph, lake *LakeConfig) error {
	driver := connectors.GetDriver(fmt.Sprintf("%s%s", lake.Provider, lake.Store), lake.Address)
	dbTopology, err := driver.Index()
	if err != nil {
		return err
	}
	if err := loadGraph(graph, dbTopology); err != nil {
		return err
	}
	return nil
}

func loadGraph(graph database.Graph, nodes []*connectors.Node) error {
	for _, n := range nodes {
		if _, err := loadNode(graph, n); err != nil {
			return err
		}
	}
	return nil
}

func loadNode(graph database.Graph, node *connectors.Node) (*database.Entity, error) {
	entityFrom, err := graph.CreateEntity(&database.Entity{Context: node.Context, Name: node.Name})
	if err != nil {
		return nil, err
	}
	if node.Children != nil {
		for _, childNode := range node.Children {
			entityTo, err := loadNode(graph, childNode)
			if err != nil {
				return nil, err
			}
			_, err = graph.CreateRelationship(&database.Relationship{Context: "owns", From: entityFrom, To: entityTo})
			if err != nil {
				return nil, err
			}
		}
	}
	return entityFrom, nil
}
