package saga

// Step represents a single step in a saga, containing an action, optional compensation, and compensation requirement.
type Step struct {
	Name             string
	action           CTXFunc
	compensation     CTXFunc
	shouldCompensate bool
	state            SagaStepState
}

func (s *Step) SetState(state SagaStepState) {
	s.state = state

}

func NewStep(name string, action CTXFunc, compensation CTXFunc, shouldCompensate bool) *Step {
	return &Step{
		Name:             name,
		action:           action,
		compensation:     compensation,
		shouldCompensate: shouldCompensate,
		state:            SagaStepStateInit,
	}
}
