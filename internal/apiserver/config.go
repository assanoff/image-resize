package apiserver

// Config ...
type Config struct {
	Port     string `yaml:"port"`
	LogLevel string `yaml:"log_level"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		Port:     ":8080",
		LogLevel: "debug",
	}
}
