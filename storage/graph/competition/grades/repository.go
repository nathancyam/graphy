package grades

import (
	"context"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
	"graphy/pkg/competition/grades"
	"graphy/storage/graph"
)

type Repository struct {
	logger *zap.Logger
	driver neo4j.Driver
}

func (r Repository) FindByID(ctx context.Context, id string) (*grades.Grade, error) {
	var grade grades.Grade

	_, err := graph.WithReadConnection(r.driver, func(tx neo4j.Transaction) (interface{}, error) {
		out, err := tx.Run(`MATCH (g:Grade { id: $id }) RETURN g LIMIT 1`, map[string]interface{}{
			"id": id,
		})
		if err != nil {
			return nil, err
		}

		if out.Next() {
			rec, ok := out.Record().Get("g")
			if !ok {
				return nil, errors.New("")
			}
			node, ok := rec.(neo4j.Node)
			if !ok {
				return nil, errors.New("")
			}
			if err := mapstructure.Decode(node.Props(), &grade); err != nil {
				return nil, err
			}
		}

		if err := out.Err(); err != nil {
			return nil, err
		}

		return grade, nil
	})

	if err != nil {
		return nil, err
	}

	return &grade, nil
}

func NewRepository(logger *zap.Logger, driver neo4j.Driver) *Repository {
	return &Repository{logger: logger, driver: driver}
}
