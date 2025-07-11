package handler

import (
	"github.com/loganrk/worker-engine/internal/core/port"
)

type handler struct {
	usecases port.SvrList // List of usecases (services) to handle business logic
	logger   port.Logger  // Logger instance for logging messages
}

// New creates and returns a new handler instance with the provided logger and service list.
func New(loggerIns port.Logger, svcList port.SvrList) *handler {
	return &handler{
		usecases: svcList,   // List of services that will handle specific business logic
		logger:   loggerIns, // Logger for capturing logs
	}
}
