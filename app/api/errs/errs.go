package errs

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func New(code ErrCode, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}

func Newf(code ErrCode, format string, v ...any) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, v...),
	}
}

func (err Error) Error() string {
	return err.Message
}

func IsError(err error) bool {
	var er Error
	return errors.As(err, &er)
}

func GetError(err error) Error {
	var er Error
	if !errors.As(err, &er) {
		return Error{}
	}

	return er
}
