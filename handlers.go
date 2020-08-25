package main

import (
	"net/http"
)

func readHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	res := graph.Read(w)

	w.Header().Set("Content-Context", "application/json")
	_, err := w.Write(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
