// Package cedragokit is the Go SDK for the Cedra blockchain network.
// Use api.NewCedra to get started.
package cedragokit

import (
	"github.com/celerfi/cedra-go-kit/api"
	"github.com/celerfi/cedra-go-kit/client"
)

// New returns a Cedra client for the given network.
func New(network client.Network) *api.Cedra {
	return api.NewCedra(network)
}

// NewWithConfig returns a Cedra client with a custom Config.
func NewWithConfig(cfg client.Config) *api.Cedra {
	return api.NewCedraWithConfig(cfg)
}
