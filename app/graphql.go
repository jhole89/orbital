package main

import (
	"context"
	"github.com/graphql-go/graphql"
	"log"
	"time"
)

func graphqlQuery(query string, schema graphql.Schema) *graphql.Result {
	graphqlCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       graphqlCtx,
	})
	if len(result.Errors) > 0 {
		log.Printf("errors: %v", result.Errors)
	}
	return result
}

var dataType = graphql.NewObject(
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

var dataQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DataQuery",
		Fields: graphql.Fields{
			"list": &graphql.Field{
				Type:        graphql.NewList(dataType),
				Description: "Get entity list",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return graph.ListEntities()
				},
			},
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
		},
	},
)

var dataSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: dataQuery,
	},
)

var adminQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AdminQuery",
		Fields: graphql.Fields{
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
		},
	},
)

var adminSchema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: adminQuery,
	},
)
