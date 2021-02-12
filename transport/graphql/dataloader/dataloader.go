package dataloader

import (
	"context"
	"net/http"
)

const loaderKey = "dataloaders"

type GradeMiddleware func(next http.HandlerFunc) http.HandlerFunc

type Loaders struct {
	RoundsByID GradeRoundLoader
}

func NewDataloaderMiddleware(g GradeDataLoaderProvider) GradeMiddleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			reqCtx := request.Context()

			ctx := context.WithValue(reqCtx, loaderKey, &Loaders{
				RoundsByID: g(reqCtx),
			})

			r := request.WithContext(ctx)
			next(writer, r)
		}
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loaderKey).(*Loaders)
}
