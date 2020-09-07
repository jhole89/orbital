package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig_getConf(t *testing.T) {

	var conf Config
	conf.getConf()

	asserter := assert.New(t)

	asserter.NotEmpty(conf.Database.Type, "Database Type should not be empty.")
	asserter.NotEmpty(conf.Database.Endpoint, "Database Endpoint should not be empty.")

	asserter.NotEmpty(conf.Service.Port, "Service Port should not be empty.")

	for _, l := range conf.Lakes {
		asserter.NotEmpty(l.Provider, "Lake Provider should not be empty.")
		asserter.NotEmpty(l.Store, "Lake Store should not be empty.")
		asserter.NotEmpty(l.Address, "Lake Address should not be empty.")
	}
}
