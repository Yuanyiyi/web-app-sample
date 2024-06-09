package common

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentrylogrus "github.com/getsentry/sentry-go/logrus"
	"github.com/sirupsen/logrus"
)

var sentryHook *sentrylogrus.Hook

func SentryInit(serverName, sentryDSN string) error {
	if len(sentryDSN) != 0 {
		logger := logrus.StandardLogger()
		logger.SetLevel(logrus.InfoLevel)
		logger.SetOutput(os.Stderr)
		logger.SetFormatter(&logrus.JSONFormatter{})

		// Send only ERROR and higher level logs to Sentry
		sentryLevels := []logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel}

		var err error
		sentryHook, err = sentrylogrus.New(sentryLevels, sentry.ClientOptions{
			Dsn:        sentryDSN,
			ServerName: serverName,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				if hint.Context != nil {
					if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
						// You have access to the original Request
						fmt.Println(req)
					}
				}
				fmt.Println(event)
				return event
			},
			AttachStacktrace: true,
		})
		if err != nil {
			return err
		}
		logger.AddHook(sentryHook)

		// Flushes before calling os.Exit(1) when using logger.Fatal
		// (else all defers are not called, and Sentry does not have time to send the event)
		logrus.RegisterExitHandler(func() { sentryHook.Flush(5 * time.Second) })
	}
	return nil
}
func SentryQuit() {
	if sentryHook != nil {
		sentryHook.Flush(5 * time.Second)
	}
}
