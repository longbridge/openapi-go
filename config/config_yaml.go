package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type YAMLConfig struct {
}

func (c *YAMLConfig) GetConfig(opts *Options) (*Config, error) {
	parseData := &parseConfig{}
	bytes, err := ioutil.ReadFile(opts.filePath)
	if err != nil {
		err = errors.Wrapf(err, "YAML ReadFile err")
		return nil, err
	}
	err = yaml.Unmarshal(bytes, parseData)
	if err != nil {
		err = errors.Wrapf(err, "YAML GetConfig err")
		return nil, err
	}
	if parseData.Longbridge == nil {
		return nil, errors.New("Longbridge config is not exist in yaml file")
	}
	return parseData.Longbridge, nil
}
