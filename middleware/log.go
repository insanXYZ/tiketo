package middleware

import (
	"tiketo/util/logger"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logger.Info(makeFieldLog(c), "Incoming request")
		return next(c)
	}
}

func makeFieldLog(c echo.Context) logrus.Fields {
	fields := logrus.Fields{
		"at": time.Now().Format(time.DateTime),
	}

	if c == nil {
		return fields
	}

	fields["method"] = c.Request().Method
	fields["uri"] = c.Request().URL.String()

	return fields
}
