package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hpds/skill-hub/pkg/errno"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	t.Run("Success creates correct response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		Success(c, gin.H{"key": "value"})

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "success", resp.Message)
		data := resp.Data.(map[string]interface{})
		assert.Equal(t, "value", data["key"])
	})

	t.Run("Success with nil data", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		Success(c, nil)

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
		assert.Equal(t, "success", resp.Message)
		assert.Nil(t, resp.Data)
	})

	t.Run("Error creates error response", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		Error(c, errno.ParamError)

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10002, resp.Code)
		assert.Equal(t, "parameter error", resp.Message)
	})

	t.Run("ErrorWithMsg uses custom message", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ErrorWithMsg(c, errno.InternalError, "custom detail")

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10001, resp.Code)
		assert.Equal(t, "custom detail", resp.Message)
	})

	t.Run("response has RequestId when set", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("RequestId", "req-abc-123")
		Success(c, nil)

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "req-abc-123", resp.RequestID)
	})
}

func TestResponse_JSON(t *testing.T) {
	t.Run("json marshal roundtrip", func(t *testing.T) {
		orig := &Response{
			Code:    0,
			Message: "test",
			Data:    map[string]interface{}{"count": 42},
		}
		b, err := json.Marshal(orig)
		assert.NoError(t, err)

		var decoded Response
		err = json.Unmarshal(b, &decoded)
		assert.NoError(t, err)
		assert.Equal(t, orig.Code, decoded.Code)
		assert.Equal(t, orig.Message, decoded.Message)
	})
}

func TestResponseWithHTTP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		Success(c, gin.H{"result": "ok"})
	})

	t.Run("http status is 200 on success", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 0, resp.Code)
	})

	r.GET("/error", func(c *gin.Context) {
		Error(c, errno.Unauthorized)
	})

	t.Run("http status is 200 even on error (gin handles)", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)

		var resp Response
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, 10003, resp.Code)
	})
}

func TestContentType(t *testing.T) {
	t.Run("content type is json", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		Success(c, nil)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	})
}
