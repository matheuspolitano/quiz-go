package api

import (
	"github.com/gin-gonic/gin"
)

// Response represents the standard JSON response structure
type Response struct {
	Status  string      `json:"status"`            // "success" or "error"
	Message string      `json:"message,omitempty"` // Success or error message
	Data    interface{} `json:"data,omitempty"`    // Payload for success responses
	Error   string      `json:"error,omitempty"`   // Error details
}

// SendError sends a JSON error response
func SendError(c *gin.Context, message string, errDetail string, statusCode int) {
	response := Response{
		Status:  "error",
		Message: message,
		Error:   errDetail,
	}
	c.JSON(statusCode, response)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
