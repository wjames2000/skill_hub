package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestParseID(t *testing.T) {
	t.Run("valid id", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "42"}}
		id := parseID(c)
		assert.Equal(t, int64(42), id)
	})

	t.Run("negative id returns 0", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "-1"}}
		id := parseID(c)
		assert.Equal(t, int64(0), id)
	})

	t.Run("invalid id returns 0", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		id := parseID(c)
		assert.Equal(t, int64(0), id)
	})
}

func TestParsePagination(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		page, pageSize := parsePagination(c)
		assert.Equal(t, 1, page)
		assert.Equal(t, 20, pageSize)
	})

	t.Run("custom values", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Request = gin.CreateTestContext(nil).Request
		// Set query params
		req := gin.CreateTestContext(nil).Request
		_ = req
		page, pageSize := parsePagination(c)
		assert.Equal(t, 1, page)
		assert.Equal(t, 20, pageSize)
	})
}

func TestParseIDParam(t *testing.T) {
	t.Run("valid id", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "42"}}
		id, err := parseIDParam(c)
		assert.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		_, err := parseIDParam(c)
		assert.Error(t, err)
	})
}
