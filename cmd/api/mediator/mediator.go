// cmd/api/mediator/mediator.go

package mediator

import (
	"car-service/cmd/api/response"
	"car-service/internal/domain/errors"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"reflect"

	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Mediator struct {
	commands map[string]any // Store any type of command handler
	queries  map[string]any // Store any type of query handler
	mu       sync.RWMutex
	db       *gorm.DB
}

func NewMediator(db *gorm.DB) *Mediator {
	return &Mediator{
		commands: make(map[string]any),
		queries:  make(map[string]any),
		db:       db,
	}
}

// RegisterCommand stores the handler with its type information erased.
func (m *Mediator) RegisterCommand(commandName string, handler any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.commands[commandName] = handler
}

// RegisterQuery stores the handler with its type information erased.
func (m *Mediator) RegisterQuery(queryName string, handler any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.queries[queryName] = handler
}

func (m *Mediator) LogRequest(c *gin.Context, body []byte) {
	log.Printf("Procesando solicitud: Method=%s, Path=%s, Content-Type=%s, Content-Length=%d",
		c.Request.Method, c.Request.URL.Path, c.GetHeader("Content-Type"), c.Request.ContentLength)
	if len(body) > 0 {
		log.Printf("Cuerpo de la solicitud: %s", string(body))
	} else {
		log.Printf("Cuerpo de la solicitud: [vacío o no aplicable]")
	}
}

func (m *Mediator) Send(c *gin.Context, actionType string, name string, requestModel any) {
	cmdCtx := &CommandContext{
		Context:   c.Request.Context(),
		decisions: []string{},
	}

	var bodyBytes []byte
	var err error

	// Read body only for relevant methods
	if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, err = io.ReadAll(c.Request.Body)
			if err != nil {
				log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
				response.JSON(c, http.StatusInternalServerError, "Error al leer el cuerpo de la solicitud", nil, []string{err.Error()}, cmdCtx.decisions)
				return
			}
			// Restore the body so it can be read again if needed.
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
	}

	m.LogRequest(c, bodyBytes)

	// Unmarshal body into requestModel only if bodyBytes is not empty and requestModel is a pointer.
	// requestModel must be a pointer for json.Unmarshal to populate it.
	if len(bodyBytes) > 0 {
		if reflect.ValueOf(requestModel).Kind() != reflect.Ptr {
			log.Printf("Error: requestModel debe ser un puntero para deserializar el cuerpo. Tipo recibido: %T", requestModel)
			response.JSON(c, http.StatusInternalServerError, "Error interno del servidor", nil, []string{"requestModel no es un puntero"}, cmdCtx.decisions)
			return
		}
		if err := json.Unmarshal(bodyBytes, requestModel); err != nil {
			log.Printf("Error al deserializar JSON: %v. Cuerpo: %s", err, string(bodyBytes))
			response.JSON(c, http.StatusBadRequest, "Error al procesar la solicitud: formato JSON inválido", nil, []string{err.Error()}, cmdCtx.decisions)
			return
		}
	}


	if actionType == Query {
		m.mu.RLock()
		handler := m.queries[name]
		m.mu.RUnlock()

		if handler == nil {
			response.JSON(c, http.StatusNotFound, "Query handler not found", nil, nil, cmdCtx.decisions)
			return
		}
		
		val, err := m.executeQueryReflect(handler, requestModel, cmdCtx.Context)
		if err != nil {
			response.JSON(c, http.StatusInternalServerError, "Error al ejecutar query", nil, []string{err.Error()}, cmdCtx.decisions)
			return
		}
		response.JSON(c, http.StatusOK, "Operación completada con éxito", val, nil, cmdCtx.decisions)

	} else { // Command
		m.mu.RLock()
		handler := m.commands[name]
		m.mu.RUnlock()

		if handler == nil {
			response.JSON(c, http.StatusNotFound, "Command handler not found", nil, nil, cmdCtx.decisions)
			return
		}

		validationMessages := m.validateCommandReflect(handler, requestModel, c, cmdCtx)
		if len(validationMessages) > 0 {
			response.JSON(c, http.StatusBadRequest, "Bad Request", nil, validationMessages, cmdCtx.decisions)
			return
		}

		result, err := m.executeCommandReflect(handler, requestModel, cmdCtx)
		if err != nil {
			if businessErr, ok := err.(*errors.BusinessError); ok {
				response.JSON(c, http.StatusConflict, "Error de negocio", nil, []string{businessErr.Message}, cmdCtx.decisions)
			} else {
				response.JSON(c, http.StatusInternalServerError, "Error al ejecutar el comando", nil, []string{err.Error()}, cmdCtx.decisions)
			}
			return
		}
		response.JSON(c, http.StatusCreated, "Operación completada con éxito", result, nil, cmdCtx.decisions)
	}
}

func (m *Mediator) executeQueryReflect(handler any, request any, ctx context.Context) (any, error) {
	handlerValue := reflect.ValueOf(handler)
	requestValue := reflect.ValueOf(request)

	// If requestModel was passed as a pointer for unmarshalling, but the handler expects a value, dereference it.
	if requestValue.Kind() == reflect.Ptr && handlerValue.Type().Method(0).Type.In(0).Kind() != reflect.Ptr {
		requestValue = requestValue.Elem()
	}


	executeMethod := handlerValue.MethodByName("Execute")
	if !executeMethod.IsValid() {
		return nil, fmt.Errorf("handler %T no tiene el método Execute", handler)
	}

	returnValues := executeMethod.Call([]reflect.Value{requestValue, reflect.ValueOf(ctx)})
	if len(returnValues) != 2 {
		return nil, fmt.Errorf("el método Execute no retornó 2 valores")
	}

	result := returnValues[0].Interface()
	errVal := returnValues[1].Interface()
	if errVal != nil {
		return result, errVal.(error)
	}
	return result, nil
}

func (m *Mediator) executeCommandReflect(handler any, request any, cmdCtx *CommandContext) (any, error) {
	tx := m.db.Begin()
	// Add transaction to context for handlers to use
	ctxWithTx := context.WithValue(cmdCtx.Context, "DB_TX", tx)
	cmdCtx.Context = ctxWithTx // Update CommandContext's context

	var panicErr error
	defer func() {
		if r := recover(); r != nil {
			panicErr = fmt.Errorf("panic: %v", r)
			tx.Rollback()
		}
	}()

	handlerValue := reflect.ValueOf(handler)
	requestValue := reflect.ValueOf(request)
	
	// If requestModel was passed as a pointer for unmarshalling, but the handler expects a value, dereference it.
	if requestValue.Kind() == reflect.Ptr && handlerValue.Type().Method(0).Type.In(0).Kind() != reflect.Ptr {
		 requestValue = requestValue.Elem()
	}

	// cmdCtxValue := reflect.ValueOf(cmdCtx) // This was unused

	executeMethod := handlerValue.MethodByName("Execute")
	if !executeMethod.IsValid() {
		tx.Rollback()
		return nil, fmt.Errorf("handler %T no tiene el método Execute", handler)
	}
	
	// The Execute method of a CommandHandler expects (request TRequest, ctx *context.Context)
	// So we pass requestValue and reflect.ValueOf(&cmdCtx.Context)
	returnValues := executeMethod.Call([]reflect.Value{requestValue, reflect.ValueOf(&cmdCtx.Context)})
	if len(returnValues) != 2 {
		tx.Rollback()
		return nil, fmt.Errorf("el método Execute no retornó 2 valores")
	}

	result := returnValues[0].Interface()
	errVal := returnValues[1].Interface()

	if errVal != nil {
		tx.Rollback()
		return result, errVal.(error)
	}
	
	if panicErr != nil {
		tx.Rollback() // Ensure rollback on panic
		return nil, panicErr
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Mediator) validateCommandReflect(handler any, request any, c *gin.Context, cmdCtx *CommandContext) []string {
	handlerValue := reflect.ValueOf(handler)
	validateMethod := handlerValue.MethodByName("Validate")

	if !validateMethod.IsValid() {
		return nil // No Validate method, assume valid
	}

	requestValue := reflect.ValueOf(request)
	// If requestModel was passed as a pointer for unmarshalling, but the handler expects a value, dereference it.
	if requestValue.Kind() == reflect.Ptr && validateMethod.Type().In(0).Kind() != reflect.Ptr {
		 requestValue = requestValue.Elem()
	}
	
	ginContextValue := reflect.ValueOf(c)
	cmdCtxValue := reflect.ValueOf(cmdCtx)

	// Validate method expects (request TRequest, c *gin.Context, ctx *CommandContext)
	returnValues := validateMethod.Call([]reflect.Value{requestValue, ginContextValue, cmdCtxValue})
	if len(returnValues) == 1 && !returnValues[0].IsNil() {
		if errors, ok := returnValues[0].Interface().([]*ValidationError); ok {
			errorMessages := make([]string, len(errors))
			for i, err := range errors {
				errorMessages[i] = fmt.Sprintf("%s: %s", err.Field, err.Message)
			}
			return errorMessages
		}
	}
	return nil
}
