package routes_test

import (
	"errors"
	"leboncoin/pkg/services/statistics"
)

type stubFizzBuzz struct {
	result []string
}

func (stub *stubFizzBuzz) Compute(_, _ int, _ int, _, _ string) []string {
	return stub.result
}

type stubStatistics struct {
	incrementedKeys []string
	mostRecent      *statistics.Statistic
}

func (stub *stubStatistics) GetMostRecent() *statistics.Statistic {
	return stub.mostRecent
}

func (stub *stubStatistics) Increment(key string) {
	stub.incrementedKeys = append(stub.incrementedKeys, key)
}

type stubProducer struct {
	produced  [][]byte
	produceErr error
}

func (s *stubProducer) Produce(message []byte) error {
	if s.produceErr != nil {
		return s.produceErr
	}

	s.produced = append(s.produced, message)

	return nil
}

func (s *stubProducer) Close() error {
	return nil
}

var errProduceFailed = errors.New("produce failed")
