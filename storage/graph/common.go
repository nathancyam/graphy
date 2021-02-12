package graph

import "github.com/neo4j/neo4j-go-driver/neo4j"

func WithReadConnection(driver neo4j.Driver, work func(tx neo4j.Transaction) (interface{}, error)) (interface{}, error) {
	session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}

	return session.ReadTransaction(work)
}

func WithWriteConnection(driver neo4j.Driver, work func(tx neo4j.Transaction) (interface{}, error)) (interface{}, error) {
	session, err := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	if err != nil {
		return nil, err
	}

	return session.WriteTransaction(work)
}
