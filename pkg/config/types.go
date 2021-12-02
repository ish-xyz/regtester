package config

type Config struct {
	Connection ConnectionCfg `yaml:"connection"`
	Registries []string      `yaml:"registries"`
	Images     []string      `yaml:"images"`
	Workload   WorkloadCfg   `yaml:"workload"`
}

type ConnectionCfg struct {
	BasicAuth       BasicAuthCfg      `yaml:"basicAuth"`
	CAPath          string            `yaml:"CAPath"`
	ExtraHeaders    map[string]string `yaml:"extraHeaders"`
	ManifestVersion string            `yaml:"manifestVersion"`
}

type BasicAuthCfg struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type WorkloadCfg struct {
	Pulls          int  `yaml:"pulls"`
	MaxConcPulls   int  `yaml:"maxConcPulls"`
	MaxConcLayers  int  `yaml:"maxConcLayers"`
	CheckIntegrity bool `yaml:"checkIntegrity"`
}
