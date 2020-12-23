package main

import (
	"context"
	"database/sql"
	"github.com/jhole89/orbital/connectors"
	"github.com/jhole89/orbital/ent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"sort"
	"testing"
)

func graphTestSetUp(ctx context.Context) *Graph {
	graph, err := newGraph("gremlin-rest", "http://127.0.0.1:8182")
	if err != nil {
		panic(err)
	}
	err = graph.deleteAll(ctx)
	if err != nil {
		panic(err)
	}
	return graph
}

func graphTestCleanUp(ctx context.Context, graph *Graph) {
	err := graph.deleteAll(ctx)
	if err != nil {
		panic(err)
	}
	err = graph.conn.Close()
	if err != nil {
		panic(err)
	}
}

func Test_createDataVertex(t *testing.T) {

	asserter := assert.New(t)
	ctx := context.Background()
	graph := graphTestSetUp(ctx)
	t.Cleanup(func() { graphTestCleanUp(ctx, graph) })

	data, err := graph.createDataVertex(ctx, "foo", "bar")

	asserter.NoError(err)
	asserter.NotNil(data)
	asserter.Equal("foo", data.Name)
	asserter.Equal("bar", data.Context)
}

func Test_createRelationship(t *testing.T) {
	asserter := assert.New(t)
	ctx := context.Background()

	graph := graphTestSetUp(ctx)
	t.Cleanup(func() { graphTestCleanUp(ctx, graph) })

	from, _ := graph.createDataVertex(ctx, "foo", "words")
	to, _ := graph.createDataVertex(ctx, "bar", "words")
	_, err := graph.createRelationship(ctx, from, to)
	fromEdges, _ := from.QueryOwns().All(ctx)

	asserter.NoError(err)
	asserter.Len(fromEdges, 1)
	asserter.Equal(to.ID, fromEdges[0].ID)
	asserter.Equal(to.Name, fromEdges[0].Name)
	asserter.Equal(to.Context, fromEdges[0].Context)
	asserter.Equal(to.Edges, fromEdges[0].Edges)
}

func Test_deleteAll(t *testing.T) {
	asserter := assert.New(t)
	ctx := context.Background()
	graph := graphTestSetUp(ctx)
	t.Cleanup(func() { graphTestCleanUp(ctx, graph) })

	data, _ := graph.createDataVertex(ctx, "deleteAll", "test")

	err := graph.deleteAll(ctx)
	asserter.NoError(err)

	d, err := graph.getDataEntity(ctx, data.ID)
	asserter.Error(err)
	asserter.Nil(d)
}

func Test_getDataEntity(t *testing.T) {
	asserter := assert.New(t)
	ctx := context.Background()
	graph := graphTestSetUp(ctx)
	t.Cleanup(func() { graphTestCleanUp(ctx, graph) })

	data, _ := graph.createDataVertex(ctx, "getDataEntity", "test")
	res, err := graph.getDataEntity(ctx, data.ID)

	asserter.NoError(err)
	asserter.Equal(data.ID, res.ID)
	asserter.Equal(data.Context, res.Context)
	asserter.Equal(data.Name, res.Name)
	asserter.Equal(data.Edges, res.Edges)
}

func Test_listDataEntities(t *testing.T) {
	asserter := assert.New(t)
	ctx := context.Background()
	graph := graphTestSetUp(ctx)
	t.Cleanup(func() { graphTestCleanUp(ctx, graph) })

	d1, _ := graph.createDataVertex(ctx, "listDataEntity1", "test")
	d2, _ := graph.createDataVertex(ctx, "listDataEntity2", "test")
	ds := []*ent.Data{d1, d2}
	sort.Slice(ds, func(i, j int) bool {
		return ds[i].ID < ds[j].ID
	})

	res, err := graph.listDataEntities(ctx)
	sort.Slice(res, func(i, j int) bool {
		return res[i].ID < res[j].ID
	})

	asserter.NoError(err)

	for i, _ := range ds {
		asserter.Equal(ds[i].ID, res[i].ID)
		asserter.Equal(ds[i].Context, res[i].Context)
		asserter.Equal(ds[i].Name, res[i].Name)
		asserter.Equal(ds[i].Edges, res[i].Edges)
	}
}

func TestGraph_getDataConnections(t *testing.T) {
	asserter := assert.New(t)
	ctx := context.Background()
	graph := graphTestSetUp(ctx)
	t.Cleanup(func() { graphTestCleanUp(ctx, graph) })

	from, _ := graph.createDataVertex(ctx, "foo", "words")
	to, _ := graph.createDataVertex(ctx, "bar", "words")
	_, _ = graph.createRelationship(ctx, from, to)

	edges, err := graph.getDataConnections(ctx, from.ID)

	asserter.NoError(err)
	asserter.Len(edges, 1)
	asserter.Equal(to.ID, edges[0].ID)
	asserter.Equal(to.Context, edges[0].Context)
	asserter.Equal(to.Name, edges[0].Name)
	asserter.Equal(to.Edges, edges[0].Edges)
}

func TestGraph_reIndex(t *testing.T) {
//	TODO
}

func TestGraph_index(t *testing.T)  {
	//TODO
}

func TestGraph_load(t *testing.T) {
	ctx := context.Background()
	graph := graphTestSetUp(ctx)
	t.Cleanup(func() { graphTestCleanUp(ctx, graph) })

	am := new(MockDriver)
	n := connectors.Node{
		Name:       "foo",
		Context:    "foos",
		Properties: nil,
		Children:   []*connectors.Node{{
			Name:       "bar",
			Context:    "bars",
			Properties: nil,
			Children: nil,
		}},
	}
	am.On("Index").Return([]*connectors.Node{&n}, nil)

	err := graph.load(ctx, am)
	assert.NoError(t, err)
}

type MockDriver struct {
	mock.Mock
}

func (a *MockDriver) Query(query string) (*sql.Rows, error) {
	panic("implement me")
}

func (a *MockDriver) Index() ([]*connectors.Node, error) {
	args := a.Called()
	return args.Get(0).([]*connectors.Node), args.Error(1)
}
