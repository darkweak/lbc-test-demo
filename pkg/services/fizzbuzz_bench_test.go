package services_test

import (
	"testing"

	"leboncoin/pkg/services"
)

func BenchmarkFizzBuzzCompute(b *testing.B) {
	cases := []struct {
		name  string
		limit int
	}{
		{"small/limit=15", 15},
		{"medium/limit=1000", 1_000},
		{"large/limit=100000", 100_000},
	}

	svc := services.NewFizzBuzz()

	for _, benchCase := range cases {
		b.Run(benchCase.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()

			var result []string

			for b.Loop() {
				result = svc.Compute(3, 5, benchCase.limit, "fizz", "buzz")
			}

			_ = result
		})
	}
}
