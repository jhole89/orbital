package main

import (
	"fmt"
	"github.com/jhole89/orbital/connectors"
	"github.com/jhole89/orbital/database"
	"strings"
)

func loadGraph(lake *LakeConfig) error {
	driver := connectors.GetDriver(fmt.Sprintf("%s%s", lake.Provider, lake.Store), lake.Address)

	dbTopology, _ := driver.Index()

	for _, node := range dbTopology {
		_, err := nodeToGraph(graph, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func nodeToGraph(graph *database.Graph, node *connectors.Node) (*database.Entity, error) {
	var propertyList []database.Property
	for k, v := range node.Properties {
		propertyList = append(propertyList, database.Property{Attribute: k, Value: v})
	}
	entityA := database.Entity{Name: node.Name, Context: node.Context, Properties: propertyList}

	_, err := (*graph).CreateEntity(entityA)
	if err != nil {
		return &database.Entity{}, err
	}
	if node.Children != nil {
		for _, childNode := range node.Children {
			entityB, err := nodeToGraph(graph, childNode)
			if err != nil {
				return &entityA, nil
			}
			relationship := database.Relationship{From: &entityA, To: entityB, Context: fmt.Sprintf("has_%s", strings.ToLower(entityB.Context))}
			_, err = (*graph).CreateRelationship(relationship)
			if err != nil {
				return entityB, err
			}
		}
	}
	return &entityA, nil
}
