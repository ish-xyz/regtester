package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Load configuration file in the path passed as input and \
// returns a copied value of it.
func Load(cfgPath string) (Config, error) {
	var c Config

	file, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

// Validate configuration against a pre-defined schema
func (c *Config) Validate() error {
	return nil
}
