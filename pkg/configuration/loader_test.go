package configuration_test

import (
	"errors"
	"testing"

	"leboncoin/pkg/configuration"
)

func TestNewConfigurationWithAddressEnvSet(t *testing.T) {
	t.Setenv("ADDRESS", ":9090")
	t.Setenv("KAFKA_ADDRESSES", ":9090")

	cfg, err := configuration.NewConfiguration()
	if err != nil {
		t.Fatalf("NewConfiguration() error = %v, want nil", err)
	}

	if cfg.Address != ":9090" {
		t.Errorf("cfg.Address = %q, want %q", cfg.Address, ":9090")
	}
}

func TestNewConfigurationWithAddressEmptyString(t *testing.T) {
	t.Setenv("ADDRESS", "")
	t.Setenv("KAFKA_ADDRESSES", "")

	cfg, err := configuration.NewConfiguration()
	if err != nil {
		t.Fatalf("NewConfiguration() error = %v, want nil", err)
	}

	if cfg.Address != "" {
		t.Errorf("cfg.Address = %q, want empty string when ADDRESS=\"\"", cfg.Address)
	}
}

func TestEnvVarNotSetErrorMessage(t *testing.T) {
	t.Parallel()

	var target configuration.EnvVarNotSetError

	err := configuration.NewErrEnvVarNotSet("MY_VAR")
	if !errors.As(err, &target) {
		t.Fatalf("error is not EnvVarNotSetError: %T %v", err, err)
	}

	want := "required environment variable not set: MY_VAR"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}
