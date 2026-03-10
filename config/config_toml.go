package config

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type TOMLConfig struct {
}

func (c *TOMLConfig) GetConfig(opts *Options) (*Config, error) {
	parseData := &parseConfig{}
	_, err := toml.DecodeFile(opts.filePath, parseData)
	if err != nil {
		err = errors.Wrapf(err, "TOML GetConfig err")
		return nil, err
	}
	if parseData.Longbridge == nil {
		return nil, errors.New("Longbridge config is not exist in toml file")
	}
	return parseData.Longbridge, nil
}
