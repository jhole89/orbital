package database

import (
	"github.com/northwesternmutual/grammes/model"
	"github.com/northwesternmutual/grammes/query"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockGremlinClient struct {
	testVertexID  interface{}
	testVertexIDs []interface{}
	queryResponse [][]byte
}

func (m *mockGremlinClient) AddVertex(label string, properties ...interface{}) (model.Vertex, error) {
	return model.Vertex{
		Type: "g:Vertex",
		Value: model.VertexValue{
			ID: model.ID{
				Type:  "g:Int64",
				Value: m.testVertexID,
			},
			Label:      label,
			Properties: nil,
		},
	}, nil
}
func (m *mockGremlinClient) AllVertices() ([]model.Vertex, error) {
	var vertices = make([]model.Vertex, len(m.testVertexIDs))
	for _, id := range m.testVertexIDs {
		vertices = append(vertices, model.Vertex{
			Type: "g:Vertex",
			Value: model.VertexValue{
				ID: model.ID{
					Type:  "g:Int64",
					Value: id,
				},
				Label:      "test",
				Properties: nil,
			},
		})
	}
	return vertices, nil
}
func (m *mockGremlinClient) DropAll() error {
	return nil
}
func (m *mockGremlinClient) ExecuteQuery(queryObj query.Query) ([][]byte, error) {
	return m.queryResponse, nil
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

	m := mockGremlinClient{
		testVertexID: 1234,
		queryResponse: [][]byte{
			[]byte("{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:VertexProperty\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":1},\"label\":\"name\",\"value\":\"mydb\"}},{\"@type\":\"g:VertexProperty\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":2},\"label\":\"context\",\"value\":\"database\"}}]}"),
		},
	}
	g := Gremlin{&m}
	expected := Entity{
		ID:         1234,
		Context:    "database",
		Name:       "mydb",
		Properties: nil,
	}

	entity, err := g.CreateEntity(&Entity{Context: "database", Name: "mydb"})

	asserter.NoError(err)
	asserter.Equal(&expected, entity)
}

func TestGremlin_GetEntity(t *testing.T) {
	asserter := assert.New(t)
	m := mockGremlinClient{
		testVertexID: 1234,
		queryResponse: [][]byte{
			[]byte("{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:VertexProperty\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":1},\"label\":\"name\",\"value\":\"mydb\"}},{\"@type\":\"g:VertexProperty\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":2},\"label\":\"context\",\"value\":\"database\"}}]}"),
		},
	}
	g := Gremlin{&m}
	expectedEnt := Entity{ID: 1234, Context: "database", Name: "mydb"}

	ent, err := g.GetEntity(m.testVertexID)

	asserter.NoError(err)
	asserter.NotNil(ent)
	asserter.Equal(&expectedEnt, ent)
}

//func TestGremlin_ListEntities(t *testing.T) {
//	asserter := assert.New(t)
//	m := mockGremlinClient{
//		testVertexIDs: []interface{}{1234, 4567},
//		//queryResponse: [][]byte{
//		//	[]byte("{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:VertexProperty\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":1},\"label\":\"name\",\"value\":\"mydb\"}},{\"@type\":\"g:VertexProperty\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":2},\"label\":\"context\",\"value\":\"database\"}}]}"),
//		//},
//	}
//	g := Gremlin{&m}
//	expectedEnts := []*Entity{{ID: 1234, Context: "database", Name: "mydb"}, {ID: 5678, Context: "table", Name: "mytable"}}
//
//	ents, err := g.ListEntities()
//
//	asserter.NoError(err)
//	asserter.NotNil(ents)
//	asserter.Equal(expectedEnts, ents)
//}

//func TestGremlin_CreateRelationship(t *testing.T) {
//	asserter := assert.New(t)
//	expectedResponse :=[]byte(
//		"{\"@type\":\"g:List\",\"@value\":[{\"@type\":\"g:Edge\",\"@value\":{\"id\":{\"@type\":\"g:Int64\",\"@value\":39},\"label\":\"has_table\",\"inVLabel\":\"table\",\"outVLabel\":\"database\",\"inV\":{\"@type\":\"g:Int64\",\"@value\":10},\"outV\":{\"@type\":\"g:Int64\",\"@value\":8}}}]}",
//	)
//
//	m := mockGremlinClient{queryResponse: expectedResponse}
//	g := Gremlin{&m}
//	e1 := Entity{
//		Context:    "database",
//		Name:       "some-database",
//		Properties: nil,
//	}
//	e2 := Entity{
//		Context:    "table",
//		Name:       "some-table",
//		Properties: nil,
//	}
//	r := Relationship{
//		From:    &e1,
//		To:      &e2,
//		Context: "has_table",
//	}
//
//	resp, err := g.CreateRelationship(&r)
//
//	asserter.NoError(err)
//	asserter.Equal(expectedResponse, resp)
//}

func TestGremlin_GetRelationships(t *testing.T) {}
