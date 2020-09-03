package connectors

import (
	"database/sql"
	"fmt"
	_ "github.com/segmentio/go-athena"
	"strings"
)

var err error

type AwsAthenaConnector struct {
	Connection *sql.DB
}

type AwsAthenaTableDetails struct {
	Column string
	Type   sql.NullString
}

type Column struct {
	Name string
	Type string
}

func (d *AwsAthenaConnector) Connect(address string) {
	d.Connection, err = sql.Open("athena", address)
	if err != nil {
		fmt.Printf("Unable to establish connection: %s\n", err)
	}
}

func (d *AwsAthenaConnector) Query(queryString string) *sql.Rows {
	rows, err := d.Connection.Query(queryString)

	if err != nil {
		fmt.Printf("Unable to perform query: %s\n", err)
	}
	return rows
}

func (d *AwsAthenaConnector) getDatabases() []string {
	query := d.Query("SHOW SCHEMAS")
	var databases []string

	for query.Next() {
		var databaseName string
		err := query.Scan(&databaseName)
		if err != nil {
			fmt.Printf("Unable to scan database: %s\n", err)
		}
		databases = append(databases, strings.TrimSpace(databaseName))
	}
	return databases
}

func (d *AwsAthenaConnector) getTables(database string) []string {
	query := d.Query(fmt.Sprintf("SHOW TABLES IN %s", database))
	var tables []string

	for query.Next() {
		var tabName string
		err := query.Scan(&tabName)
		if err != nil {
			fmt.Printf("Unable to scan table: %s\n", err)
		}
		tables = append(tables, strings.TrimSpace(tabName))
	}
	return tables
}

func (d *AwsAthenaConnector) describeTables(database string, table string) []Column {
	query := d.Query(fmt.Sprintf("DESCRIBE %s.%s", database, table))
	var columns []Column

	for query.Next() {
		var tableAttribute = AwsAthenaTableDetails{}

		err := query.Scan(&tableAttribute.Column, &tableAttribute.Type)
		if err != nil {
			fmt.Printf("Unable to scan TableAttributes: %s\n", err)
		}
		c := strings.Split(tableAttribute.Column, "\t")

		if len(c) == 2 {
			col := Column{Name: strings.TrimSpace(c[0]), Type: strings.TrimSpace(c[1])}
			columns = append(columns, col)
		}
	}
	return columns
}

func (d *AwsAthenaConnector) getTableMeta(database string, table string) []string {
	query := d.Query(fmt.Sprintf("SHOW TBLPROPERTIES %s.%s", database, table))
	var tableMeta []string

	for query.Next() {
		var propAttr string
		var propValue sql.NullString
		err := query.Scan(&propAttr, &propValue)
		if err != nil {
			fmt.Printf("Unable to scan tableProp: %s\n", err)
		}
		fmt.Printf("---- ---- --- TableProp: %s -- %s\n", propAttr, propValue.String)
		tableMeta = append(tableMeta, propAttr)
	}
	return tableMeta
}

func (d *AwsAthenaConnector) Index() []*Node {

	var databases []*Node
	for _, database := range d.getDatabases() {

		var tables []*Node
		for _, table := range d.getTables(database) {

			var fields []*Node
			for _, field := range d.describeTables(database, table) {
				fields = append(fields, &Node{Name: field.Name, Context: "field", Properties: map[string]string{"data-type": field.Type}})
			}
			tables = append(tables, &Node{Name: table, Context: "table", Children: fields})
		}
		databases = append(databases, &Node{Name: database, Context: "database", Children: tables})
	}
	return databases
}
