package connectors

import (
	"database/sql"
	"fmt"
	_ "github.com/segmentio/go-athena"
)

var err error

type AwsAthenaConnector struct {
	Connection *sql.DB
}

type AwsAthenaTableDetails struct {
	Attribute string
	Value sql.NullString
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
		var database string
		err := query.Scan(&database)
		if err != nil {
			fmt.Printf("Unable to scan database: %s\n", err)
		}
		databases = append(databases, database)
	}
	return databases
}

func (d *AwsAthenaConnector) getTables(database string) []string {
	query := d.Query(fmt.Sprintf("SHOW TABLES IN %s", database))
	var tables []string

	for query.Next() {
		var table string
		err := query.Scan(&table)
		if err != nil {
			fmt.Printf("Unable to scan table: %s\n", err)
		}
		tables = append(tables, table)
	}
	return tables
}

func (d *AwsAthenaConnector) describeTables(database string, table string) []string {
	query := d.Query(fmt.Sprintf("DESCRIBE %s.%s", database, table))
	var tableAttributes []string

	for query.Next() {
		var tableAttribute = AwsAthenaTableDetails{}

		err := query.Scan(&tableAttribute.Attribute, &tableAttribute.Value)
		if err != nil {
			fmt.Printf("Unable to scan TableAttributes: %s\n", err)
		}
		tableAttributes = append(tableAttributes, tableAttribute.Attribute)
	}
	return tableAttributes
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
			for _, tableAttr := range d.describeTables(database, table) {
				fields = append(fields, &Node{Name: tableAttr,Context:  "field"})
			}
			tables = append(tables, &Node{Name: table, Context: "table", Children: fields})
		}
		databases = append(databases, &Node{Name: database, Context: "database", Children: tables})
	}
	return databases
}
