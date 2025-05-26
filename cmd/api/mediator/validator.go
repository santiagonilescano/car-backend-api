package mediator

import "github.com/gin-gonic/gin"

// ValidationError is defined in types.go

type CommandValidator interface {
	Validate(c *gin.Context, ctx *CommandContext) []*ValidationError
}
