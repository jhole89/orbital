package connectors

import (
	"database/sql"
	"github.com/jhole89/discovery-backend/connectors/aws"
)

var Driver = map[string] map[string] func(address string) *sql.DB  {
	"aws": awsDrivers,
}

var awsDrivers = map[string] func(address string) *sql.DB {
	"athena": aws.Connect,
}
