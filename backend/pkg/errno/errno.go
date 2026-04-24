package errno

var (
	OK                = &Errno{Code: 0, Message: "success"}
	InternalError     = &Errno{Code: 10001, Message: "internal server error"}
	ParamError        = &Errno{Code: 10002, Message: "invalid parameter"}
	NotFound          = &Errno{Code: 10003, Message: "not found"}
	Unauthorized      = &Errno{Code: 10004, Message: "unauthorized"}
	Forbidden         = &Errno{Code: 10005, Message: "forbidden"}
	TooManyRequests   = &Errno{Code: 10006, Message: "too many requests"}
	DBError           = &Errno{Code: 20001, Message: "database error"}
	RedisError        = &Errno{Code: 20002, Message: "redis error"}
)

type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Errno) Error() string {
	return e.Message
}
