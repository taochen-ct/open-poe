package response

import "net/http"

func NewError(httpCode, errorCode int, errorMsg string) *Error {
	return &Error{
		httpCode:  httpCode,
		errorCode: errorCode,
		errorMsg:  errorMsg,
	}
}

func BadRequest(errorMsg string, errorCode ...int) *Error {
	errCode := DefaultError
	if len(errorCode) > 0 {
		errCode = errorCode[0]
	}
	return NewError(http.StatusBadRequest, errCode, errorMsg)
}

func BusinessFail(errorMsg string, errorCode ...int) *Error {
	errCode := DefaultError
	if len(errorCode) > 0 {
		errCode = errorCode[0]
	}
	return NewError(http.StatusOK, errCode, errorMsg)
}

func TooManyRequestsErr(errorMsg string) *Error {
	return NewError(http.StatusTooManyRequests, TooManyRequests, errorMsg)
}

func InternalServer(errorMsg string) *Error {
	return NewError(http.StatusInternalServerError, ServerError, errorMsg)
}

func (e *Error) HttpCode() int {
	return e.httpCode
}

func (e *Error) ErrorCode() int {
	return e.errorCode
}

func (e *Error) Error() string {
	return e.errorMsg
}
