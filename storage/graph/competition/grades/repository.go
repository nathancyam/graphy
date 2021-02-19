package grade

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

func NewRepository(driver neo4j.Driver, logger *zap.Logger) *Repository {
	return &Repository{driver: driver, logger: logger}
}

func (r Repository) FindByID(ctx context.Context, id string) (*grades.Grade, error) {
	var model grades.Grade

	_, err := graph.WithReadConnection(r.driver, func(tx neo4j.Transaction) (interface{}, error) {
		out, err := tx.Run(`MATCH (grade:Grade) WHERE res.id = $id RETURN grade LIMIT 1`, map[string]interface{}{
			"id": id,
		})
		if err != nil {
			return nil, err
		}

		if out.Next() {
			node, ok := out.Record().Get("grade")
			if !ok {
				return nil, errors.New("")
			}

			if err = hydrateStruct(&model, node); err != nil {
				return nil, err
			}
		}

		if err := out.Err(); err != nil {
			return nil, err
		}

		return model, nil
	})

	if err != nil {
		return nil, err
	}

	return &model, nil
}

func toList(i interface{}) ([]grades.Grade, error) {
	col, ok := i.([]interface{})
	if !ok {
		return nil, graph.ErrNotSlice
	}

	var out = make([]grades.Grade, len(col))
	for index, i := range col {
		j, ok := i.(grades.Grade)
		if !ok {
			return nil, graph.ErrUnmarshal
		}
		out[index] = j
	}
	return out, nil
}

func hydrateStruct(model *grades.Grade, val interface{}) error {
	node, ok := val.(neo4j.Node)
	if !ok {
		return graph.ErrNotNode
	}

	return mapstructure.Decode(node.Props(), model)
}

