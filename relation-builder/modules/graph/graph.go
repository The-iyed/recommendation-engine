package graph

import (
	"r-builder/config"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jDB struct {
	Driver neo4j.Driver
}

func NewNeo4jDB() (*Neo4jDB, error) {
	cfg := config.LoadConfig()
	driver, err := neo4j.NewDriver(cfg.Neo4jURI, neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""))
	if err != nil {
		return nil, err
	}
	return &Neo4jDB{Driver: driver}, nil
}

func (db *Neo4jDB) Close() {
	db.Driver.Close()
}

