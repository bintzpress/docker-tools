package setConfig

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func LoadConfig(dir string) (*SetConfig, error) {
	data, err := ioutil.ReadFile(dir + "docker-builder-set.yml")
	if err != nil {
		return nil, err
	}

	var sc SetConfig
	err = yaml.Unmarshal(data, &sc)
	if err == nil {
		err = validate(&sc)
	}
	return &sc, err
}

func validate(sc *SetConfig) error {
	if sc.Version != "1.0" {
		return errors.New("Invalid set config version")
	}
	return nil
}
