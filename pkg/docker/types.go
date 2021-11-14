package docker

type DockerClient struct {
	Username     string
	Password     string
	ExtraHeaders map[string]string
}
