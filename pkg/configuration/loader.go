package configuration

import (
	"fmt"
	"os"
)

type Configuration struct {
	Address string
}

func NewConfiguration() (Configuration, error) {
	address, err := loadEnv("ADDRESS", ":8080", false)
	if err != nil {
		return Configuration{}, fmt.Errorf("error loading address: %w", err)
	}

	return Configuration{
		Address: address,
	}, nil
}

func loadEnv[T any](name string, defaultValue T, required bool) (T, error) {
	v, ok := os.LookupEnv(name)
	if !ok {
		if required {
			return defaultValue, fmt.Errorf("required environment variable not set: %s", name)
		}

		return defaultValue, nil
	}

	switch any(defaultValue).(type) {
	case string:
		return any(v).(T), nil
	}

	return defaultValue, nil
}
