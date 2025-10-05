package saga

import "fmt"

type SagaStep struct {
	Name       string
	Action     func() error
	Compensate func()
}

type Saga struct {
	steps []SagaStep
}

func (s *Saga) AddStep(name string, action func() error, compensate func()) {
	s.steps = append(s.steps, SagaStep{
		Name:       name,
		Action:     action,
		Compensate: compensate,
	})
}

func (s *Saga) Execute() error {
	var executed []SagaStep
	for _, step := range s.steps {
		if err := step.Action(); err != nil {
			for i := len(executed) - 1; i >= 0; i-- {
				executed[i].Compensate()
			}
			return fmt.Errorf("saga step %s failed: %w", step.Name, err)
		}
		executed = append(executed, step)
	}
	return nil
}
