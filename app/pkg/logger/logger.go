package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/sohlich/elogrus.v7"
)

type Logger struct {
	*logrus.Logger
}

func New() *Logger {
	return &Logger{
		Logger: logrus.New(),
	}
}

func (log *Logger) AddElasticHook(e *echo.Echo, logname string) error {
	elasticsearchURL := os.Getenv("ELASTICSEARCH_URL")
	client, err := elastic.NewClient(elastic.SetURL(elasticsearchURL))
	if err != nil {
		return errors.Wrap(err, "elastic.NewClient")
	}
	hostname, err := os.Hostname()
	if err != nil {
		return errors.Wrap(err, "os.Hostname")
	}

	hook, err := elogrus.NewAsyncElasticHookWithFunc(client, hostname, logrus.DebugLevel, func() string {
		return logname + "-" + time.Now().Format("2006-01-02")
	})
	if err != nil {
		return errors.Wrap(err, "elogrus.NewAsyncElasticHook")
	}
	log.Hooks.Add(hook)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:      true,
		LogRemoteIP:     true,
		LogHost:         true,
		LogMethod:       true,
		LogURI:          true,
		LogURIPath:      true,
		LogRoutePath:    true,
		LogStatus:       true,
		LogError:        true,
		LogResponseSize: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			entry := log.WithFields(logrus.Fields{
				"latency(ms)":  values.Latency.Milliseconds(),
				"remoteIP":     values.RemoteIP,
				"host":         values.Host,
				"method":       values.Method,
				"URI":          values.URI,
				"URIPath":      values.URIPath,
				"routePath":    values.RoutePath,
				"status":       values.Status,
				"error":        values.Error,
				"responseSize": values.ResponseSize,
			})
			if values.Status < http.StatusInternalServerError {
				entry.Info()
			} else {
				entry.Warn()
			}
			return nil
		},
	}))
	return nil
}
