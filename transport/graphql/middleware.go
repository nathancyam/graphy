package graphql

import (
	"graphy/transport/graphql/dataloader"
	"net/http"
)

type Middlewares []func(next http.HandlerFunc) http.HandlerFunc

func NewMiddlewares(g dataloader.GradeMiddleware) Middlewares {
	return []func(next http.HandlerFunc) http.HandlerFunc{g}
}
