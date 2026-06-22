package routes_test

import "leboncoin/pkg/services"

type stubFizzBuzz struct {
	result []string
}

func (stub *stubFizzBuzz) Compute(_, _ int, _ int, _, _ string) []string {
	return stub.result
}

type stubStatistics struct {
	incrementedKeys []string
	mostRecent      *services.Statistic
}

func (stub *stubStatistics) GetMostRecent() *services.Statistic {
	return stub.mostRecent
}

func (stub *stubStatistics) Increment(key string) {
	stub.incrementedKeys = append(stub.incrementedKeys, key)
}
