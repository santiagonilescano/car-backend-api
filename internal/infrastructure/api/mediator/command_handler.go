// internal/infrastructure/api/command_handler.go

package api

import "context"

// Definición de la interfaz CommandRequest
type CommandRequest[TResponse any] interface {
	// Aquí puedes agregar métodos que sean comunes a todos los CommandRequests
}

// Actualización de la interfaz CommandHandler
type CommandHandler[TCommandRequest CommandRequest[TResponse], TResponse any] interface {
	Execute(request TCommandRequest, ctx context.Context) (TResponse, error)
}
