package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jhole89/orbital/ent"
	"log"
	"strings"
	"time"
)

// newGraph establishes a new connection to a supported GraphDB passed by string name
func newGraph(graphName string, dsn string) (*ent.Client, error) {

	retryCount := 10
	for {
		log.Println("Attempting to connect to server at: " + dsn)

		client, err := ent.Open(strings.ToLower(graphName), dsn)

		if err != nil {
			if retryCount == 0 {
				log.Println("Unable to connect to server: " + err.Error())
				return nil, err
			}

			log.Printf("Could not connect to server. Waiting 2 seconds. %d retries left...\n", retryCount)
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Connected to server at: " + dsn)
			defer client.Close()
			return client, nil
		}
	}
}

func deleteAll(ctx context.Context, graph *ent.Client) error {
	numData, err := graph.Data.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	log.Printf("#Data Vertices deleted: %d\n", numData)

	return nil
}

func createDataVertex(ctx context.Context, client *ent.Client, name, context string) (*ent.Data, error) {
	d, err := client.Data.
		Create().
		SetName(name).
		SetContext(context).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed creating Data Vertex: %v", err)
	}
	log.Printf("Data Vertex was created: %s\n", d.String())

	return d, nil
}

func createRelationship(ctx context.Context, from, to *ent.Data) (*ent.Data, error) {
	switch to.Context {
	case "table":
		return from.Update().AddHasTable(to).Save(ctx)
	case "field":
		return from.Update().AddHasField(to).Save(ctx)
	default:
		return nil, errors.New(fmt.Sprintf("%s is not a valid context", to.Context))
	}
}

func listDataEntities(ctx context.Context, graph *ent.Client) ([]*ent.Data, error) {
	return graph.Data.Query().All(ctx)
}

func getDataEntity(ctx context.Context, graph *ent.Client, id int) (*ent.Data, error) {
	return graph.Data.Get(ctx, id)
}

type Property struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}
