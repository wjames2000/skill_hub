package errno

var (
	OK              = &Errno{Code: 0, Message: "success"}
	InternalError   = &Errno{Code: 10001, Message: "internal server error"}
	ParamError      = &Errno{Code: 10002, Message: "invalid parameter"}
	NotFound        = &Errno{Code: 10003, Message: "not found"}
	Unauthorized    = &Errno{Code: 10004, Message: "unauthorized"}
	Forbidden       = &Errno{Code: 10005, Message: "forbidden"}
	TooManyRequests = &Errno{Code: 10006, Message: "too many requests"}
	DBError         = &Errno{Code: 20001, Message: "database error"}
	RedisError      = &Errno{Code: 20002, Message: "redis error"}

	UserNotFound    = &Errno{Code: 20101, Message: "user not found"}
	UserExists      = &Errno{Code: 20102, Message: "user already exists"}
	InvalidPassword = &Errno{Code: 20103, Message: "invalid password"}
	InvalidToken    = &Errno{Code: 20104, Message: "invalid token"}
	ExpiredToken    = &Errno{Code: 20105, Message: "expired token"}

	APIKeyNotFound   = &Errno{Code: 20201, Message: "api key not found"}
	AlreadyFavorited = &Errno{Code: 20301, Message: "already favorited"}
	NotFavorited     = &Errno{Code: 20302, Message: "not favorited"}
	DuplicatedReview = &Errno{Code: 20401, Message: "already reviewed"}

	CategoryNotFound = &Errno{Code: 20501, Message: "category not found"}
	CategoryExists   = &Errno{Code: 20502, Message: "category slug already exists"}
)

type Errno struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Errno) Error() string {
	return e.Message
}
