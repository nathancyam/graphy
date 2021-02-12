package graphql

import (
	"context"
	"go.uber.org/zap"
	"graphy/pkg/competition/grades"
	"graphy/pkg/competition/rounds"
	"time"
)

type Resolver struct {
	RoundService rounds.Service
	GradeSvc     *grades.Service
	Logger       *zap.Logger
}

func NewResolver(rounds rounds.Service, gradeSvc *grades.Service, logger *zap.Logger) *Resolver {
	return &Resolver{
		Logger:       logger,
		RoundService: rounds,
		GradeSvc:     gradeSvc,
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
