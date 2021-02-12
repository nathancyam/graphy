//+build wireinject

package inject

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"graphy/cmd/graphy/neo"
	"graphy/pkg/competition/grades"
	"graphy/pkg/competition/rounds"
	grades2 "graphy/storage/graph/competition/grades"
	rounds2 "graphy/storage/graph/competition/rounds"
	"graphy/transport/graphql"
	"graphy/transport/graphql/dataloader"
	"graphy/transport/http"
)

var neoSet = wire.NewSet(
	neo.NewBasicAuth,
	neo.NewDriver,
	http.NewNoOpHealthCheck,
	wire.Bind(new(http.HealthCheck), new(*http.NoOpHealthCheck)),
)

var roundSet = wire.NewSet(
	rounds2.NewRoundRepository,
	wire.Bind(new(rounds.Repository), new(*rounds2.RoundRepository)),
	rounds.NewUpdateService,
	rounds.NewServiceImpl,
	wire.Bind(new(rounds.Service), new(*rounds.ServiceImpl)),
)

var gradeSet = wire.NewSet(
	grades2.NewRepository,
	wire.Bind(new(grades.Repository), new(*grades2.Repository)),
	grades.NewService,
)

var graphqlSet = wire.NewSet(
	graphql.NewResolver,
	dataloader.NewGradeDLoader,
	dataloader.NewDataloaderMiddleware,
	graphql.NewMiddlewares,
)

func InitialiseAppServer(logger *zap.Logger) (*http.AppServer, func(), error) {
	panic(wire.Build(
		neoSet,
		roundSet,
		gradeSet,
		graphqlSet,
		http.New,
	))
}
