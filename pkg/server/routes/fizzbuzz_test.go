package routes_test

import (
	"encoding/json"
	"leboncoin/pkg/server/routes"
	"leboncoin/pkg/services/fizzbuzz"
	"leboncoin/pkg/services/pubsub"
	"net/http"
	"net/http/httptest"
	"testing"
)

const stubFizzResult = "fizz"

func newFizzBuzzMux(fb fizzbuzz.FizzBuzz, producer pubsub.Producer) *http.ServeMux {
	mux := http.NewServeMux()
	route := routes.NewFizzBuzz(fb, producer)
	route.Register(mux)

	return mux
}

func TestFizzBuzzQueryParamsHappyPath(t *testing.T) {
	t.Parallel()

	fbStub := &stubFizzBuzz{result: []string{"1", "2", stubFizzResult}}
	producer := &stubProducer{produced: nil, produceErr: nil}
	mux := newFizzBuzzMux(fbStub, producer)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&limit=3&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []string

	err := json.Unmarshal(rec.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	want := []string{"1", "2", stubFizzResult}
	if len(got) != len(want) {
		t.Fatalf("len(body) = %d, want %d", len(got), len(want))
	}

	for index := range want {
		if got[index] != want[index] {
			t.Errorf("body[%d] = %q, want %q", index, got[index], want[index])
		}
	}
}

func TestFizzBuzzQueryParamsMissingInt1(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int2=5&limit=15&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzQueryParamsMissingInt2(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&limit=15&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzQueryParamsMissingLimit(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzQueryParamsInvalidInt1(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=abc&int2=5&limit=15&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzQueryParamsInvalidInt2(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&int2=xyz&limit=15&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzQueryParamsInvalidLimit(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&limit=nope&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzQueryParamsProducesMessage(t *testing.T) {
	t.Parallel()

	fbStub := &stubFizzBuzz{result: []string{"1", "2", stubFizzResult}}
	producer := &stubProducer{produced: nil, produceErr: nil}
	mux := newFizzBuzzMux(fbStub, producer)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&limit=3&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if len(producer.produced) != 1 {
		t.Fatalf("Produce called %d times, want 1", len(producer.produced))
	}

	wantMsg := "3-5-3-fizz-buzz"
	if string(producer.produced[0]) != wantMsg {
		t.Errorf("produced message = %q, want %q", string(producer.produced[0]), wantMsg)
	}
}

func TestFizzBuzzPathParamsHappyPath(t *testing.T) {
	t.Parallel()

	fbStub := &stubFizzBuzz{result: []string{"1", "2", stubFizzResult, "4", "buzz"}}
	producer := &stubProducer{produced: nil, produceErr: nil}
	mux := newFizzBuzzMux(fbStub, producer)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz/3/5/5/fizz/buzz?str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var got []string

	err := json.Unmarshal(rec.Body.Bytes(), &got)
	if err != nil {
		t.Fatalf("response body is not valid JSON: %v; body: %s", err, rec.Body.String())
	}

	if len(got) != 5 {
		t.Errorf("len(body) = %d, want 5", len(got))
	}
}

func TestFizzBuzzPathParamsInvalidInt1(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz/abc/5/15/fizz/buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzPathParamsInvalidInt2(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz/3/xyz/15/fizz/buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzPathParamsInvalidLimit(t *testing.T) {
	t.Parallel()

	mux := newFizzBuzzMux(
		&stubFizzBuzz{result: nil},
		&stubProducer{produced: nil, produceErr: nil},
	)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz/3/5/nope/fizz/buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestFizzBuzzPathParamsStr1Str2FromQueryString(t *testing.T) {
	t.Parallel()

	fbStub := &stubFizzBuzz{result: []string{stubFizzResult}}
	producer := &stubProducer{produced: nil, produceErr: nil}
	mux := newFizzBuzzMux(fbStub, producer)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz/3/5/1/ignored/ignored?str1=foo&str2=bar",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if len(producer.produced) == 0 {
		t.Fatal("Produce was not called")
	}

	wantMsg := "3-5-1-foo-bar"
	if string(producer.produced[0]) != wantMsg {
		t.Errorf("produced message = %q, want %q", string(producer.produced[0]), wantMsg)
	}
}

func TestFizzBuzzPathParamsProducesMessage(t *testing.T) {
	t.Parallel()

	fbStub := &stubFizzBuzz{result: []string{"1"}}
	producer := &stubProducer{produced: nil, produceErr: nil}
	mux := newFizzBuzzMux(fbStub, producer)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz/3/5/1/fizz/buzz?str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if len(producer.produced) != 1 {
		t.Fatalf("Produce called %d times, want 1", len(producer.produced))
	}
}

func TestFizzBuzzQueryParamsProduceErrorReturns500(t *testing.T) {
	t.Parallel()

	fbStub := &stubFizzBuzz{result: []string{"1", "2", "fizz"}}
	producer := &stubProducer{produced: nil, produceErr: errProduceFailed}
	mux := newFizzBuzzMux(fbStub, producer)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&limit=3&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestFizzBuzzPathParamsProduceErrorReturns500(t *testing.T) {
	t.Parallel()

	fbStub := &stubFizzBuzz{result: []string{"1"}}
	producer := &stubProducer{produced: nil, produceErr: errProduceFailed}
	mux := newFizzBuzzMux(fbStub, producer)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz/3/5/1/fizz/buzz?str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

func TestFizzBuzzRegisterQueryParamRoute(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	route := routes.NewFizzBuzz(
		&stubFizzBuzz{result: []string{}},
		&stubProducer{produced: nil, produceErr: nil},
	)
	route.Register(mux)

	req := httptest.NewRequestWithContext(
		t.Context(), http.MethodGet,
		"/fizzbuzz?int1=3&int2=5&limit=0&str1=fizz&str2=buzz",
		nil,
	)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
