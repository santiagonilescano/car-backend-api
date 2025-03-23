// internal/infrastructure/api/mediator.go

package api

import (
	"context"
	"errors"
	"log"
	"sync"
)

type Mediator struct {
	handlers map[string]CommandHandler[CommandRequest[any], any]
	mu       sync.RWMutex
}

func NewMediator() *Mediator {
	return &Mediator{
		handlers: make(map[string]CommandHandler[CommandRequest[any], any]),
	}
}

func (m *Mediator) Register(command string, handler CommandHandler[CommandRequest[any], any]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[command] = handler
}

func (m *Mediator) Send(command string, data *CommandRequest[any]) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	handler, exists := m.handlers[command]
	if !exists {
		return errors.New("handler not found")
	}

	// Aquí puedes agregar el pipeline de middlewares
	if err := m.executeWithPipeline(handler, data); err != nil {
		return err
	}

	return nil
}

func (m *Mediator) executeWithPipeline(handler CommandHandler[CommandRequest[any], any], data CommandRequest[any]) error {
	// 1. Logging
	log.Printf("Executing command: %T", handler)

	// 2. Exception Handling
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	// 3. Database Transaction (ejemplo simplificado)
	// Aquí deberías iniciar una transacción y pasarla al handler
	// db.Begin() y luego commit o rollback según el resultado

	_, err := handler.Execute(data, context.Background())
	if err != nil {
		return err
	}

	return nil
}
