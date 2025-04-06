// cmd/api/response/response.go
package response

import "github.com/gin-gonic/gin"

type StandardResponse struct {
	Message   string      `json:"message"`
	Errors    []string    `json:"errors,omitempty"`
	Decisions []string    `json:"decisions,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, httpStatus int, message string,
	data interface{}, errors []string, decisions []string) {

	response := StandardResponse{
		Message:   message,
		Errors:    errors,
		Decisions: decisions,
		Data:      data,
	}

	c.JSON(httpStatus, response)
}
