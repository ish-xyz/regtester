package docker

import (
	"github.com/ish-xyz/regtester/pkg/config"
	"github.com/ish-xyz/regtester/pkg/controller"
)

func NewDockerClient(c config.Config) *DockerClient {
	return &DockerClient{
		Username:     c.Connection.BasicAuth.Username,
		Password:     c.Connection.BasicAuth.Password,
		ExtraHeaders: c.Connection.ExtraHeaders,
	}
}

func Pull() *controller.Report {
	return nil
}
