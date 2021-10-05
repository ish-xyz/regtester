package config

type Config struct {
	Connection ConnectionCfg `yaml:"connection"`
	Registries []string      `yaml:"registries"`
	Images     []string      `yaml:"images"`
	Spec       SpecCfg       `yaml:"spec"`
}

type ConnectionCfg struct {
	BasicAuth    BasicAuthCfg      `yaml:"basicAuth"`
	CAPath       string            `yaml:"CAPath"`
	ExtraHeaders map[string]string `yaml:"extraHeaders"`
}

type BasicAuthCfg struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type SpecCfg struct {
	Pulls         int `yaml:"pulls"`
	ParallelPulls int `yaml:"parallelPulls"`
}
