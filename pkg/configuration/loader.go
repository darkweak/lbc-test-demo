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

type EnvVarNotSetError struct {
	env string
}

func (e EnvVarNotSetError) Error() string {
	return "required environment variable not set: " + e.env
}

func newErrEnvVarNotSet(env string) error {
	return EnvVarNotSetError{
		env: env,
	}
}

func loadEnv[T any](name string, defaultValue T, required bool) (T, error) {
	val, ok := os.LookupEnv(name)
	if !ok {
		if required {
			return defaultValue, newErrEnvVarNotSet(name)
		}

		return defaultValue, nil
	}

	if fmt.Sprintf("%T", defaultValue) == fmt.Sprintf("%T", name) {
		typedVal, ok := any(val).(T)
		if ok {
			return typedVal, nil
		}
	}

	return defaultValue, nil
}
