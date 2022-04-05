package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

func ReadConfig(path string) (Config, error) {
	var conf Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, fmt.Errorf("reading config file: %w", err)
	}

	err = toml.Unmarshal(data, &conf)
	if err != nil {
		return conf, fmt.Errorf("umarshalling data: %w", err)
	}
	return conf, nil
}
