package main

import (
	"encoding/json"
	"net/http"
)

func dynamicQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Expected POST not GET", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, "Please send a request body", http.StatusBadRequest)
		return
	}

	type query struct {
		Query string `json:"query"`
	}
	var q query
	err := json.NewDecoder(r.Body).Decode(&q)

	res, err := (*graph).Query(q.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Context", "application/json")
	_, err = w.Write(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
