package config

import "fmt"

type Config struct {
	Host  string
	Port  int
	Debug bool
}

func (c *Config) URL() string {
	return fmt.Sprintf("http://%s:%d", c.Host, c.Port)
}
