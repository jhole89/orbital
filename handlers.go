package main

import (
	"github.com/jhole89/discovery-backend/connectors/aws"
	"github.com/jhole89/discovery-backend/database"
	"net/http"
)

func readHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
	}

	res := database.Read(w)

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func index() {

	res := aws.Query("SHOW DATABASES")
	aws.ReadRows(res)

}
