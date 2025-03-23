package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware para registrar información de las solicitudes HTTP
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Tiempo de inicio
		start := time.Now()

		// Procesar request
		c.Next()

		// Tiempo de finalización
		end := time.Now()
		latency := end.Sub(start)

		// Registrar información de la solicitud
		log.Printf(
			"[%s] %s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
		)
	}
}
