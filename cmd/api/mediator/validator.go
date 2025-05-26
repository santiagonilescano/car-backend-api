package mediator

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type CommandValidator interface {
	Validate(request CommandRequest[any], ctx *CommandContext) []*ValidationError
}
