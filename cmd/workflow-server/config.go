package main

import (
	"fmt"

	env "github.com/Netflix/go-env"
)

type appConfig struct {
	StateMachineARN  string `env:"STATE_MACHINE_ARN"`
	AwsRegion        string `env:"AWS_REGION"`
	S3ForcePathStyle bool   `env:"S3_FORCE_PATH_STYLE"`
	Endpoint         string `env:"ENDPOINT"`
}

func NewConfig() (*appConfig, error) {
	var cfg appConfig
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
