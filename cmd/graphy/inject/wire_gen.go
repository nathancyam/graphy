// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package inject

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"graphy/cmd/graphy/neo"
	"graphy/pkg/competition/grades"
	rounds2 "graphy/pkg/competition/rounds"
	"graphy/store/neo4j"
	"graphy/store/neo4j/competition/grades"
	"graphy/store/neo4j/competition/rounds"
	"graphy/transport/graphql"
	"graphy/transport/graphql/dataloader"
	"graphy/transport/http"
)

// Injectors from wire.go:

func InitialiseAppServer(logger *zap.Logger) (*http.AppServer, func(), error) {
	authToken := neo.NewBasicAuth()
	driver, cleanup, err := neo.NewDriver(authToken, logger)
	if err != nil {
		return nil, nil, err
	}
	roundRepository := rounds.NewRoundRepository(driver, logger)
	updateService := rounds2.NewUpdateService(roundRepository)
	serviceImpl := rounds2.NewServiceImpl(updateService, roundRepository)
	repository := grade.NewRepository(driver, logger)
	service := grades.NewService(repository)
	resolver := graphql.NewResolver(serviceImpl, service, logger)
	noOpHealthCheck := http.NewNoOpHealthCheck()
	gradeRoundLoaderProvider := dataloader.NewGradeDLoader(serviceImpl)
	gradeMiddleware := dataloader.NewDataloaderMiddleware(gradeRoundLoaderProvider)
	requestIDMiddleware := http.NewRequestIDMiddleware(logger)
	middleware := neo4j.NewNeo4jSessionMiddleware(logger, driver)
	middlewares := http.NewMiddlewares(gradeMiddleware, requestIDMiddleware, middleware)
	appServer := http.New(resolver, logger, noOpHealthCheck, middlewares)
	return appServer, func() {
		cleanup()
	}, nil
}

// wire.go:

var neoSet = wire.NewSet(neo.NewBasicAuth, neo.NewDriver, http.NewNoOpHealthCheck, wire.Bind(new(http.HealthCheck), new(*http.NoOpHealthCheck)))

var roundSet = wire.NewSet(rounds.NewRoundRepository, wire.Bind(new(rounds2.Repository), new(*rounds.RoundRepository)), rounds2.NewUpdateService, rounds2.NewServiceImpl, wire.Bind(new(rounds2.Service), new(*rounds2.ServiceImpl)))

var gradeSet = wire.NewSet(grade.NewRepository, wire.Bind(new(grades.Repository), new(*grade.Repository)), grades.NewService)

var graphqlSet = wire.NewSet(graphql.NewResolver, dataloader.NewGradeDLoader, dataloader.NewDataloaderMiddleware, neo4j.NewNeo4jSessionMiddleware, http.NewRequestIDMiddleware, http.NewMiddlewares)
