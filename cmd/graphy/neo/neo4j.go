package neo

import (
	"github.com/cenkalti/backoff"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
	"graphy/cmd/graphy/config"
	"time"
)

var (
	MaxRetryAttempts = 5
	b                = backoff.NewConstantBackOff(3 * time.Second)
)

func NewBasicAuth() neo4j.AuthToken {
	return neo4j.BasicAuth(config.Neo4JUser, config.Neo4jPassword, "")
}

func NewDriver(authToken neo4j.AuthToken, logger *zap.Logger) (neo4j.Driver, func(), error) {
	driver, err := neo4j.NewDriver(config.Neo4JHost, authToken)
	if err != nil {
		return nil, nil, err
	}

	times := 0
	err = backoff.Retry(func() error {
		if err = driver.VerifyConnectivity(); err != nil {
			logger.Error("failed to verify connection to neo4j", zap.Error(err), zap.Int("attempt", times))
			if times > MaxRetryAttempts {
				return backoff.Permanent(err)
			}
			times++
			return err
		}
		return nil
	}, b)

	if err != nil {
		return nil, nil, err
	}

	logger.Info("connection to Neo4j established", zap.String("user", config.Neo4JUser), zap.String("uri", config.Neo4JHost))
	cleanup := func() {
		closeErr := driver.Close()
		if closeErr != nil {
			logger.Fatal("failed to close database connection", zap.Error(err))
		}
		logger.Info("closed neo4j driver", zap.String("neo4j", config.Neo4JHost))
	}

	return driver, cleanup, nil
}

type HealthCheck func() error

func (h HealthCheck) Do() error {
	return h()
}

func NewHealthCheck(driver neo4j.Driver) HealthCheck {
	return driver.VerifyConnectivity
}
