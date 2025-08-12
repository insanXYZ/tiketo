package httpresponse

import (
	"net/http"
	"tiketo/dto"
	"tiketo/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func Success(c echo.Context, message string, data any, statusCode ...int) error {
	if len(statusCode) == 0 {
		statusCode = append(statusCode, http.StatusOK)
	}

	return c.JSON(statusCode[0], dto.Response{
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, message string, err error, statusCode ...int) error {
	if len(statusCode) == 0 {
		statusCode = append(statusCode, http.StatusBadRequest)
	}

	res := dto.Response{
		Message: message,
	}

	if err != nil {
		if ValidationErrors, ok := err.(validator.ValidationErrors); ok {
			res.Errors = util.GetErrorValidateMessageStruct(ValidationErrors)
		} else {
			res.Errors = err.Error()
		}
	}

	return c.JSON(statusCode[0], res)
}
