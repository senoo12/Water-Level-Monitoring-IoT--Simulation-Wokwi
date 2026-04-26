package response

import "github.com/gin-gonic/gin"

type BaseResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, message string, data interface{}) {
	res := BaseResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	c.JSON(200, res)
}

func Error(c *gin.Context, statusCode int, message string) {
	res := BaseResponse{
		Success: false,
		Message: message,
		Data:    nil,
	}
	c.JSON(statusCode, res)
}