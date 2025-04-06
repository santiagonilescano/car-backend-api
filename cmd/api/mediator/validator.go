package mediator

import "github.com/gin-gonic/gin"

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type CommandValidator interface {
	Validate(c *gin.Context) []*ValidationError
}
