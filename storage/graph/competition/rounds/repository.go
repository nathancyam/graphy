package rounds

import (
	"context"
	"errors"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
	"graphy/pkg/competition/rounds"
	"graphy/storage/graph"
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
	res, err := graph.WithReadConnection(r.driver, func(tx neo4j.Transaction) (interface{}, error) {
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

	if err != nil {
		return nil, err
	}

	return toRounds(res)
}

func (r RoundRepository) FindGradeRounds(ctx context.Context, gradeIDs []string) (map[string][]rounds.Round, error) {
	var res = make(map[string][]rounds.Round, len(gradeIDs))

	_, err := graph.WithReadConnection(r.driver, func(tx neo4j.Transaction) (interface{}, error) {
		out, err := tx.Run(`
MATCH (g:Grade)-[:HAS_ROUND]->(r:Round)
WHERE g.id IN $gradeIDs
RETURN { id: g.id, items: COLLECT(r) } as out
`, map[string]interface{}{
			"gradeIDs": gradeIDs,
		})
		if err != nil {
			return nil, err
		}

		if out.Next() {
			o, ok := out.Record().Get("out")
			if !ok {
				return nil, errors.New("")
			}

			m := o.(map[string]interface{})
			id := m["id"].(string)
			items := m["items"].([]interface{})

			var rs = make([]rounds.Round, len(items))
			for index, i := range items {
				var r rounds.Round
				n, ok := i.(neo4j.Node)
				if !ok {
					return nil, graph.ErrNotNode
				}

				if err := mapstructure.Decode(n.Props(), &r); err != nil {
					return nil, err
				}
				rs[index] = r
			}

			res[id] = rs
		}

		return nil, err
	})

	return res, err
}

func toRounds(i interface{}) ([]rounds.Round, error) {
	col, ok := i.([]interface{})
	if !ok {
		return nil, graph.ErrNotSlice
	}

	var out = make([]rounds.Round, len(col))
	for index, i := range col {
		j, ok := i.(rounds.Round)
		if !ok {
			return nil, graph.ErrNotRound
		}
		out[index] = j
	}
	return out, nil
}

func hydrateStruct(round *rounds.Round, val interface{}) error {
	node, ok := val.(neo4j.Node)
	if !ok {
		return graph.ErrNotNode
	}

	return mapstructure.Decode(node.Props(), round)
}
