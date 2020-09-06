package database

import (
	"github.com/schwartzmx/gremtune"
	"testing"
)

type MockGremlin struct {
	Gremlin
}

func (m *MockGremlin) runQuery() ([]gremtune.Response, error) {
	return []gremtune.Response{{RequestID: "abc"}}, nil
}

func TestGremlin_Clean(t *testing.T) {}

func TestGremlin_CreateEntity(t *testing.T) {}

func TestGremlin_CreateRelationship(t *testing.T) {}

func TestGremlin_Query(t *testing.T) {}

func TestGremlin_Read(t *testing.T) {}

func TestGremlin_httpQuery(t *testing.T) {}

func TestGremlin_runQuery(t *testing.T) {}

func TestNewGremlin(t *testing.T) {}

func Test_unmarshall(t *testing.T) {}
