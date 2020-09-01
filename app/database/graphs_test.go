package database

import (
	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestGetGraph(t *testing.T) {
	asserter := assert.New(t)

	var g *Gremlin
	monkey.PatchInstanceMethod(reflect.TypeOf(g), "Connect", func(_ *Gremlin, _ string) {})
	asserter.Implements((*Graph)(nil), GetGraph("gremlin", ""), "Gremlin graph should establish connection.")
	monkey.UnpatchInstanceMethod(reflect.TypeOf(g), "Connect")

	asserter.Nil(GetGraph("not-a-real-driver", ""), "Unsupported graph should return nil.")
}
