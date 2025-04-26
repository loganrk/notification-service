package handler

import (
	"github.com/loganrk/message-service/internal/core/port"
)

type handler struct {
	logger port.Logger // Logger instance for logging messages
}

// New creates and returns a new handler instance with the provided logger.
func New(loggerIns port.Logger) *handler {
	return &handler{
		logger: loggerIns, // Logger for capturing logs
	}
}
