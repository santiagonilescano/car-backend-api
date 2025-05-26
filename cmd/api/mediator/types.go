package mediator

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
)

const Query = "query"
const Command = "command"

// ValidationError defines the structure for a validation error.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// CommandHandler defines the interface for command handlers.
// TRequest is the type of the request (command) data.
// TResponse is the type of the response data.
type CommandHandler[TRequest any, TResponse any] interface {
	Execute(request TRequest, ctx *context.Context) (TResponse, error)
	Validate(request TRequest, c *gin.Context, cmdCtx *CommandContext) []*ValidationError
}

// QueryHandler defines the interface for query handlers.
// TRequest is the type of the request (query) data.
// TResponse is the type of the response data.
type QueryHandler[TRequest any, TResponse any] interface {
	Execute(request TRequest, ctx context.Context) (TResponse, error)
}

// CommandContext holds context for command execution, including decisions.
type CommandContext struct {
	context.Context
	decisions []string
	mu        sync.Mutex
}

// AddDecision records a decision made during command processing.
func (c *CommandContext) AddDecision(decision string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.decisions = append(c.decisions, decision)
}
