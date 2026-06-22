package routes_test

import (
	"encoding/json"
	"leboncoin/pkg/services/statistics"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"leboncoin/pkg/server/routes"
)

func newStatisticsMux(stats statistics.Statistics) *http.ServeMux {
	mux := http.NewServeMux()
	route := routes.NewStatistics(stats)
	route.Register(mux)

	return mux
}

func TestStatisticsHandlerNoData(t *testing.T) {
	t.Parallel()

	mux := newStatisticsMux(&stubStatistics{incrementedKeys: nil, mostRecent: nil})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/statistics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestStatisticsHandlerWithData(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	stat := &statistics.Statistic{
		LastCall: now,
		Hit:      42,
		Key:      "3-5-15-fizz-buzz",
	}

	mux := newStatisticsMux(&stubStatistics{incrementedKeys: nil, mostRecent: stat})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/statistics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got statistics.Statistic

	err := json.Unmarshal(rec.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if got.Key != stat.Key {
		t.Errorf("body.Key = %q, want %q", got.Key, stat.Key)
	}

	if got.Hit != stat.Hit {
		t.Errorf("body.Hit = %d, want %d", got.Hit, stat.Hit)
	}
}

func TestStatisticsHandlerResponseIsValidJSON(t *testing.T) {
	t.Parallel()

	stat := &statistics.Statistic{
		LastCall: time.Now(),
		Hit:      1,
		Key:      "some-key",
	}

	mux := newStatisticsMux(&stubStatistics{incrementedKeys: nil, mostRecent: stat})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/statistics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var raw map[string]any

	err := json.Unmarshal(rec.Body.Bytes(), &raw)
	if err != nil {
		t.Errorf("response is not valid JSON: %v; body: %s", err, rec.Body.String())
	}
}

func TestStatisticsHandlerJSONFieldNames(t *testing.T) {
	t.Parallel()

	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	stat := &statistics.Statistic{
		LastCall: now,
		Hit:      7,
		Key:      "1-2-10-foo-bar",
	}

	mux := newStatisticsMux(&stubStatistics{incrementedKeys: nil, mostRecent: stat})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/statistics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var raw map[string]any

	err := json.Unmarshal(rec.Body.Bytes(), &raw)
	if err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	for _, field := range []string{"last_call", "hit", "key"} {
		if _, ok := raw[field]; !ok {
			t.Errorf("JSON response missing field %q", field)
		}
	}
}

func TestStatisticsRegisterRoute(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	route := routes.NewStatistics(&stubStatistics{incrementedKeys: nil, mostRecent: nil})
	route.Register(mux)

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/statistics", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d (route should be registered)", rec.Code, http.StatusNotFound)
	}
}

func TestStatisticsImplementsRoute(t *testing.T) {
	t.Parallel()

	var _ routes.Route = routes.NewStatistics(&stubStatistics{incrementedKeys: nil, mostRecent: nil})
}
