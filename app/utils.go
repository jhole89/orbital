package main

import (
	"fmt"
	"github.com/jhole89/orbital/connectors"
	"github.com/jhole89/orbital/database"
	"strings"
)

func nodeToGraph(graph database.Graph, node *connectors.Node) (database.Entity, error) {
	entityA := database.Entity{Name: node.Name, Context: node.Context}
	_, err := graph.CreateEntity(entityA)
	if err != nil {
		return database.Entity{}, err
	}
	if node.Children != nil {
		for _, childNode := range node.Children {
			entityB, err := nodeToGraph(graph, childNode)
			if err != nil {
				return entityA, nil
			}
			relationship := database.Relationship{From: entityA, To: entityB, Context: fmt.Sprintf("has_%s", strings.ToLower(entityB.Context))}
			_, err = graph.CreateRelationship(relationship)
			if err != nil {
				return entityB, err
			}
		}
	}
	return entityA, nil
}
