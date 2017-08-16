package api

import raven "github.com/getsentry/raven-go"

// Sentry DSN for internal error logging.
func init() {
	raven.SetDSN("https://71d8d1bc8e4843eeba979fdaadebe48b:df30d2048fc44b5185809f04ba9d2294@sentry.io/186627")
}

// LogSentry submits private errors to sentry.
func LogSentry(err error) {
	raven.CaptureError(err, nil)
}
