package connectors

import (
	"database/sql"
	"fmt"
	_ "github.com/segmentio/go-athena"
)

type AwsAthenaConnector struct {
	Address string
}

type AwsAthenaTableDetails struct {
	Attribute string
	Value sql.NullString
}

func (d *AwsAthenaConnector) Connect() *sql.DB {
	db, err := sql.Open("athena", d.Address)
	if err != nil {
		fmt.Printf("Unable to establish connection: %s\n", err)
	}
	return db
}

func (d *AwsAthenaConnector) Query(db *sql.DB, queryString string) *sql.Rows {
	rows, err := db.Query(queryString)

	if err != nil {
		fmt.Printf("Unable to perform query: %s\n", err)
	}
	return rows
}

func (d *AwsAthenaConnector) getDatabases(db *sql.DB) []string {
	query := d.Query(db, "SHOW SCHEMAS")
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

func (d *AwsAthenaConnector) getTables(db *sql.DB, database string) []string {
	query := d.Query(db, fmt.Sprintf("SHOW TABLES IN %s", database))
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

func (d *AwsAthenaConnector) describeTables(db *sql.DB, database string, table string) []string {
	query := d.Query(db, fmt.Sprintf("DESCRIBE %s.%s", database, table))
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

func (d *AwsAthenaConnector) getTableMeta(db *sql.DB, database string, table string) []string {
	query := d.Query(db, fmt.Sprintf("SHOW TBLPROPERTIES %s.%s", database, table))
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


func (d *AwsAthenaConnector) Index(db *sql.DB) []*Node {

	var databases []*Node
	for _, database := range d.getDatabases(db) {

		var tables []*Node
		for _, table := range d.getTables(db, database) {

			var fields []*Node
			for _, tableAttr := range d.describeTables(db, database, table) {
				fields = append(fields, &Node{Name: tableAttr,Context:  "field"})
			}
			tables = append(tables, &Node{Name: table, Context: "table", Children: fields})
		}
		databases = append(databases, &Node{Name: database, Context: "database", Children: tables})
	}
	return databases
}
