package connectors

import (
	"database/sql"
	"log"
	"strings"
)

type Driver interface {
	Connect(address string)
	Query(query string) *sql.Rows
	Index() []*Node
}

type Node struct {
	Name string
	Context string
	Children []*Node
}

func GetDriver(name string, address string) Driver {

	var supportedConnectors = map[string]Driver{
		"awsathena": &AwsAthenaConnector{},
	}

	c, ok := supportedConnectors[strings.ToLower(name)]

	if ok {
		c.Connect(address)
		return c
	} else {
		keys := make([]string, len(supportedConnectors))
		for k := range supportedConnectors {
			keys = append(keys, k)
		}
		log.Printf("Connecting to %s is not supported. Please specifiy a supported connector in your config.yaml.\nValid connectors's: %s", name, keys)
		return nil
	}
}
