package http

type HealthCheck interface {
	Do() error
}

type NoOpHealthCheck struct {
}

func (n NoOpHealthCheck) Do() error {
	return nil
}
