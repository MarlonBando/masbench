package models

// Config holds the application configuration settings.
type Config struct {
	ServerPath          string `yaml:"ServerPath"`
	LevelsDir           string `yaml:"LevelsDir"`
	BenchmarkFolder     string `yaml:"BenchmarkFolder"`
	ClientCommand       string `yaml:"ClientCommand"`
	Timeout             int    `yaml:"Timeout"`
	AlgorithmFlagFormat string `yaml:"AlgorithmFlagFormat"`
}

var DefaultConfiguration Config = Config{
	ServerPath:          "path/to/your/server_executable",
	LevelsDir:           "path/to/your/levels_directory",
	BenchmarkFolder:     "benchmarks",
	ClientCommand:       "your_client_command --level {level_path}",
	Timeout:             180,
	AlgorithmFlagFormat: "-%s",
}
