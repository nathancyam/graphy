//+build wireinject

package inject

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"graphy/cmd/graphy/neo"
	"graphy/pkg/rounds"
	"graphy/storage/graph"
	"graphy/transport/graphql"
	"graphy/transport/http"
)

func InitialiseAppServer(logger *zap.Logger) (*http.AppServer, func(), error) {
	wire.Build(
		neo.NewBasicAuth,
		neo.NewDriver,
		neo.NewHealthCheck,
		graph.NewRoundRepository,
		wire.Bind(new(rounds.Repository), new(*graph.RoundRepository)),
		wire.Bind(new(http.HealthCheck), new(neo.HealthCheck)),
		graphql.NewResolver,
		http.New,
	)
	return &http.AppServer{}, nil, nil
}
