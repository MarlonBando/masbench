package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
	"masbench/internals/models"
)

var (
	instance *models.Config
	once     sync.Once
)

// GetConfig returns the singleton configuration instance.
// It loads the configuration from a file on the first call.
func GetConfig() *models.Config {
	once.Do(func() {
		loadConfig()
	})
	return instance
}

func loadConfig(filePath ...string) {
	configPath := "masbench_config.yml"
	isDefaultPath := true
	if len(filePath) > 0 && filePath[0] != "" {
		configPath = filePath[0]
		isDefaultPath = false
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) && isDefaultPath {
			fmt.Println("\033[31mNo masbench_config.yml found!\033[0m")
			fmt.Println("Run 'masbench init' to create one or provide the path of an existing configuration file using the appropriate flag.")
			fmt.Println("See the flag with 'masbench --help'")
			os.Exit(1)
		}
		fmt.Printf("\033[31mError reading config file %s: %v\033[0m\n", configPath, err)
		os.Exit(1)
	}

	instance = &models.Config{}
	if err := yaml.Unmarshal(data, instance); err != nil {
		fmt.Printf("\033[31mError parsing config file: %v\033[0m\n", err)
		os.Exit(1)
	}

	// This is necessary for people that already initialize their masbench
	// prior the algo update, this will automatically add the AlgorithmFlagFormat
	// to their configuration
	if instance.AlgorithmFlagFormat == "" {
		defaultFormat := models.DEFAULT_CONFIGURATION.AlgorithmFlagFormat
		instance.AlgorithmFlagFormat = defaultFormat

		if isDefaultPath {
			appendLine := fmt.Sprintf("AlgorithmFlagFormat: %s\n", defaultFormat)
			// Ensure there's a newline before the new field
			if len(data) > 0 && data[len(data)-1] != '\n' {
				appendLine = "\n" + appendLine
			}
			_ = os.WriteFile(configPath, append(data, []byte(appendLine)...), 0644)
		}
	}
}
