package docker

type DockerClient struct {
	Username        string
	Password        string
	ExtraHeaders    map[string]string
	ManifestVersion string
	CheckIntegrity  bool
	MaxConcLayers   int
}

type ImageManifestV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	}
	Layers []Layer
}

type Layer struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}
