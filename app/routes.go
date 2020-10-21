package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"net/http"
	"time"
)

func registerRoutes() {

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(graphqlQuery(r.URL.Query().Get("query"), dataSchema)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	})
}

func graphqlQuery(query string, schema graphql.Schema) *graphql.Result {
	graphqlCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		Context:       graphqlCtx,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}
