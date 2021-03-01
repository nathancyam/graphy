package neo4j

import (
	"bytes"
	"context"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

const neo4jSession = "neo4jSession"

type Middleware func(next http.HandlerFunc) http.HandlerFunc

func readRewindRequest(r *http.Request) []byte {
	var b []byte
	if r.Body != nil {
		b, _ = ioutil.ReadAll(r.Body)
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	return b
}

func NewNeo4jSessionMiddleware(l *zap.Logger, d neo4j.Driver) Middleware {
	introspection := []byte("Introspection")
	mutation := []byte(`"query":"mutation`)
	errMsg := []byte(`{ "error": "unable to process your request, please try again later" }`)

	readSessionCfg := neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead}
	writeSessionCfg := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite}

	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			// How slow/fast is this?
			b := readRewindRequest(request)
			isMutation := bytes.Contains(b, mutation)
			isIntrospection := bytes.Contains(b, introspection)

			if isIntrospection {
				next(writer, request)
				return
			}

			accessMode := readSessionCfg
			if isMutation {
				accessMode = writeSessionCfg
			}

			session, err := d.NewSession(accessMode)
			reqCtx := request.Context()
			requestID := reqCtx.Value("request-id").(string)

			if err != nil {
				l.Error("unable to create neo4j session", zap.Error(err), zap.String("request-id", requestID))
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Header().Set("Content-Type", "application/json")
				writer.Write(errMsg)
				return
			}

			defer func () {
				session.Close()
				l.Info("neo4j session closed for request", zap.String("request-id", requestID))
			}()

			requestCtx := context.WithValue(request.Context(), neo4jSession, session)

			l.Info("neo4j session created for request", zap.String("request-id", requestID))
			next(writer, request.WithContext(requestCtx))
		}
	}
}

func For(ctx context.Context) neo4j.Session {
	return ctx.Value(neo4jSession).(neo4j.Session)
}