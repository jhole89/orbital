package main

import (
	"github.com/graphql-go/graphql"
)

var dataType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Data",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"context": &graphql.Field{
				Type: graphql.String,
			},
			//"properties": &graphql.Field{
			//	Type: graphql.NewList(propertyType),
			//},
		},
	},
)

//
//var propertyType = graphql.NewObject(
//	graphql.ObjectConfig{
//		Name: "Property",
//		Fields: graphql.Fields{
//			"attribute": &graphql.Field{
//				Type: graphql.String,
//			},
//			"value": &graphql.Field{
//				Type: graphql.String,
//			},
//		},
//	},
//)
//
var dataQuery = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "DataQuery",
		Fields: graphql.Fields{
			"data": &graphql.Field{
				Type:        dataType,
				Description: "Get data-entity by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						data, err := getDataEntity(p.Context, graph, id)
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
				Description: "Get data-entity list",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return listDataEntities(p.Context, graph)
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
