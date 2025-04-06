package mediator

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
)

const Query = "query"
const Command = "command"

type CommandRequest[TResponse any] struct {
	Data TResponse
}

type QueryRequest[TResponse any] struct {
	Data TResponse
}

type CommandHandler[T any, R any] interface {
	Execute(request T, ctx *context.Context) (R, error)
	Validate(c *gin.Context) []*ValidationError
}

type QueryHandler[T any, R any] interface {
	Execute(request T, ctx context.Context) (R, error)
}

type CommandContext struct {
	context.Context
	decisions []string
	mu        sync.Mutex
}

func (c *CommandContext) AddDecision(decision string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.decisions = append(c.decisions, decision)
}
