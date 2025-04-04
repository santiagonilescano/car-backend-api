// cmd/api/mediator/mediator.go

package mediator

import (
	"car-service/cmd/api/response"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommandContext mantiene el contexto y las decisiones
type CommandContext struct {
	context.Context
	decisions []string
	mu        sync.Mutex
}

// AddDecision agrega una decisión al contexto
func (c *CommandContext) AddDecision(decision string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.decisions = append(c.decisions, decision)
}

// Definición de la interfaz CommandRequest
type CommandRequest[TResponse any] struct {
	Data TResponse
	// Aquí puedes agregar métodos que sean comunes a todos los CommandRequests
}

type QueryRequest[TResponse any] struct {
	Data TResponse
}

// Actualización de la interfaz CommandHandler
type CommandHandler[T any, R any] interface {
	Execute(request T, ctx *context.Context) (R, error)
}

type QueryHandler[T any, R any] interface {
	Execute(request T, ctx context.Context) (R, error)
}

type Mediator struct {
	handlers map[string]CommandHandler[CommandRequest[any], any]
	queries  map[string]QueryHandler[QueryRequest[any], any]
	mu       sync.RWMutex
	db       *gorm.DB
}

func NewMediator(db *gorm.DB) *Mediator {
	return &Mediator{
		handlers: make(map[string]CommandHandler[CommandRequest[any], any]),
		queries:  make(map[string]QueryHandler[QueryRequest[any], any]),
		db:       db,
	}
}

func (m *Mediator) Register(command string, handler CommandHandler[CommandRequest[any], any]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[command] = handler
}

func (m *Mediator) RegisterQuery(query string, handler QueryHandler[QueryRequest[any], any]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queries[query] = handler
}

// HandleGinRequest maneja una solicitud completa desde Gin
func (m *Mediator) HandleGinRequest(c *gin.Context, commandName string, requestType interface{}) {
	// Crear contexto enriquecido
	cmdCtx := &CommandContext{
		Context:   c.Request.Context(),
		decisions: []string{fmt.Sprintf("Solicitud desde IP: %s", c.ClientIP())},
	}

	// Intentar vincular JSON
	if err := c.ShouldBindJSON(requestType); err != nil {
		response.JSON(c, http.StatusBadRequest, http.StatusBadRequest,
			"Error al procesar la solicitud", nil, []string{err.Error()}, cmdCtx.decisions)
		return
	}

	// Crear solicitud de comando
	request := &CommandRequest[any]{Data: requestType}

	// Ejecutar comando
	result, err := m.executeCommand(commandName, request, cmdCtx)

	// Manejar resultado
	if err != nil {
		response.JSON(c, http.StatusInternalServerError, http.StatusInternalServerError,
			"Error al ejecutar el comando", nil, []string{err.Error()}, cmdCtx.decisions)
		return
	}

	// Respuesta exitosa
	response.JSON(c, http.StatusOK, http.StatusOK,
		"Operación completada con éxito", result, nil, cmdCtx.decisions)
}

// executeCommand busca y ejecuta un comando
func (m *Mediator) executeCommand(commandName string, request *CommandRequest[any], ctx *CommandContext) (any, error) {
	// Buscar handler
	m.mu.RLock()
	handler, exists := m.handlers[commandName]
	m.mu.RUnlock()

	if !exists {
		return nil, errors.New("handler not found")
	}

	// Ejecutar con pipeline
	return m.executeWithPipeline(handler, request, ctx)
}

func (m *Mediator) SendQuery(query string, data *QueryRequest[any]) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	queries, exists := m.queries[query]
	if !exists {
		return nil, errors.New("handler not found")
	}

	resp, err := m.executeWithPipelineQuery(queries, data)
	return resp, err
}

// executeWithPipeline ejecuta un comando dentro de una transacción y recoge decisiones
func (m *Mediator) executeWithPipeline(handler CommandHandler[CommandRequest[any], any],
	data *CommandRequest[any],
	ctx *CommandContext) (any, error) {
	log.Printf("Executing Command: %T", handler)
	tx := m.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.AddDecision(fmt.Sprintf("Recovered from panic: %v", r))
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// Ejecutar el comando pasando el contexto enriquecido
	response, err := handler.Execute(*data, &ctx.Context)
	if err != nil {
		tx.Rollback()
		ctx.AddDecision(fmt.Sprintf("Error ejecutando el comando: %v", err))
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		ctx.AddDecision(fmt.Sprintf("Error al confirmar la transacción: %v", err))
		return nil, err
	}

	return response, nil
}

func (m *Mediator) executeWithPipelineQuery(handler QueryHandler[QueryRequest[any], any], data *QueryRequest[any]) (any, error) {
	// 1. Logging
	log.Printf("Executing query: %T", handler)

	// 2. Exception Handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	reponse, err := handler.Execute(*data, context.Background())
	if err != nil {
		return nil, err
	}

	return reponse, nil
}
