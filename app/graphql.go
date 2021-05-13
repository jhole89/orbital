package main

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/jhole89/orbital/database"
	"net/http"
)

func resolveRelationships(id interface{}, p graphql.ResolveParams) ([]*database.Entity, error){
	context, ok := p.Args["context"].(string)
	if ok {
		entities, err := graph.GetRelationships(id, context)
		if err != nil {
			return nil, err
		}
		return entities, err
	}
	return nil, nil
}

func createEntityHandler() (*handler.Handler, error) {

	entityType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Entity",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.NewNonNull(graphql.ID),
					Description: "The ID of the Entity",
				},
				"name": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
					Description: "The name of the Entity",
				},
				"context": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
					Description: "The context of the Entity",
				},
			},
		},
	)

	nonNullEntityType := graphql.NewNonNull(entityType)
	nonNullEntityListType := graphql.NewNonNull(graphql.NewList(nonNullEntityType))

	connections := graphql.Field{
		Type:        nonNullEntityListType,
		Description: "Get Entity connections by ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.ID,
				Description: "The ID of the Entity",
			},
			"context": &graphql.ArgumentConfig{
				Type: graphql.String,
				Description: "The context of the Connection",
				DefaultValue: "owns",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			switch t := p.Source.(type) {
			case *database.Entity:
				id := t.ID
				return resolveRelationships(id, p)
			case interface{}:
				id, ok := p.Args["id"].(string)
				if ok {
					return resolveRelationships(id, p)
				}
				return nil, nil
			default:
				fmt.Println("Unknown Source Type received")
				return nil, nil
			}
		},
	}

	entityType.AddFieldConfig("connections", &connections)

	fields := graphql.Fields{
		"entity": &graphql.Field{
			Type:        entityType,
			Description: "Get entity by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
					Description: "The ID of the Entity",
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(string)
				if ok {
					entity, err := graph.GetEntity(id)
					if err != nil {
						return nil, err
					}
					return entity, nil
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        nonNullEntityListType,
			Description: "Get list of Entities",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return graph.ListEntities()
			},
		},
		"connections": &connections,
	}

	query := graphql.ObjectConfig{Name: "EntityQuery", Fields: fields}
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: graphql.NewObject(query)})
	if err != nil {
		return nil, err
	}
	return handler.New(&handler.Config{Schema: &schema, Pretty: true, GraphiQL: true}), nil
}

func createAdminHandler() (*handler.Handler, error) {
	var fields = graphql.Fields{
		"rebuild": &graphql.Field{
			Type:        graphql.String,
			Description: "Rebuild graph",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				err := reIndex(graph, conf.Lakes)
				if err != nil {
					return nil, err
				}
				return "started", nil
			},
		},
	}

	var query = graphql.ObjectConfig{Name: "AdminQuery", Fields: fields}
	var schema, err = graphql.NewSchema(graphql.SchemaConfig{Query: graphql.NewObject(query)})
	if err != nil {
		return nil, err
	}
	return handler.New(&handler.Config{Schema: &schema, Pretty: true, GraphiQL: true}), nil
}

func disableCors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, Accept-Encoding")

		h.ServeHTTP(w, r)
	})
}
