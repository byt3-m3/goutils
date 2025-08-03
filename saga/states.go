package saga

type SagaState int

const (
	SagaStateInit     SagaState = iota
	SagaStateStarted  SagaState = iota
	SagaStateDone     SagaState = iota
	SagaStatePoisoned SagaState = iota
	SagaStateErrored  SagaState = iota
)

type SagaStepState int

const (
	SagaStepStateInit             SagaStepState = iota
	SagaStepStateStarted          SagaStepState = iota
	SagaStepStateCompensating     SagaStepState = iota
	SagaStepStateCompensationDone SagaStepState = iota
	SagaStepStateDone             SagaStepState = iota
	SagaStepStateError            SagaStepState = iota
)
