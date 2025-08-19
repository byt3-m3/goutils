package saga

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"log/slog"
)

type Orchestrator struct {
	Sagas  map[string]*Saga
	Logger *slog.Logger
}

type NewOrchestratorOpts func(*Orchestrator)

func WithSagas(sagas []*Saga) NewOrchestratorOpts {
	return func(o *Orchestrator) {
		for _, saga := range sagas {
			o.AddSaga(saga)
		}
	}
}

func WithLogger(logger *slog.Logger) NewOrchestratorOpts {
	return func(o *Orchestrator) {
		o.Logger = logger
	}
}

func NewOrchestrator(opts ...NewOrchestratorOpts) *Orchestrator {
	o := &Orchestrator{
		Sagas: make(map[string]*Saga),
	}

	for _, opt := range opts {
		opt(o)
	}

	if o.Logger == nil {
		o.Logger = slog.Default()
	}
	// test
	return o
}

// AddSaga adds a saga to the orchestrator's collection, associating it with its name for subsequent execution management.
func (o *Orchestrator) AddSaga(saga *Saga) {
	o.Sagas[saga.Name] = saga
}

// Run executes all sagas managed by the Orchestrator within the specified context and timeout duration.
// It logs any errors encountered during execution or context timeout and returns the corresponding error if any.
func (o *Orchestrator) Run(ctx context.Context) error {

	var errs error
	for _, saga := range o.Sagas {
		o.Logger.Info("orchestrating saga",
			slog.String("saga_name", saga.Name),
			slog.Int("steps_count", len(saga.steps)),
		)
		if err := saga.executeActions(ctx); err != nil {
			errs = errors.Wrap(err, fmt.Sprintf("issue executing saga %s", saga.Name))
			o.Logger.Info("end orchestrating saga",
				slog.String("saga_name", saga.Name),
				slog.Int("steps_count", len(saga.steps)),
				slog.String("duration", saga.duration().String()),
				slog.String("start_time", saga.startTime.String()),
				slog.String("end_time", saga.endTime.String()),
				slog.String("errors", errs.Error()),
			)
		} else {
			o.Logger.Info("end orchestrating saga",
				slog.String("saga_name", saga.Name),
				slog.Int("steps_count", len(saga.steps)),
				slog.String("duration", saga.duration().String()),
				slog.String("start_time", saga.startTime.String()),
				slog.String("end_time", saga.endTime.String()),
			)
		}

	}

	if errs != nil {
		return errs
	}

	return nil
}
