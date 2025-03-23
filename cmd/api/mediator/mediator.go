// cmd/api/mediator/mediator.go

package mediator

import (
	"context"
	"errors"
	"log"
	"sync"

	"gorm.io/gorm"
)

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
	Execute(request T, ctx context.Context) (R, error)
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

func (m *Mediator) Send(command string, data *CommandRequest[any]) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	handler, exists := m.handlers[command]
	if !exists {
		return nil, errors.New("handler not found")
	}

	// Aquí puedes agregar el pipeline de middlewares
	resp, err := m.executeWithPipeline(handler, data)
	return resp, err
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

func (m *Mediator) executeWithPipeline(handler CommandHandler[CommandRequest[any], any], data *CommandRequest[any]) (any, error) {

	log.Printf("Executing Command: %T", handler)
	tx := m.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	response, err := handler.Execute(*data, context.Background())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
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
