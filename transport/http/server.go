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
	"time"
)

type AppServer struct {
	Server      *http.Server
	Logger      *zap.Logger
	Middlewares graphql.Middlewares

	mux         *http.ServeMux
	resolver    *graphql.Resolver
	healthCheck HealthCheck
}

func New(resolver *graphql.Resolver, logger *zap.Logger, hChk HealthCheck, mdlwares graphql.Middlewares) *AppServer {
	srv := &AppServer{
		Server:      nil,
		Logger:      logger,
		Middlewares: mdlwares,
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

	graphHandler := func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		requestID := RandomString(32)

		writer.Header().Set("request-id", requestID)
		defer func(begin time.Time) {
			took := time.Since(begin)
			s.Logger.Info("request complete", zap.Duration("took", took), zap.Int64("ms", took.Milliseconds()), zap.String("request-id", requestID))
		}(start)

		reqCtx := context.WithValue(request.Context(), "request-id", requestID)
		graphSrv.ServeHTTP(writer, request.WithContext(reqCtx))
	}

	graphEndpoint := "/"
	if config.PlaygroundEnabled {
		graphEndpoint = "/repogen"
		s.Logger.Info("playground enabled, setting GraphQL endpoint to /repogen")
		s.mux.HandleFunc("/", playground.Handler("GraphQL playground", "/repogen"))
	}

	var acc = graphHandler
	for _, m := range s.Middlewares {
		acc = m(acc)
	}

	s.mux.HandleFunc(graphEndpoint, acc)
}

func (s *AppServer) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *AppServer) ListenAndServe() error {
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
