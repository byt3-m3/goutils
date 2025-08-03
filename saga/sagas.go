package saga

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
	"time"
)

type CTXFunc func(ctx context.Context) error

type Saga struct {
	Name      string
	steps     []*Step
	logger    *slog.Logger
	state     SagaState
	startTime time.Time
	endTime   time.Time
}

type NewSagaOpts func(*Saga)

// WithSteps sets the steps for a saga during its initialization.
func WithSteps(steps []*Step) NewSagaOpts {
	return func(s *Saga) {
		s.steps = steps
	}
}

// NewSaga initializes a new Saga with a specified name, logger, and optional configurations.
// name specifies the saga's name, logger handles logging, and opts allows for additional setup, such as pre-defined steps.
func NewSaga(name string, logger *slog.Logger, opts ...NewSagaOpts) *Saga {
	s := &Saga{
		Name:   name,
		steps:  make([]*Step, 0),
		logger: logger,
		state:  SagaStateInit,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// AddStep adds a new step to the saga with a name, action, compensation, and whether compensation is required.
func (s *Saga) AddStep(name string, action CTXFunc, compensation CTXFunc, shouldCompensate bool) {
	step := &Step{
		Name:             name,
		action:           action,
		compensation:     compensation,
		shouldCompensate: shouldCompensate,
		state:            SagaStepStateInit,
	}

	s.steps = append(s.steps, step)
}

// executeActions runs all the actions defined in the saga steps sequentially.
// If a step fails, compensations are executed for all completed steps in reverse order.
// Logs execution details at each step and returns an error if any step or compensation fails.
func (s *Saga) executeActions(ctx context.Context) error {

	s.logger.Info("executing saga",
		slog.String("saga_name", s.Name),
		slog.Int("steps_count", len(s.steps)),
	)
	var actionErrs error
	var compErrs error
	s.setState(SagaStateStarted)
	s.startTime = time.Now()
	for stepIndex, step := range s.steps {
		s.logger.Info("executing step",
			slog.String("step_name", step.Name),
		)
		step.SetState(SagaStepStateStarted)
		err := step.action(ctx)
		if err != nil {
			slog.Error("step action failed",
				slog.String("step_name", step.Name),
				slog.Any("error", err),
			)
			step.SetState(SagaStepStateError)
			s.setState(SagaStateErrored)
			actionErrs = errors.Wrap(err, fmt.Sprintf("step failed: %s", step.Name))
			err = s.executeCompensations(ctx, stepIndex)
			if err != nil {
				compErrs = errors.Wrap(err, "compensations failed")
				s.setState(SagaStatePoisoned)
				s.endTime = time.Now()

			}
			return errors.Wrap(compErrs, actionErrs.Error())

		}

		step.SetState(SagaStepStateDone)

	}

	s.setState(SagaStateDone)
	s.endTime = time.Now()

	return nil
}

// executeCompensations performs the compensation logic for all steps up to the given stepIndex in reverse order.
// It logs execution details and accumulates errors if compensations fail. Returns an error if any step fails to compensate.
func (s *Saga) executeCompensations(ctx context.Context, stepIndex int) error {

	var errs error
	for i := stepIndex; i >= 0; i-- {
		step := s.steps[i]
		if step.shouldCompensate {
			s.logger.Info("executing compensation step",
				slog.String("step_name", step.Name),
			)
			step.SetState(SagaStepStateCompensating)
			if err := step.compensation(ctx); err != nil {
				slog.Error("compensation step failed",
					slog.String("step_name", step.Name),
				)
				step.SetState(SagaStepStateError)
				errs = errors.Wrap(err, fmt.Sprintf("compensation failed for step %s", step.Name))
				continue
			}

			step.SetState(SagaStepStateCompensationDone)
		}
	}

	if errs != nil {
		return errs
	}

	return nil

}

func (s *Saga) setState(newState SagaState) {
	s.state = newState
}

func (s *Saga) duration() time.Duration {
	return s.endTime.Sub(s.startTime)
}
