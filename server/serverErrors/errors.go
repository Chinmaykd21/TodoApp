package serverErrors

import "fmt"

type errMessage struct {
	Code    ErrorCode `json:"errCode"`
	Message string    `json:"errMessage"`
}

func (e *errMessage) Error() string {
	return fmt.Sprintf("Error: %d \nMessage: %s", e.Code, e.Message)
}

func New(code ErrorCode, message string) error {
	e := &errMessage{
		Code:    code,
		Message: message,
	}
	if message == "" {
		e.Message = code.String()
	}
	return e
}
