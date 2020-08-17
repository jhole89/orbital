package main

import (
    "fmt"
    "github.com/jhole89/discovery-backend/connectors"
    "github.com/jhole89/discovery-backend/database"
    "net/http"
)

func main() {

    var conf Config
    conf.getConf()

    database.Start(conf.Database.Endpoint)

    for _, lake := range conf.Lakes {
        connectors.Driver[lake.Provider][lake.Store](lake.Address)
        index()
    }

    registerRoutes()
    fmt.Printf("Server running at http://127.0.0.1:%d\n", conf.Service.Port)
    err := http.ListenAndServe(fmt.Sprintf(":%d", conf.Service.Port), nil)
    panicCheck(err, "Unable to serve")
}
