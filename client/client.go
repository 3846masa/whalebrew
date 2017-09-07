package client

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DefaultVersion is the Engine API version used by Whalebrew
const DefaultVersion string = "1.20"

// NewClient returns a Docker client configured for Whalebrew
func NewClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return cli, err
	}
	cli.NegotiateAPIVersionPing(types.Ping{
		APIVersion: DefaultVersion,
	})
	return cli, nil
}
