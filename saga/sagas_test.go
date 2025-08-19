package saga

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestSaga(t *testing.T) {

	t.Run("test when saga is executed successfully", func(t *testing.T) {
		saga := NewSaga("executing my saga", slog.Default())

		saga.AddStep("step 1",
			func(ctx context.Context) error {
				fmt.Println("executing step 1")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 1")
				return nil
			},
			true,
		)

		saga.AddStep("step 2",
			func(ctx context.Context) error {
				fmt.Println("executing step 2")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 2")
				return nil
			},
			true,
		)

		err := saga.executeActions(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, SagaStateDone, saga.state)
		for _, step := range saga.steps {
			assert.Equal(t, SagaStepStateDone, step.state)
		}
	})

	t.Run("test when saga step 3 failed and compensating actions", func(t *testing.T) {
		saga := NewSaga("executing my saga", slog.Default())

		saga.AddStep("step 1",
			func(ctx context.Context) error {
				fmt.Println("executing step 1")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 1")
				return nil
			},
			true,
		)

		saga.AddStep("step 2",
			func(ctx context.Context) error {
				fmt.Println("executing step 2")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 2")
				return nil
			},
			true,
		)

		saga.AddStep("step 3",
			func(ctx context.Context) error {
				fmt.Println("executing step 3")
				return errors.New("unable to execute step 3")
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 3")
				return nil
			},
			true,
		)

		saga.AddStep("step 4",
			func(ctx context.Context) error {
				fmt.Println("executing step 4")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 4")
				return nil
			},
			true,
		)

		err := saga.executeActions(context.Background())
		assert.NoError(t, err)

		assert.Equal(t, SagaStateErrored, saga.state)
		assert.Equal(t, SagaStepStateCompensationDone, saga.steps[0].state)
		assert.Equal(t, SagaStepStateCompensationDone, saga.steps[1].state)
		assert.Equal(t, SagaStepStateCompensationDone, saga.steps[2].state)
		assert.Equal(t, SagaStepStateInit, saga.steps[3].state)
	})

	t.Run("test when saga step 3 failed and compensating action for step 2", func(t *testing.T) {
		saga := NewSaga("executing my saga", slog.Default())

		saga.AddStep("step 1",
			func(ctx context.Context) error {
				fmt.Println("executing step 1")
				return nil
			},

			func(ctx context.Context) error {
				fmt.Println("compensating step 1")
				return nil
			},
			true,
		)

		saga.AddStep("step 2",
			func(ctx context.Context) error {
				fmt.Println("executing step 2")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 2")
				return errors.New("unable to execute compensating step 2")
			},
			true,
		)

		saga.AddStep("step 3",
			func(ctx context.Context) error {
				fmt.Println("executing step 3")
				return errors.New("unable to execute step 3")
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 3")
				return nil
			},
			true,
		)

		saga.AddStep("step 4",
			func(ctx context.Context) error {
				fmt.Println("executing step 4")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("compensating step 4")
				return nil
			},
			true,
		)

		err := saga.executeActions(context.Background())
		assert.Error(t, err)

		assert.Equal(t, SagaStatePoisoned, saga.state)
		assert.Equal(t, SagaStepStateCompensationDone, saga.steps[0].state)
		assert.Equal(t, SagaStepStateError, saga.steps[1].state)
		assert.Equal(t, SagaStepStateCompensationDone, saga.steps[2].state)
		assert.Equal(t, SagaStepStateInit, saga.steps[3].state)
	})

	t.Run("test when saga step fails", func(t *testing.T) {
		saga := NewSaga("create a new user",
			slog.Default(),
			WithSteps([]*Step{
				{
					Name: "create_new_user",
					action: func(ctx context.Context) error {
						return errors.New("unable to create users")
					},
					compensation: func(ctx context.Context) error {
						fmt.Println("cleaning up resources")
						return nil
					},
					shouldCompensate: true,
				},
			}))

		err := saga.executeActions(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, SagaStateErrored, saga.state)
		assert.Equal(t, SagaStepStateCompensationDone, saga.steps[0].state)
	})

}
