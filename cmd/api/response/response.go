// cmd/api/response/response.go
package response

import "github.com/gin-gonic/gin"

// StandardResponse es la estructura estándar de respuesta
type StandardResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Errors     []string    `json:"errors,omitempty"`
	Decisions  []string    `json:"decisions,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// JSON envía una respuesta JSON estandarizada
func JSON(c *gin.Context, httpStatus int, internalCode int, message string,
	data interface{}, errors []string, decisions []string) {

	response := StandardResponse{
		StatusCode: internalCode,
		Message:    message,
		Errors:     errors,
		Decisions:  decisions,
		Data:       data,
	}

	c.JSON(httpStatus, response)
}
