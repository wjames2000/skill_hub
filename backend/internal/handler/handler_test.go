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
		id, err := parseID(c, "id")
		assert.NoError(t, err)
		assert.Equal(t, int64(42), id)
	})

	t.Run("id is zero returns error", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "0"}}
		_, err := parseID(c, "id")
		assert.Error(t, err)
	})

	t.Run("negative id returns error", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "-1"}}
		_, err := parseID(c, "id")
		assert.Error(t, err)
	})

	t.Run("invalid id returns error", func(t *testing.T) {
		c, _ := gin.CreateTestContext(nil)
		c.Params = []gin.Param{{Key: "id", Value: "abc"}}
		_, err := parseID(c, "id")
		assert.Error(t, err)
	})
}
