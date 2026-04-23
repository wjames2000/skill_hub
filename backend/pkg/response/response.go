package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/pkg/errno"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, e *errno.Errno) {
	c.JSON(http.StatusOK, Response{
		Code:    e.Code,
		Message: e.Message,
	})
}

func ServerError(c *gin.Context) {
	Error(c, errno.InternalError)
}
