package graphql

import (
	"context"
	"go.uber.org/zap"
	"graphy/pkg/rounds"
	"time"
)

type Resolver struct {
	RoundRepo rounds.Repository
	Logger    *zap.Logger
}

func NewResolver(rounds rounds.Repository, logger *zap.Logger) *Resolver {
	return &Resolver{
		RoundRepo: rounds,
		Logger:    logger,
	}
}

func (r Resolver) LogDuration(ctx context.Context, resolver string) func() {
	start := time.Now()
	return func() {
		dur := time.Since(start)
		ms := dur.Milliseconds()
		requestID := ctx.Value("request-id")
		r.Logger.Info("completed resolver", zap.Duration("took", dur), zap.String("resolver", resolver), zap.Int64("ms", ms), zap.String("request-id", requestID.(string)))
	}
}
