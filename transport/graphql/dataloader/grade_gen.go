// Code generated by dlgen, DO NOT EDIT.

package dataloader

import (
	"context"
	"errors"
	"graphy/pkg/competition/rounds"
	"graphy/transport/graphql/model"
	"time"
)

type GradeRoundLoaderProvider func(con context.Context) GradeRoundLoader

func NewGradeDLoader(svc rounds.Service) GradeRoundLoaderProvider {
	return func(ctx context.Context) GradeRoundLoader {
		return GradeRoundLoader{
			fetch: func(keys []string) ([][]*model.Round, []error) {
				// TODO: Generate the service here
				m, err := svc.FindGradeRounds(ctx, keys)

				if err != nil {
					// TODO: Make error message better
					return nil, []error{errors.New("")}
				}

				var out = make([][]*model.Round, len(keys))
				for index, keyID := range keys {
					resolvedVal := m[keyID]
					var decodedItems = make([]*model.Round, len(resolvedVal))
					for i, rawItem := range resolvedVal {
						decoded, _ := decode(rawItem)
						decodedItems[i] = &decoded
					}

					out[index] = decodedItems
				}

				return out, nil
			},
			wait:     1 * time.Millisecond,
			maxBatch: 100,
		}
	}
}

func decode(m rounds.Round) (model.Round, error) {
	return model.Round(m), nil
}