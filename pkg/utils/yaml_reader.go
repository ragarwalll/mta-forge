package utils

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadYAML reads a YAML file from the given path and unmarshals it into the provided interface
func ReadYAML(path string, out interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, out)
}

// ReadYAMLString reads a YAML file from the given path and returns it as a string
func ReadYAMLString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// ReadYAMLDir reads all .yml files in a directory and returns a map
// with filename (without extension) as key and parsed YAML content as value
func ReadYAMLDir(dirPath string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml")) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var data interface{}
			err = yaml.Unmarshal(content, &data)

			if err != nil {
				return err
			}

			filename := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			result[filename] = data
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
