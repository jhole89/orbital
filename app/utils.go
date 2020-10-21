package main

import (
	"fmt"
	"github.com/jhole89/orbital/connectors"
	"github.com/jhole89/orbital/ent"
	"log"
)

func loadGraph(graph *ent.Client, lake *LakeConfig) error {
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

func nodeToGraph(graph *ent.Client, node *connectors.Node) (*ent.Data, error) {
	var propertyList []Property
	for k, v := range node.Properties {
		propertyList = append(propertyList, Property{Attribute: k, Value: v})
	}

	entityFrom, err := createDataVertex(ctx, graph, node.Name, node.Context)
	if err != nil {
		return nil, err
	}
	if node.Children != nil {
		for _, childNode := range node.Children {
			entityTo, err := nodeToGraph(graph, childNode)
			if err != nil {
				return entityFrom, nil
			}
			_, err = createRelationship(ctx, entityFrom, entityTo)
			if err != nil {
				return entityTo, err
			}
			log.Printf("Created relationship between %s and %s\n", entityFrom.String(), entityTo.String())
		}
	}
	return entityFrom, nil
}
