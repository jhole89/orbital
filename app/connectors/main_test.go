package connectors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDriver(t *testing.T) {

	asserter := assert.New(t)

	asserter.Implements((*Driver)(nil), GetDriver("awsathena", ""), "Supported driver should establish connection.")
	asserter.Nil(GetDriver("not-a-real-driver", ""), "Unsupported driver should return nil.")
}
