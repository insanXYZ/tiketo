package httpresponse

import (
	"net/http"
	"tiketo/dto"
	"tiketo/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func Success(c echo.Context, message string, data any) error {
	return c.JSON(http.StatusOK, dto.Response{
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, message string, err error) error {
	res := dto.Response{
		Message: message,
	}

	if err != nil {
		if ValidationErrors, ok := err.(validator.ValidationErrors); ok {
			res.Errors = util.GetErrorValidateMessageStruct(ValidationErrors)
		} else {
			res.Message = err.Error()
		}
	}

	return c.JSON(http.StatusBadRequest, res)
}
