package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseID(c *gin.Context, param string) (int64, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, errors.New("invalid id: must be positive")
	}
	return id, nil
}
