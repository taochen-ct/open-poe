package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Fail(c *gin.Context, httpCode int, errorCode int, errorMsg string) {
	c.JSON(httpCode, Response{
		ErrorCode: errorCode,
		Data:      nil,
		Message:   errorMsg,
	})
	c.Abort()
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		0,
		data,
		"ok",
	})
	c.Abort()
}

func ServiceError(c *gin.Context, err error) {
	var v *Error
	ok := errors.As(err, &v)
	if ok {
		Fail(c, v.HttpCode(), v.ErrorCode(), v.Error())
	} else {
		Fail(c, http.StatusInternalServerError, ServerError, err.Error())
	}
}

func BusinessError(c *gin.Context, err error) {
	var v *Error
	ok := errors.As(err, &v)
	if ok {
		Fail(c, v.HttpCode(), v.ErrorCode(), v.Error())
	} else {
		Fail(c, http.StatusOK, DefaultError, err.Error())
	}
}

func BadRequestError(c *gin.Context, err error) {
	var v *Error
	ok := errors.As(err, &v)
	if ok {
		Fail(c, v.HttpCode(), v.ErrorCode(), v.Error())
	} else {
		Fail(c, http.StatusBadRequest, ValidateError, err.Error())
	}
}
