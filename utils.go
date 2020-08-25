package main

import (
	"fmt"
	"github.com/jhole89/discovery-backend/connectors"
	"github.com/jhole89/discovery-backend/database"
	"strings"
)

func nodeToGraph(graph database.Graph, node *connectors.Node) database.Entity {
	entityA := database.Entity{Name: node.Name, Context: node.Context}
	graph.CreateEntity(entityA)
	if node.Children != nil {
		for _, childNode := range node.Children {
			entityB := nodeToGraph(graph, childNode)
			relationship := database.Relationship{From: entityA, To: entityB, Context: fmt.Sprintf("has_%s", strings.ToLower(entityB.Context))}
			graph.CreateRelationship(relationship)
		}
	}
	return entityA
}
