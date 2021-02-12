package dataloader

import (
	"context"
	"errors"
	"graphy/pkg/competition/rounds"
	"graphy/transport/graphql/model"
	"time"
)

type GradeDataLoaderProvider func(ctx context.Context) GradeRoundLoader

func NewGradeDLoader(roundSvc rounds.Service) GradeDataLoaderProvider {
	return func(ctx context.Context) GradeRoundLoader {
		return GradeRoundLoader{
			fetch: func(gradeIDs []string) ([][]*model.Round, []error) {
				m, err := roundSvc.FindGradeRounds(ctx, gradeIDs)
				if err != nil {
					return nil, []error{errors.New("")}
				}

				var out = make([][]*model.Round, len(gradeIDs))
				for index, id := range gradeIDs {
					rr := m[id]
					var rs = make([]*model.Round, len(rr))
					for j, r := range rr {
						o := model.Round(r)
						rs[j] = &o
					}

					out[index] = rs
				}

				return out, nil
			},
			wait:     1 * time.Millisecond,
			maxBatch: 100,
		}
	}
}
