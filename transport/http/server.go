package http

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"go.uber.org/zap"
	"graphy/cmd/graphy/config"
	"graphy/transport/graphql"
	"graphy/transport/graphql/generated"
	"net/http"
	"strconv"
)

type AppServer struct {
	Server      *http.Server
	Logger      *zap.Logger
	Middlewares Middlewares

	mux         *http.ServeMux
	resolver    *graphql.Resolver
	healthCheck HealthCheck
}

func New(resolver *graphql.Resolver, logger *zap.Logger, hChk HealthCheck, middlewares Middlewares) *AppServer {
	srv := &AppServer{
		Server:      nil, // Need to build everything else before. We could omit this, but it's more explicit this way.
		Logger:      logger,
		Middlewares: middlewares,
		resolver:    resolver,
		healthCheck: hChk,
		mux:         http.NewServeMux(),
	}

	srv.attachGqlgenHandlers()
	srv.attachHealthHandler()
	srv.build()

	return srv
}

func (s *AppServer) attachHealthHandler() {
	okRes := []byte(`{"status": "green", "reason": "ok"}`)
	failedRes := []byte(`{"status": "yellow", "reason": "neo4j connection could not be established"}`)

	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if s.healthCheck == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(okRes)
			return
		}

		if err := s.healthCheck.Do(); err != nil {
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write(failedRes)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(okRes)
		return
	})
}

func (s *AppServer) attachGqlgenHandlers() {
	graphSrv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: s.resolver}))

	graphEndpoint := "/"
	if config.PlaygroundEnabled {
		graphEndpoint = "/query"
		s.Logger.Info("playground enabled, setting GraphQL endpoint to /query")
		s.mux.HandleFunc("/", playground.Handler("GraphQL playground", graphEndpoint))
	}

	var acc = graphSrv.ServeHTTP
	for _, m := range s.Middlewares {
		acc = m(acc)
	}

	s.mux.HandleFunc(graphEndpoint, acc)
}

func (s *AppServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *AppServer) ListenAndServe() error {
	if s.Server == nil {
		s.Logger.Error("server struct not initialised. *AppServer.build() should be called.")
	}

	err := s.Server.ListenAndServe()
	if err != nil {
		if err != http.ErrServerClosed {
			s.Logger.Error(err.Error())
			return err
		}
		s.Logger.Info("server closed")
		return nil
	}

	return nil
}

func (s *AppServer) build() {
	s.Server = &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: s.mux,
	}
}
