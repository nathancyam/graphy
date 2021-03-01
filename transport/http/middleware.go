package http

import (
	"context"
	"go.uber.org/zap"
	"graphy/store/neo4j"
	"graphy/transport/graphql/dataloader"
	"net/http"
	"time"
)

type Middlewares []func(next http.HandlerFunc) http.HandlerFunc

type RequestIDMiddleware func(next http.HandlerFunc) http.HandlerFunc

func NewRequestIDMiddleware(log *zap.Logger) RequestIDMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			requestID := RandomString(32)
			writer.Header().Set("request-id", requestID)

			defer func(begin time.Time) {
				took := time.Since(begin)
				log.Info("request complete", zap.Duration("took", took), zap.Int64("ms", took.Milliseconds()), zap.String("request-id", requestID))
			}(start)

			reqCtx := context.WithValue(request.Context(), "request-id", requestID)
			next(writer, request.WithContext(reqCtx))
		}
	}
}

func NewMiddlewares(g dataloader.GradeMiddleware, rID RequestIDMiddleware, neo4j neo4j.Middleware) Middlewares {
	// Middleware order should be:
	//  Request ID, Neo4j, Authentication, Dataloader
	return []func(next http.HandlerFunc) http.HandlerFunc{g, neo4j, rID}
}

