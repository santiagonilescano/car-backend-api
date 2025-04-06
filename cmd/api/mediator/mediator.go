// cmd/api/mediator/mediator.go

package mediator

import (
	"car-service/cmd/api/response"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Mediator struct {
	commands map[string]CommandHandler[CommandRequest[any], any]
	queries  map[string]QueryHandler[QueryRequest[any], any]
	mu       sync.RWMutex
	db       *gorm.DB
}

func NewMediator(db *gorm.DB) *Mediator {
	return &Mediator{
		commands: make(map[string]CommandHandler[CommandRequest[any], any]),
		queries:  make(map[string]QueryHandler[QueryRequest[any], any]),
		db:       db,
	}
}

func (m *Mediator) RegisterCommand(command string, handler CommandHandler[CommandRequest[any], any]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.commands[command] = handler
}

func (m *Mediator) RegisterQuery(query string, handler QueryHandler[QueryRequest[any], any]) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queries[query] = handler
}

func (m *Mediator) LogRequest(c *gin.Context, cmdCtx *CommandContext, requestType any) {
	log.Printf("Procesando solicitud: Content-Type=%s, Content-Length=%d, Path=%s", c.GetHeader("Content-Type"), c.Request.ContentLength, c.Request.URL.Path)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
	} else {
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		log.Printf("Cuerpo de la solicitud: %s", string(body))
		if err := json.Unmarshal(body, requestType); err != nil {
			log.Printf("Error al deserializar JSON: %v", err)
			return
		}
	}
}

func (m *Mediator) Validate(c *gin.Context, command CommandHandler[CommandRequest[any], any], cmdCtx *CommandContext) []string {
	if validator, ok := command.(CommandValidator); ok {
		validationErrors := validator.Validate(c)
		if len(validationErrors) > 0 {
			errorMessages := make([]string, len(validationErrors))
			for i, err := range validationErrors {
				errorMessages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
			}
			return errorMessages

		}
		return nil
	}
	return nil
}

func (m *Mediator) Send(c *gin.Context, actionType string, name string, requestType any) {
	cmdCtx := &CommandContext{
		Context:   c.Request.Context(),
		decisions: []string{},
	}

	m.LogRequest(c, cmdCtx, requestType)

	if actionType == "query" {
		m.mu.RLock()
		selectedQuery := m.queries[name]
		m.mu.RUnlock()
		queryRequest := &QueryRequest[any]{Data: requestType}
		result, err := m.ExecuteQuery(selectedQuery, queryRequest)

		if err != nil {
			response.JSON(c, http.StatusInternalServerError, http.StatusInternalServerError,
				"Error al ejecutar query", nil, []string{err.Error()}, cmdCtx.decisions)
			return
		} else {
			response.JSON(c, http.StatusOK, http.StatusOK,
				"Operación completada con éxito", result, nil, cmdCtx.decisions)
			return
		}
	} else {
		m.mu.RLock()
		selectedCommand := m.commands[name]
		m.mu.RUnlock()
		commandRequest := &CommandRequest[any]{Data: requestType}
		validationsErrors := m.Validate(c, selectedCommand, cmdCtx)
		if validationsErrors != nil {
			response.JSON(c, http.StatusBadRequest, http.StatusBadRequest,
				"Bad Request", nil, validationsErrors, cmdCtx.decisions)
			return
		} else {
			result, err := m.ExecuteCommand(selectedCommand, commandRequest, cmdCtx)

			if err != nil {
				response.JSON(c, http.StatusInternalServerError, http.StatusInternalServerError,
					"Error al ejecutar el comando", nil, []string{err.Error()}, cmdCtx.decisions)
				return
			} else {
				response.JSON(c, http.StatusCreated, http.StatusCreated,
					"Operación completada con éxito", result, nil, cmdCtx.decisions)
				return
			}
		}
	}
}

func (m *Mediator) ExecuteQuery(query QueryHandler[QueryRequest[any], any], data *QueryRequest[any]) (any, error) {
	log.Printf("Executing query: %T", query)

	var panicErr error
	defer func() {
		if r := recover(); r != nil {
			panicErr = fmt.Errorf("panic: %v", r)
		}
	}()

	if panicErr != nil {
		return nil, panicErr
	}

	reponse, err := query.Execute(*data, context.Background())
	if err != nil {
		return nil, err
	}

	return reponse, nil
}

func (m *Mediator) ExecuteCommand(command CommandHandler[CommandRequest[any], any], data *CommandRequest[any], ctx *CommandContext) (any, error) {
	log.Printf("Executing Command: %T", command)
	tx := m.db.Begin()

	var panicErr error
	defer func() {
		if r := recover(); r != nil {
			panicErr = fmt.Errorf("panic: %v", r)
			tx.Rollback()
		}
	}()

	if panicErr != nil {
		return nil, panicErr
	}

	response, err := command.Execute(*data, &ctx.Context)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return response, nil
}
