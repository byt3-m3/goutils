package saga

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestOrchestrator(t *testing.T) {

	t.Run("test when orchestrator is successful", func(t *testing.T) {

		testSaga := NewSaga("test saga 1", slog.Default(),
			WithSteps([]*Step{
				NewStep("step 1",
					func(ctx context.Context) error {
						return nil
					},
					func(ctx context.Context) error {
						return nil
					},
					true,
				),
			}),
		)

		orchestrator := NewOrchestrator(WithSagas([]*Saga{testSaga}))

		orchestrator.AddSaga(testSaga)

		ctx := context.Background()

		err := orchestrator.Run(ctx)
		assert.NoError(t, err)
		assert.Equal(t, SagaStateDone, testSaga.state)
	})

	t.Run("test when saga step action fails", func(t *testing.T) {

		testSaga := NewSaga("test saga 1", slog.Default(),
			WithSteps([]*Step{
				{
					Name: "step 1",
					action: func(ctx context.Context) error {
						return nil
					},
					compensation: func(ctx context.Context) error {
						return nil
					},
					shouldCompensate: true,
				},
				{
					Name: "step 2",
					action: func(ctx context.Context) error {
						return errors.New("unable to execute step 2")
					},
					compensation: func(ctx context.Context) error {
						return nil
					},
					shouldCompensate: true,
				},
			}),
		)

		orchestrator := NewOrchestrator(WithSagas([]*Saga{testSaga}))

		orchestrator.AddSaga(testSaga)

		ctx := context.Background()

		err := orchestrator.Run(ctx)
		assert.NoError(t, err)

	})

	t.Run("test when saga step action and compensation fails", func(t *testing.T) {

		testSaga := NewSaga("test saga 1", slog.Default(),
			WithSteps([]*Step{
				{
					Name: "step 1",
					action: func(ctx context.Context) error {
						slog.Info("step 1")
						return nil
					},
					compensation: func(ctx context.Context) error {
						slog.Info("step 1 compensation")
						return nil
					},
					shouldCompensate: true,
				},
				{
					Name: "step 2",
					action: func(ctx context.Context) error {
						slog.Info("step 2")
						return nil
					},
					compensation: func(ctx context.Context) error {
						slog.Info("step 2 compensation")
						return errors.New("unable to execute step 2 compensation")
					},
					shouldCompensate: true,
				},
				{
					Name: "step 3",
					action: func(ctx context.Context) error {
						slog.Info("step 3")
						return errors.New("error executing step 3")
					},
					compensation: func(ctx context.Context) error {
						slog.Info("step 3 compensation")
						return nil
					},
					shouldCompensate: true,
				},
			}),
		)

		orchestrator := NewOrchestrator(WithSagas([]*Saga{testSaga}))

		ctx := context.Background()

		err := orchestrator.Run(ctx)
		assert.Error(t, err)
	})

}
