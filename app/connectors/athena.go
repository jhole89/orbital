package connectors

import (
	"database/sql"
	"fmt"
	_ "github.com/segmentio/go-athena"
	"strings"
)

type AwsAthena struct {
	Connection athenaClient
}

type athenaClient interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Close() error
}

func NewAwsAthena(dsn string) (Driver, error) {

	db, err := sql.Open("athena", dsn)

	if err != nil {
		fmt.Printf("Unable to establish connection: %s\n", err)
		return nil, err
	}

	return &AwsAthena{Connection: db}, nil
}

func (a *AwsAthena) Close() {
	_ = a.Connection.Close()
}

type AwsAthenaTableDetails struct {
	Column string
	Type   sql.NullString
}

type Column struct {
	Name string
	Type string
}

func (a *AwsAthena) Query(queryString string) (*sql.Rows, error) {

	rows, err := a.Connection.Query(queryString)

	if err != nil {
		fmt.Printf("Unable to execute query: %s\n", err)
		return nil, err
	}

	return rows, nil
}

func (a *AwsAthena) getDatabases() ([]string, error) {
	query, err := a.Query("SHOW SCHEMAS")

	if err != nil {
		return nil, err
	}

	var databases = make([]string, 0)

	for query.Next() {
		var databaseName string
		err := query.Scan(&databaseName)

		if err != nil {
			fmt.Printf("Unable to scan database: %s\n", err)
			return nil, err
		}
		databases = append(databases, strings.TrimSpace(databaseName))
	}
	return databases, nil
}

func (a *AwsAthena) getTables(database string) ([]string, error) {
	query, err := a.Query(fmt.Sprintf("SHOW TABLES IN %s", database))

	if err != nil {
		return nil, err
	}

	var tables = make([]string, 0)

	for query.Next() {

		var tabName string
		err := query.Scan(&tabName)

		if err != nil {
			fmt.Printf("Unable to scan table: %s\n", err)
			return nil, err
		}
		tables = append(tables, strings.TrimSpace(tabName))
	}
	return tables, nil
}

func (a *AwsAthena) describeTables(database, table string) ([]Column, error) {
	query, err := a.Query(fmt.Sprintf("DESCRIBE %s.%s", database, table))

	if err != nil {
		return nil, err
	}

	var columns = make([]Column, 0)

	for query.Next() {
		var tableAttribute = AwsAthenaTableDetails{}

		err := query.Scan(&tableAttribute.Column, &tableAttribute.Type)

		if err != nil {
			fmt.Printf("Unable to scan TableAttributes: %s\n", err)
			return nil, err
		}

		c := strings.Split(tableAttribute.Column, "\t")

		if len(c) == 2 {
			col := Column{Name: strings.TrimSpace(c[0]), Type: strings.TrimSpace(c[1])}
			columns = append(columns, col)
		}
	}
	return columns, nil
}

func (a *AwsAthena) Index() ([]*Node, error) {
	return index(a.getDatabases, a.getTables, a.describeTables)
}

func index(getDatabases func() ([]string, error), getTables func(string) ([]string, error), getColumns func(string, string) ([]Column, error)) ([]*Node, error) {
	dbs, err := getDatabases()
	if err != nil {
		return nil, err
	}

	var dbNodes = make([]*Node, 0)
	for _, database := range dbs {

		tables, err := getTables(database)
		if err != nil {
			return nil, err
		}

		var tblNodes = make([]*Node, 0)
		for _, table := range tables {

			fields, err := getColumns(database, table)
			if err != nil {
				return nil, err
			}

			var fieldNodes = make([]*Node, 0)
			for _, field := range fields {
				fieldNodes = append(fieldNodes, &Node{Name: field.Name, Context: "field", Properties: map[string]string{"data-type": field.Type}})
			}
			tblNodes = append(tblNodes, &Node{Name: table, Context: "table", Children: fieldNodes})
		}
		dbNodes = append(dbNodes, &Node{Name: database, Context: "database", Children: tblNodes})
	}
	return dbNodes, nil
}
