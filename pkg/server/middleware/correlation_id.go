package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const CorrelationIDCtxKey ctxKey = "correlation_id"

func CorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		correlationID := uuid.New().String()

		request.Header.Set("X-Correlation-ID", correlationID)

		next.ServeHTTP(
			writer,
			request.WithContext(context.WithValue(request.Context(), CorrelationIDCtxKey, correlationID)),
		)
	})
}
