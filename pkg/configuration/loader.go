package configuration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Configuration struct {
	Address        string
	FTKafka        bool
	KafkaAddresses []string
}

func NewConfiguration() (Configuration, error) {
	cfg := Configuration{}

	var err error

	cfg.Address, err = loadEnv("ADDRESS", ":8080", false)
	if err != nil {
		return Configuration{}, fmt.Errorf("error loading address: %w", err)
	}

	cfg.FTKafka, err = loadEnv("ENABLE_KAFKA", false, false)
	if err != nil {
		return Configuration{}, fmt.Errorf("error loading kafka feature toggle: %w", err)
	}

	if cfg.FTKafka {
		cfg.KafkaAddresses, err = loadEnv("KAFKA_ADDRESSES", []string{}, true)
		if err != nil {
			return Configuration{}, fmt.Errorf("error loading kafka addresses: %w", err)
		}
	}

	return cfg, nil
}

type EnvVarNotSetError struct {
	env string
}

func (e EnvVarNotSetError) Error() string {
	return "required environment variable not set: " + e.env
}

func NewErrEnvVarNotSet(env string) error {
	return EnvVarNotSetError{
		env: env,
	}
}

func loadEnv[T any](name string, defaultValue T, required bool) (T, error) {
	val, ok := os.LookupEnv(name)
	if !ok {
		if required {
			return defaultValue, NewErrEnvVarNotSet(name)
		}

		return defaultValue, nil
	}

	switch any(defaultValue).(type) {
	case string:
		typedVal, ok := any(val).(T)
		if ok {
			return typedVal, nil
		}
	case bool:
		typedVal, err := strconv.ParseBool(val)
		if err == nil {
			return any(typedVal).(T), nil
		}
	case []string:
		str, ok := any(val).(string)
		if ok {
			return any(strings.Split(str, ",")).(T), nil
		}

		return defaultValue, nil
	}

	return defaultValue, nil
}
