package services_test

import (
	"fmt"
	"testing"

	"leboncoin/pkg/services"
)

func BenchmarkStatisticsIncrement(b *testing.B) {
	b.ReportAllocs()

	svc := services.NewStatistics()

	svc.Increment("bench-key")
	svc.Increment("bench-key")

	b.ResetTimer()

	for b.Loop() {
		svc.Increment("bench-key")
	}
}

func BenchmarkStatisticsIncrementManyKeys(b *testing.B) {
	const numKeys = 100

	keys := make([]string, numKeys)
	for idx := range numKeys {
		keys[idx] = fmt.Sprintf("key-%d", idx)
	}

	b.ReportAllocs()

	svc := services.NewStatistics()

	b.ResetTimer()

	for b.Loop() {
		svc.Increment(keys[b.N%numKeys])
	}
}

func BenchmarkStatisticsGetMostRecent(b *testing.B) {
	b.ReportAllocs()

	svc := services.NewStatistics()
	svc.Increment("bench-key")
	svc.Increment("bench-key")

	b.ResetTimer()

	var result *services.Statistic

	for b.Loop() {
		result = svc.GetMostRecent()
	}

	_ = result
}

func BenchmarkStatisticsIncrementParallel(b *testing.B) {
	b.ReportAllocs()

	svc := services.NewStatistics()

	svc.Increment("shared-key")
	svc.Increment("shared-key")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			svc.Increment("shared-key")
		}
	})
}

func BenchmarkStatisticsIncrementParallelManyKeys(b *testing.B) {
	const numKeys = 100

	keys := make([]string, numKeys)
	for idx := range numKeys {
		keys[idx] = fmt.Sprintf("key-%d", idx)
	}

	b.ReportAllocs()

	svc := services.NewStatistics()

	b.ResetTimer()

	var counter int

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			svc.Increment(keys[counter%numKeys])

			counter++
		}
	})
}
