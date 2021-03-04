package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func createDataHandler() (*handler.Handler, error) {

	dataType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Data",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"context": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	fields := graphql.Fields{
		"data": &graphql.Field{
			Type:        dataType,
			Description: "Get entity by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"]
				if ok {
					data, err := graph.GetEntity(id)
					if err != nil {
						return nil, err
					}
					return data, nil
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type:        graphql.NewList(dataType),
			Description: "Get entity list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return graph.ListEntities()
			},
		},
		"connections": &graphql.Field{
			Type:        graphql.NewList(dataType),
			Description: "Get Entity connections by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"context": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"]
				if ok {
					context, ok := p.Args["context"].(string)
					if ok {
						data, err := graph.GetRelationships(id, context)
						if err != nil {
							return nil, err
						}
						return data, err
					}
				}
				return nil, nil
			},
		},
	}

	query := graphql.ObjectConfig{Name: "DataQuery", Fields: fields}
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
