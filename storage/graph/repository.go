package graph

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
	"graphy/pkg/rounds"
)

var RoundsQuery = `MATCH (res:Round) WHERE res.id IN $roundIDs RETURN res LIMIT 10`

type RoundRepository struct {
	logger *zap.Logger
	driver neo4j.Driver
}

func NewRoundRepository(driver neo4j.Driver, logger *zap.Logger) *RoundRepository {
	return &RoundRepository{driver: driver, logger: logger}
}

func (r RoundRepository) FindRoundsByID(_ context.Context, roundIDs []string) ([]rounds.Round, error) {
	session, err := r.driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	if err != nil {
		return nil, err
	}

	defer session.Close()

	res, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		res, err := tx.Run(RoundsQuery, map[string]interface{}{
			"roundIDs": roundIDs,
		})
		if err != nil {
			return nil, err
		}

		var rs []interface{}
		for res.Next() {
			for _, v := range res.Record().Values() {
				var round rounds.Round
				if err = hydrateStruct(&round, v); err != nil {
					return nil, err
				}

				rs = append(rs, round)
			}
		}

		if res.Err() != nil {
			return nil, res.Err()
		}

		return rs, nil
	})

	return toRounds(res)
}

func toRounds(i interface{}) ([]rounds.Round, error) {
	col, ok := i.([]interface{})
	if !ok {
		return nil, ErrNotSlice
	}

	var out []rounds.Round
	for _, i := range col {
		j, ok := i.(rounds.Round)
		if !ok {
			return nil, ErrNotRound
		}
		out = append(out, j)
	}
	return out, nil
}

func hydrateStruct(round *rounds.Round, val interface{}) error {
	node, ok := val.(neo4j.Node)
	if !ok {
		return ErrNotNode
	}
	if err := mapstructure.Decode(node.Props(), round); err != nil {
		return fmt.Errorf("invalid node properties %w", err)
	}

	return nil
}
