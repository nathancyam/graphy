package http

type HealthCheck interface {
	Do() error
}

type NoOpHealthCheck struct {
}

func NewNoOpHealthCheck() *NoOpHealthCheck {
	return &NoOpHealthCheck{}
}

func (n NoOpHealthCheck) Do() error {
	return nil
}
