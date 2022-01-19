package templateConfig

import (
	"errors"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadConfig(dir string) (*TemplateConfig, error) {
	var err error
	var data []byte

	_, err = os.Stat(dir + "config.yml")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else {
		data, err = ioutil.ReadFile(dir + "config.yml")
		if err != nil {
			return nil, err
		}

		var tc TemplateConfig
		err = yaml.Unmarshal(data, &tc)
		return &tc, err
	}
}
