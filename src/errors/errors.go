package errors

import "net/http"

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

// Error implements the Error interface
func (e Error) Error() string {
	return e.Err
}

func ErrUnAuthorized() Error {
	return Error{
		Code: http.StatusUnauthorized,
		Err:  "unauthorized request",
	}
}
