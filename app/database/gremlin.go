package database

import (
	"github.com/northwesternmutual/grammes"
	"github.com/northwesternmutual/grammes/model"
	"github.com/northwesternmutual/grammes/query"
	"log"
	"strings"
	"time"
)

type Gremlin struct {
	Client gremlinClient
}

type gremlinClient interface {
	AddVertex(label string, properties ...interface{}) (model.Vertex, error)
	AllVertices() ([]model.Vertex, error)
	DropAll() error
	ExecuteQuery(queryObj query.Query) ([][]byte, error)
}

func newGremlin(dsn string) (Graph, error) {

	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Printf("Lost connection to the database: %s\n", err.Error())
	}(errs)

	retryCount := 10
	for {
		log.Println("Connecting to Gremlin database at: " + dsn)
		conn, err := grammes.DialWithWebSocket(dsn, grammes.WithErrorChannel(errs))
		if err != nil {
			if retryCount == 0 {
				log.Println("Unable to connect to Gremlin database at: " + dsn)
				return &Gremlin{}, err
			}

			log.Printf("Could not connect to Gremlin server. Wait 2 seconds. %d retries left...\n", retryCount)
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Connected to Gremlin database at: " + dsn)
			return &Gremlin{Client: conn}, nil
		}
	}
}

func (g *Gremlin) Clean() error {
	err := g.Client.DropAll()
	if err != nil {
		return err
	}
	log.Println("All vertices deleted, database is now empty.")
	return nil
}

func (g *Gremlin) CreateEntity(e *Entity) (*Entity, error) {
	vertex, err := g.Client.AddVertex(e.Context, "name", e.Name, "context", e.Context)
	if err != nil {
		return &Entity{}, err
	}
	ent, err := g.GetEntity(vertex.ID())
	if err != nil {
		return &Entity{}, err
	}
	log.Printf("Created Entity: {ID: %v, Name: %s, Context: %s}\n", ent.ID, ent.Name, ent.Context)
	return ent, nil
}

func (g *Gremlin) CreateRelationship(r *Relationship) (*Relationship, error) {
	t := grammes.Traversal()
	resp, err := g.Client.ExecuteQuery(t.AddE(r.Context).From(t.V(r.From.ID)).To(t.V(r.To.ID)))
	if err != nil {
		return &Relationship{}, err
	}
	edges, err := grammes.UnmarshalEdgeList(resp)
	if err != nil {
		return &Relationship{}, err
	}
	r.ID = edges[0].ID()

	log.Printf("Created Relationship: {ID: %v, Context: %s, From: %s (ID: %v), To: %s (ID: %v)}\n", r.ID, r.Context, r.From.Name, r.From.ID, r.To.Name, r.To.ID)
	return r, nil
}

func (g *Gremlin) GetEntity(id interface{}) (*Entity, error) {
	resp, err := g.Client.ExecuteQuery(grammes.Traversal().V(id).Properties())
	if err != nil {
		return &Entity{}, err
	}

	var e = Entity{ID: id}

	props, err := grammes.UnmarshalPropertyList(resp)
	if err != nil {
		return &Entity{}, err
	}

	for _, p := range props {
		v := p.GetValue().(string)
		switch strings.ToLower(p.Value.Label) {
		case "name":
			e.Name = v
		case "context":
			e.Context = v
		}
	}
	return &e, nil
}

func (g *Gremlin) GetRelationships(id interface{}, context string) ([]*Entity, error) {
	var entities []*Entity
	resp, err := g.Client.ExecuteQuery(grammes.Traversal().V(id).Out(context))
	if err != nil {
		return entities, err
	}
	vertices, err := grammes.UnmarshalVertexList(resp)
	if err != nil {
		return entities, err
	}
	for _, v := range vertices {
		ent, err := g.GetEntity(v.ID())
		if err != nil {
			return entities, err
		}
		entities = append(entities, ent)
	}
	return entities, nil
}

func (g *Gremlin) ListEntities() ([]*Entity, error) {
	vertices, err := g.Client.AllVertices()
	var entities []*Entity
	if err != nil {
		return entities, err
	}

	for _, ent := range vertices {
		entity, err := g.GetEntity(ent.ID())
		if err != nil {
			return entities, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}
