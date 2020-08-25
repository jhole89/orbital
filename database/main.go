package database

import "net/http"

type Graph interface {
	Connect()
	Query(queryString string) string
	CreateEntity(e Entity) string
	CreateRelationship(r Relationship) string
	Read(w http.ResponseWriter) []byte
}

type Entity struct {
	Context    string
	Name       string
	Properties []Property
}

type Property struct {
	Attribute string
	Value string
}

type Relationship struct {
	From    Entity
	To      Entity
	Context string
}
