package config

import "fmt"

type Config struct {
	Host string
	Port int
}

func (config Config) GetAddr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}
