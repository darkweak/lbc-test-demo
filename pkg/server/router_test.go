package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"leboncoin/pkg/server"
)

type routeStub struct {
	pattern string
	called  bool
}

func (routeStub *routeStub) Register(mux *http.ServeMux) {
	routeStub.called = true
	mux.HandleFunc(routeStub.pattern, func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})
}

func TestNewRouterRegistersAllRoutes(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	routeA := &routeStub{pattern: "/a", called: false}
	routeB := &routeStub{pattern: "/b", called: false}

	server.NewRouter(mux, routeA, routeB)

	if !routeA.called {
		t.Error("routeA.Register was not called")
	}

	if !routeB.called {
		t.Error("routeB.Register was not called")
	}
}

func TestNewRouterRoutesAreReachable(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	server.NewRouter(mux, &routeStub{pattern: "/ping", called: false})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}

func TestNewRouterNoRoutesIsNoop(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	server.NewRouter(mux)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
