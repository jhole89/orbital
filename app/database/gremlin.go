package database

import (
	"encoding/json"
	"fmt"
	"github.com/jhole89/orbital/database/gremlin-rest"
	"github.com/schwartzmx/gremtune"
	"log"
	"time"
)

type Gremlin struct {
	Client gremlinClient
}

type gremlinClient interface {
	Execute(query string) ([]gremtune.Response, error)
}

func newGremlin(dsn string) (Graph, error) {

	errs := make(chan error)
	go func(chan error) {
		err := <-errs
		log.Printf("Lost connection to the database: %s\n", err.Error())
	}(errs)

	dialer := gremtune.NewDialer(dsn)

	retryCount := 10
	for {
		log.Println("Connecting to Gremlin database at: " + dsn)
		conn, err := gremtune.Dial(dialer, errs)
		if err != nil {
			if retryCount == 0 {
				log.Println("Unable to connect to Gremlin database at: " + err.Error())
				return nil, err
			}

			log.Printf("Could not connect to Gremlin server. Wait 2 seconds. %d retries left...\n", retryCount)
			retryCount--
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Connected to Gremlin database at: " + dsn)
			return &Gremlin{Client: &conn}, nil
		}
	}
}

func (g *Gremlin) Clean() error {
	_, err := g.Query("g.V().drop().iterate()")
	if err != nil {
		return err
	}
	log.Println("All vertices deleted, database is now empty.")
	return nil
}

func (g *Gremlin) CreateEntity(e *Entity) (*Entity, error) {
	queryString := fmt.Sprintf("g.addV('%s').property('name', '%s').property('context', '%s')", e.Context, e.Name, e.Context)
	for _, property := range e.Properties {
		queryString += fmt.Sprintf(".property('%s', '%s')", property.Attribute, property.Value)
	}
	resp, err := g.Query(queryString)
	if err != nil {
		return nil, err
	}

	var vlc gremlin_rest.VertexList
	if err := json.Unmarshal(resp, &vlc); err != nil {
		return nil, err
	}
	e.ID = vlc.Value[0].Value.ID.Value

	log.Printf("Created Entity: {ID: %v, Name: %s, Context: %s}\n", e.ID, e.Name, e.Context)
	return e, nil
}

func (g *Gremlin) CreateRelationship(r *Relationship) (*Relationship, error) {
	resp, err := g.Query(fmt.Sprintf("g.addE('%s').from(g.V(%v)).to(g.V(%v))", r.Context, r.From.ID, r.To.ID))
	if err != nil {
		return nil, err
	}

	var elc gremlin_rest.EdgeList
	if err := json.Unmarshal(resp, &elc); err != nil {
		return nil, err
	}
	r.ID = elc.Value[0].Value.ID.Value

	log.Printf("Created Relationship: {ID: %v, Context: %s, From: %s (ID: %v), To: %s (ID: %v)}\n", r.ID, r.Context, r.From.Name, r.From.ID, r.To.Name, r.To.ID)
	return r, nil
}

func (g *Gremlin) GetEntity(id interface{}) (*Entity, error) {
	resp, err := g.Query(fmt.Sprintf("g.V(%v).properties()", id))
	if err != nil {
		return nil, err
	}
	var r gremlin_rest.VertexPropertyList
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}
	var e = Entity{ID: id}
	for _, prop := range r.Value {
		switch prop.Value.Label {
			case "name":
				e.Name = prop.Value.Value
			case "context":
				e.Context = prop.Value.Value
		}
	}
	return &e, nil
}

func (g *Gremlin) GetRelationships(id interface{}, context string) ([]*Entity, error) {
	//resp, err := g.Query(fmt.Sprintf("g.V(%v).out('%s')", id, context))
	return nil, nil
}

func (g *Gremlin) ListEntities() ([]*Entity, error) {
	resp, err := g.Query(fmt.Sprintf("g.V()"))
	if err != nil {
		return nil, err
	}
	var r gremlin_rest.VertexList
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	var entities []*Entity
	for _, ent := range r.Value {
		entity, err := g.GetEntity(ent.Value.ID.Value)
		if err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
}

func (g *Gremlin) Query(queryString string) ([]byte, error) {
	resp, err := g.Client.Execute(queryString)
	if err != nil {
		log.Printf("Unable to execute query: %s. Err: %s\n", queryString, err.Error())
		return nil, err
	}
	return marshallResponse(resp)
}

func marshallResponse(resp []gremtune.Response) ([]byte, error) {
	j, err := json.Marshal(resp[0].Result.Data)
	if err != nil {
		return nil, err
	}
	return j, nil
}
