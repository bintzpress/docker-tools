package localConfig

import (
	"bufio"
	"os"
	"strings"
)

type LocalConfig struct {
	Properties map[string]string
}

func NewLocalConfig() *LocalConfig {
	var lc LocalConfig
	lc.Properties = map[string]string{}
	return &lc
}

func LoadConfig(filename string) (*LocalConfig, error) {
	lc := NewLocalConfig()

	if len(filename) == 0 {
		return lc, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				lc.Properties[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lc, nil
}
