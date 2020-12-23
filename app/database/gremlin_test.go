package database

import (
	"github.com/schwartzmx/gremtune"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockGremlinClient struct {
	expectedResponse []byte
}

func (m *mockGremlinClient) Execute(_ string) ([]gremtune.Response, error) {
	return []gremtune.Response{{
		RequestID: "123",
		Status: gremtune.Status{
			Message:    "Success",
			Code:       200,
			Attributes: nil,
		},
		Result: gremtune.Result{
			Data: m.expectedResponse,
			Meta: nil,
		},
	}}, nil
}

func TestGremlin_Clean(t *testing.T) {
	asserter := assert.New(t)

	m := mockGremlinClient{}
	g := Gremlin{&m}
	err := g.Clean()

	asserter.NoError(err)
}

func TestGremlin_CreateEntity(t *testing.T) {
	asserter := assert.New(t)
	expectedResponse := []byte(
		"{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:Vertex\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":40},\"label\":\"database\"}}]}",
	)

	m := mockGremlinClient{expectedResponse}
	g := Gremlin{&m}
	e := Entity{
		Context:    "database",
		Name:       "analytics",
		Properties: nil,
	}
	resp, err := g.CreateEntity(e)

	asserter.NoError(err)
	asserter.Equal(expectedResponse, resp)
}

func TestGremlin_CreateRelationship(t *testing.T) {
	asserter := assert.New(t)
	expectedResponse := []byte(
		"{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:Edge\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":39},\"label\":\"has_table\",\"inVLabel\":\"table\",\"outVLabel\":\"database\",\"inV\":{\"@type\":\"g:Int64\",\"@value\":10},\"outV\":{\"@type\":\"g:Int64\",\"@value\":8}}}]}",
	)

	m := mockGremlinClient{expectedResponse}
	g := Gremlin{&m}
	e1 := Entity{
		Context:    "database",
		Name:       "some-database",
		Properties: nil,
	}
	e2 := Entity{
		Context:    "table",
		Name:       "some-table",
		Properties: nil,
	}
	r := Relationship{
		From:    &e1,
		To:      &e2,
		Context: "has_table",
	}

	resp, err := g.CreateRelationship(r)

	asserter.NoError(err)
	asserter.Equal(expectedResponse, resp)
}

func TestGremlin_Query(t *testing.T) {
	asserter := assert.New(t)

	m := mockGremlinClient{expectedResponse: []byte("{\"Query\":\"FakeQuery\"}")}
	g := Gremlin{&m}
	resp, err := g.Query("Not a real query")
	asserter.NoError(err)
	asserter.Equal([]byte("{\"Query\":\"FakeQuery\"}"), resp)
}

func Test_unmarshall(t *testing.T) {
	r := []gremtune.Response{{
		RequestID: "123",
		Status: gremtune.Status{
			Message:    "Success",
			Code:       200,
			Attributes: nil,
		},
		Result: gremtune.Result{
			Data: []byte("{\"Foo\":\"Bar\"}"),
			Meta: nil,
		},
	}}

	res, err := marshallResponse(r)

	asserter := assert.New(t)
	asserter.NoError(err)
	asserter.Equal([]byte("{\"Foo\":\"Bar\"}"), res)
}