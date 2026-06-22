package server_test

import (
	"net/http"
	"testing"

	"leboncoin/pkg/configuration"
	"leboncoin/pkg/server"
)

func TestNewServerSetsAddress(t *testing.T) {
	t.Parallel()

	cfg := configuration.Configuration{Address: ":9999"}
	srv := server.NewServer(cfg)

	if srv.Addr != ":9999" {
		t.Errorf("srv.Addr = %q, want %q", srv.Addr, ":9999")
	}
}

func TestNewServerHandlerIsServeMux(t *testing.T) {
	t.Parallel()

	cfg := configuration.Configuration{Address: ":8080"}
	srv := server.NewServer(cfg)

	if _, ok := srv.Handler.(*http.ServeMux); !ok {
		t.Errorf("srv.Handler is %T, want *http.ServeMux", srv.Handler)
	}
}

func TestNewServerReturnsHTTPServer(t *testing.T) {
	t.Parallel()

	cfg := configuration.Configuration{Address: ":8080"}
	srv := server.NewServer(cfg)

	if srv == nil {
		t.Fatal("NewServer() = nil, want *http.Server")
	}
}
