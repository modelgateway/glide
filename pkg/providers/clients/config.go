package clients

import "time"

type ClientConfig struct {
	Timeout *time.Duration `yaml:"timeout,omitempty" json:"timeout" swaggertype:"primitive,integer"`
}

func DefaultClientConfig() *ClientConfig {
	defaultTimeout := 10 * time.Second

	return &ClientConfig{
		Timeout: &defaultTimeout,
	}
}
