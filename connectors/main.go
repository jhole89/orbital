package connectors

import "database/sql"

type Driver interface {
	Connect(address string) *sql.DB
	Query(query string) *sql.Rows
	Index() []*Node
}

type Node struct {
	Name string
	Context string
	Children []*Node
}
