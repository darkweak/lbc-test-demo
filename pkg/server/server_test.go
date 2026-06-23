package server_test

import (
	"net/http"
	"testing"

	"leboncoin/pkg/configuration"
	"leboncoin/pkg/server"
)

func getDefaultConfiguration(address string) configuration.Configuration {
	return configuration.Configuration{
		Address:        address,
		FTKafka:        false,
		KafkaAddresses: nil,
	}
}

func TestNewServerSetsAddress(t *testing.T) {
	t.Parallel()

	srv := server.NewServer(getDefaultConfiguration(":9999"))

	if srv.Addr != ":9999" {
		t.Errorf("srv.Addr = %q, want %q", srv.Addr, ":9999")
	}
}

func TestNewServerHandlerIsServeMux(t *testing.T) {
	t.Parallel()

	srv := server.NewServer(getDefaultConfiguration(":8080"))

	if _, ok := srv.Handler.(*http.ServeMux); !ok {
		t.Errorf("srv.Handler is %T, want *http.ServeMux", srv.Handler)
	}
}

func TestNewServerReturnsHTTPServer(t *testing.T) {
	t.Parallel()

	srv := server.NewServer(getDefaultConfiguration(":8080"))

	if srv == nil {
		t.Fatal("NewServer() = nil, want *http.Server")
	}
}
