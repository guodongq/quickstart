package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Error struct {
	code    string
	message string
	cause   error
	stack   *stack
}

func New(msg string) *Error {
	return &Error{
		code:    GeneralCode,
		message: msg,
		stack:   newStack(),
	}
}

func Wrap(err error, msg string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		message: msg,
		cause:   err,
		stack:   newStack(),
	}
}

func (e *Error) Error() string {
	if e.cause == nil {
		return e.message
	}
	return e.message + ": " + e.cause.Error()
}

func (e *Error) StackTrace() string {
	return e.stack.frames().format()
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    e.code,
		Message: e.message,
	})
}

func (e *Error) Code() string    { return e.code }
func (e *Error) Message() string { return e.message }
func (e *Error) Cause() error    { return e.cause }

func (e *Error) WithMessagef(format string, v ...interface{}) *Error {
	e.message = fmt.Sprintf(format, v...)
	return e
}

func (e *Error) WithMessage(message string) *Error {
	e.message = message
	return e
}
func (e *Error) WithCode(code string) *Error {
	e.code = code
	return e
}

func (e *Error) WithCause(err error) *Error {
	e.cause = err
	return e
}

func (e *Error) Unwrap() error { return e.cause }

type errorFactory func(error) *Error

func makeErrorFactory(code, message string) errorFactory {
	return func(cause error) *Error {
		return &Error{
			code:    code,
			message: message,
			cause:   cause,
			stack:   newStack(),
		}
	}
}

var (
	NotFoundError           = makeErrorFactory(NotFoundCode, "resource not found")
	ConflictError           = makeErrorFactory(ConflictCode, "resource conflict")
	UnauthorizedError       = makeErrorFactory(UnAuthorizedCode, "unauthorized")
	BadRequestError         = makeErrorFactory(BadRequestCode, "bad request")
	ForbiddenError          = makeErrorFactory(ForbiddenCode, "forbidden")
	MethodNotAllowedError   = makeErrorFactory(MethodNotAllowedCode, "method not allowed")
	PreconditionFailedError = makeErrorFactory(PreconditionCode, "precondition failed")
	RateLimitError          = makeErrorFactory(RateLimitCode, "rate limit exceeded")
	UnknownError            = makeErrorFactory(GeneralCode, "unknown")
)

func makeIsError(code string) func(error) bool {
	return func(err error) bool {
		var e *Error
		return errors.As(err, &e) && e.code == code
	}
}

var (
	IsNotFoundErr        = makeIsError(NotFoundCode)
	IsConflictErr        = makeIsError(ConflictCode)
	IsUnauthorized       = makeIsError(UnAuthorizedCode)
	IsBadRequestErr      = makeIsError(BadRequestCode)
	IsForbiddenErr       = makeIsError(ForbiddenCode)
	IsMethodNotAllowed   = makeIsError(MethodNotAllowedCode)
	IsRateLimitError     = makeIsError(RateLimitCode)
	IsPreconditionFailed = makeIsError(PreconditionCode)
	IsUnknownErr         = makeIsError(GeneralCode)
)

func Cause(err error) error {
	for {
		var e *Error
		ok := errors.As(err, &e)
		if !ok || e.cause == nil {
			break
		}
		err = e.cause
	}
	return err
}

func ErrCode(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) && e.code != "" {
		return e.code
	}
	return GeneralCode
}
