package forger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ragarwalll/mta-forge.git/pkg/utils"
	"github.com/sagikazarmark/slog-shim"
	"gopkg.in/yaml.v3"
)

type Forger struct {
	BaseDir   string
	OutputDir string
}

func NewForger(baseDir, outputDir string) *Forger {
	return &Forger{BaseDir: baseDir, OutputDir: outputDir}
}

func (f *Forger) CreateDescriptor(prefix string) (string, error) {
	var resources, modules, config, sharedConfigs map[string]interface{}

	// Read the base.yml
	basePath := filepath.Join(f.BaseDir, prefix, "base.yml")
	slog.Info("Reading base configuration file", "path", basePath)
	err := utils.ReadYAML(basePath, &config)

	if err != nil {
		slog.Error("Failed to read base.yml", "error", err)
		return "", fmt.Errorf("error reading base.yml: %w", err)
	}

	// Read the resources
	resourcesDir := filepath.Join(f.BaseDir, prefix, "resources")
	if _, err = os.Stat(resourcesDir); !os.IsNotExist(err) {
		slog.Debug("Reading resources directory", "path", resourcesDir)
		resources, err = utils.ReadYAMLDir(resourcesDir)

		if err != nil {
			slog.Error("Failed to load resources", "error", err)
			return "", fmt.Errorf("error loading resources: %w", err)
		}
	} else {
		slog.Info("Resources directory not found", "path", resourcesDir)
	}

	// Read the modules
	modulesDir := filepath.Join(f.BaseDir, prefix, "modules")
	if _, err = os.Stat(modulesDir); !os.IsNotExist(err) {
		slog.Debug("Reading modules directory", "path", modulesDir)
		modules, err = utils.ReadYAMLDir(modulesDir)

		if err != nil {
			slog.Error("Failed to load modules", "error", err)
			return "", fmt.Errorf("error loading modules: %w", err)
		}
	} else {
		slog.Info("Modules directory not found", "path", modulesDir)
	}

	// Read the shared configurations
	sharedDir := filepath.Join(f.BaseDir, prefix, "shared")
	if _, err = os.Stat(sharedDir); !os.IsNotExist(err) {
		slog.Debug("Reading shared configurations directory", "path", sharedDir)
		sharedConfigs, err = utils.ReadYAMLDir(sharedDir)

		if err != nil {
			slog.Error("Failed to load shared configs", "error", err)
			return "", fmt.Errorf("error loading shared configs: %w", err)
		}

		// Process modules
		slog.Debug("Processing shared configurations")

		for sharedConfigType, configs := range sharedConfigs {
			configsMap, ok := configs.(map[string]interface{})
			if !ok {
				slog.Warn("Invalid shared config type", "type", sharedConfigType)
				continue
			}

			for configName, config := range configsMap {
				slog.Debug("Applying shared config", "type", sharedConfigType, "name", configName)
				f.applySharedConfig(config, sharedConfigType, modules)
			}
		}
	} else {
		slog.Info("Shared configurations directory not found", "path", sharedDir)
	}

	mta := map[string]interface{}{}

	for key, value := range config {
		mta[key] = value
	}

	mta["modules"] = f.mapToSlice(modules)
	mta["resources"] = f.mapToSlice(resources)

	yamlData, err := yaml.Marshal(mta)
	if err != nil {
		return "", fmt.Errorf("error marshaling MTA: %w", err)
	}

	return string(yamlData), nil
}

func (f *Forger) applySharedConfig(config interface{}, configType string, modules map[string]interface{}) {
	configMap, ok := config.(map[string]interface{})
	if !ok {
		return
	}

	appliesTo, ok := configMap["applies-to"].([]interface{})
	if !ok {
		return
	}

	values, ok := configMap["values"]
	if !ok {
		return
	}

	for _, moduleInterface := range appliesTo {
		module, ok := moduleInterface.(string)
		if !ok {
			continue
		}

		moduleConfig, ok := modules[module].(map[string]interface{})
		if !ok {
			continue
		}

		if _, exists := moduleConfig[configType]; !exists {
			if _, isSlice := values.([]interface{}); isSlice {
				moduleConfig[configType] = []interface{}{}
			} else {
				moduleConfig[configType] = map[string]interface{}{}
			}
		}

		switch existingConfig := moduleConfig[configType].(type) {
		case []interface{}:
			if valuesSlice, ok := values.([]interface{}); ok {
				moduleConfig[configType] = append(existingConfig, valuesSlice...)
			}
		case map[string]interface{}:
			if valuesMap, ok := values.(map[string]interface{}); ok {
				for k, v := range valuesMap {
					existingConfig[k] = v
				}
			}
		}
	}
}

func (f *Forger) mapToSlice(m map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(m))

	for _, v := range m {
		result = append(result, v)
	}

	return result
}
