package saga

import (
	"log"
)

type Saga struct {
	steps []SagaStep
}

type SagaStep struct {
	action       func() error
	compensation func() error
}

func NewSaga() *Saga {
	return &Saga{
		steps: make([]SagaStep, 0, 1),
	}
}

func (s *Saga) AddStep(action func() error, compensation func() error) {
	step := SagaStep{action, compensation}
	s.steps = append(s.steps, step)
}

func (s *Saga) Execute() (bool, error) {
	for i, step := range s.steps {
		if err := step.action(); err != nil {
			log.Println("Saga failed. Starting compensation transactions...")
			return false, s.Compensate(i)
		}
	}
	return true, nil
}

func (s *Saga) Compensate(step int) error {
	for i := step; i >= 0; i-- {
		if err := s.steps[i].compensation(); err != nil {
			return err
		}
	}
	return nil
}
