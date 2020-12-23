package main

import (
	"context"
	"fmt"
	"github.com/jhole89/orbital/connectors"
	"github.com/jhole89/orbital/ent"
	"github.com/jhole89/orbital/ent/data"
	"log"
	"strings"
	"time"
)

type Graph struct {
	conn *ent.Client
}


// newGraph establishes a new connection to a supported GraphDB passed by string name
func newGraph(graphName string, dsn string) (*Graph, error) {

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
			return &Graph{conn: client}, nil
		}
	}
}

func (g *Graph) deleteAll(ctx context.Context) error {
	numData, err := g.conn.Data.Delete().Exec(ctx)
	if err != nil {
		return err
	}
	log.Printf("#Data Vertices deleted: %d\n", numData)

	return nil
}

func (g *Graph) createDataVertex(ctx context.Context, name, context string) (*ent.Data, error) {
	d, err := g.conn.Data.
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

func (g *Graph) createRelationship(ctx context.Context, from, to *ent.Data) (*ent.Data, error) {
	return from.Update().AddOwns(to).Save(ctx)
}

func (g *Graph) listDataEntities(ctx context.Context) ([]*ent.Data, error) {
	return g.conn.Data.Query().All(ctx)
}

func (g *Graph) getDataEntity(ctx context.Context, id int) (*ent.Data, error) {
	return g.conn.Data.Get(ctx, id)
}

func (g *Graph) getDataConnections(ctx context.Context, id int) ([]*ent.Data, error) {
	return g.conn.Data.Query().Where(data.ID(id)).QueryOwns().All(ctx)
}

func (g *Graph) index(ctx context.Context, lakes []*LakeConfig) error {
	for _, lake := range lakes {
		driver := connectors.GetDriver(fmt.Sprintf("%s%s", lake.Provider, lake.Store), lake.Address)
		if err := g.load(ctx, driver); err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph) reIndex(ctx context.Context, lakes []*LakeConfig) error {
	if err := g.deleteAll(ctx); err != nil {
		return err
	}
	if err := g.index(ctx, lakes); err != nil {
		return err
	}
	return nil
}

func (g *Graph) load(ctx context.Context, driver connectors.Driver) error {
	dbTopology, _ := driver.Index()

	for _, node := range dbTopology {
		_, err := g.fromNode(ctx, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Graph) fromNode(ctx context.Context, node *connectors.Node) (*ent.Data, error) {
	var propertyList []Property
	for k, v := range node.Properties {
		propertyList = append(propertyList, Property{Attribute: k, Value: v})
	}

	entityFrom, err := g.createDataVertex(ctx, node.Name, node.Context)
	if err != nil {
		return nil, err
	}
	if node.Children != nil {
		for _, childNode := range node.Children {
			entityTo, err := g.fromNode(ctx, childNode)
			if err != nil {
				return entityFrom, nil
			}
			_, err = g.createRelationship(ctx, entityFrom, entityTo)
			if err != nil {
				return entityTo, err
			}
			log.Printf("Created relationship between %s and %s\n", entityFrom.String(), entityTo.String())
		}
	}
	return entityFrom, nil
}

type Property struct {
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}
