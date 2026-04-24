package errno

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrno(t *testing.T) {
	t.Run("success errno", func(t *testing.T) {
		assert.Equal(t, 0, Success.Code)
		assert.Equal(t, "success", Success.Message)
	})

	t.Run("common errors have correct codes", func(t *testing.T) {
		assert.Equal(t, 10001, InternalError.Code)
		assert.Equal(t, 10002, ParamError.Code)
		assert.Equal(t, 10003, Unauthorized.Code)
		assert.Equal(t, 10004, Forbidden.Code)
		assert.Equal(t, 10005, NotFound.Code)
		assert.Equal(t, 10006, TooManyRequests.Code)
	})

	t.Run("user errors have correct codes", func(t *testing.T) {
		assert.Equal(t, 20101, UserNotFound.Code)
		assert.Equal(t, 20102, UserExists.Code)
		assert.Equal(t, 20103, InvalidPassword.Code)
	})

	t.Run("skill errors have correct codes", func(t *testing.T) {
		assert.Equal(t, 30101, SkillNotFound.Code)
	})

	t.Run("Errno implements error interface", func(t *testing.T) {
		err := InternalError
		assert.Equal(t, "internal server error", err.Error())
	})

	t.Run("Errno with custom message", func(t *testing.T) {
		code := 99999
		msg := "custom error"
		e := &Errno{Code: code, Message: msg}
		assert.Equal(t, code, e.Code)
		assert.Equal(t, msg, e.Message)
		assert.Equal(t, msg, e.Error())
	})

	t.Run("Error method returns message", func(t *testing.T) {
		tests := []struct {
			e    *Errno
			want string
		}{
			{Success, "success"},
			{InternalError, "internal server error"},
			{ParamError, "parameter error"},
			{Unauthorized, "unauthorized"},
			{Forbidden, "forbidden"},
			{NotFound, "not found"},
			{TooManyRequests, "too many requests"},
			{UserNotFound, "user not found"},
			{UserExists, "user already exists"},
			{InvalidPassword, "invalid password"},
			{SkillNotFound, "skill not found"},
		}
		for _, tt := range tests {
			assert.Equal(t, tt.want, tt.e.Error(), "code %d", tt.e.Code)
		}
	})
}
