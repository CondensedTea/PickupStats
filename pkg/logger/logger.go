package logger

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echologrus "github.com/spirosoik/echo-logrus"
)

func SetLogger(e *echo.Echo, lvl string) (*logrus.Logger, error) {
	logger := logrus.New()

	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return logger, err
	}

	logger.SetLevel(level)
	mw := echologrus.NewLoggerMiddleware(logger)
	e.Logger = mw
	e.Use(mw.Hook())
	return logger, nil
}
