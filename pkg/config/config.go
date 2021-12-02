package config

import (
	"errors"
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
	if c.Workload.MaxConcPulls < 1 {
		return errors.New("You can't have parallel value lower than 1")
	}

	if c.Workload.Pulls < 1 {
		return errors.New("You can't have total downloads lower than 1")
	}

	return nil
}
