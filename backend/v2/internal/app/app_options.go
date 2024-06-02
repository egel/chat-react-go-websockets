package app

import "time"

// AppOptions contains options for application server
type AppOptions struct {
	GracefulShutdownTime time.Duration
}

// NewAppOptions returns pointer to new instance of application's options
func NewAppOptions(gracefulShutdownTime time.Duration) *AppOptions {
	return &AppOptions{
		GracefulShutdownTime: gracefulShutdownTime,
	}
}
