package statistics_test

import (
	"leboncoin/pkg/services/statistics"
	"sync"
	"testing"
	"time"
)

const testKey1 = "key1"

func TestStatisticsGetMostRecentEmpty(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	got := svc.GetMostRecent()
	if got != nil {
		t.Errorf("GetMostRecent() on empty statistics = %v, want nil", got)
	}
}

func TestStatisticsGetMostRecentAfterOnlyOneIncrementIsNil(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()
	svc.Increment(testKey1)

	got := svc.GetMostRecent()
	if got != nil {
		t.Errorf("GetMostRecent() after single Increment = %v, want nil (max_hit not set on first call)", got)
	}
}

func TestStatisticsIncrementSecondCallSetsMaxHit(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	before := time.Now()

	svc.Increment(testKey1)
	svc.Increment(testKey1)

	after := time.Now()

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil after two Increments, want non-nil")
	}

	if got.Key != testKey1 {
		t.Errorf("GetMostRecent().Key = %q, want %q", got.Key, testKey1)
	}

	if got.Hit != 2 {
		t.Errorf("GetMostRecent().Hit = %d, want 2", got.Hit)
	}

	if got.LastCall.Before(before) || got.LastCall.After(after) {
		t.Errorf("GetMostRecent().LastCall = %v, want between %v and %v", got.LastCall, before, after)
	}
}

func TestStatisticsIncrementUpdatesHit(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	svc.Increment(testKey1)
	svc.Increment(testKey1)
	svc.Increment(testKey1)

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil, want non-nil")
	}

	if got.Hit != 3 {
		t.Errorf("GetMostRecent().Hit = %d after 3 increments, want 3", got.Hit)
	}

	if got.Key != testKey1 {
		t.Errorf("GetMostRecent().Key = %q, want %q", got.Key, testKey1)
	}
}

func TestStatisticsIncrementUpdatesLastCall(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	svc.Increment(testKey1)

	before := time.Now()

	svc.Increment(testKey1)

	after := time.Now()

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil, want non-nil")
	}

	if got.LastCall.Before(before) || got.LastCall.After(after) {
		t.Errorf("LastCall not updated: got %v, want between %v and %v", got.LastCall, before, after)
	}
}

func TestStatisticsMostRecentTracksMaxHit(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	svc.Increment(testKey1)
	svc.Increment(testKey1)
	svc.Increment("key2")
	svc.Increment("key2")
	svc.Increment("key2")

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil, want non-nil")
	}

	if got.Key != "key2" {
		t.Errorf("GetMostRecent().Key = %q, want %q (highest hit count)", got.Key, "key2")
	}

	if got.Hit != 3 {
		t.Errorf("GetMostRecent().Hit = %d, want 3", got.Hit)
	}
}

func TestStatisticsMostRecentStaysWhenNewKeyHasFewerHits(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	for range 5 {
		svc.Increment(testKey1)
	}

	svc.Increment("key2")

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil, want non-nil")
	}

	if got.Key != testKey1 {
		t.Errorf("GetMostRecent().Key = %q, want %q", got.Key, testKey1)
	}
}

func TestStatisticsMostRecentSwitchesToNewLeader(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	svc.Increment("keyA")
	svc.Increment("keyA")

	svc.Increment("keyB")
	svc.Increment("keyB")
	svc.Increment("keyB")

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil, want non-nil")
	}

	if got.Key != "keyB" {
		t.Errorf("GetMostRecent().Key = %q after keyB overtakes keyA, want %q", got.Key, "keyB")
	}
}

func TestStatisticsMultipleKeysFirstCallEachReturnsNil(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	for _, key := range []string{"alpha", "beta", "gamma"} {
		svc.Increment(key)
	}

	got := svc.GetMostRecent()
	if got != nil {
		t.Errorf("GetMostRecent() = %v, want nil (no key incremented twice)", got)
	}
}

func TestStatisticsMultipleKeysOneIncrementedTwice(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	svc.Increment("alpha")
	svc.Increment("beta")
	svc.Increment("beta")

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil, want non-nil")
	}

	if got.Key != "beta" {
		t.Errorf("GetMostRecent().Key = %q, want %q", got.Key, "beta")
	}

	if got.Hit != 2 {
		t.Errorf("GetMostRecent().Hit = %d, want 2", got.Hit)
	}
}

func TestStatisticsConcurrentIncrement(t *testing.T) {
	t.Parallel()

	svc := statistics.NewStatistics()

	const (
		goroutines = 50
		iterations = 20
	)

	var waitGroup sync.WaitGroup

	waitGroup.Add(goroutines)

	for range goroutines {
		go func() {
			defer waitGroup.Done()

			for range iterations {
				svc.Increment("shared-key")
			}
		}()
	}

	waitGroup.Wait()

	got := svc.GetMostRecent()
	if got == nil {
		t.Fatal("GetMostRecent() = nil after concurrent increments, want non-nil")
	}

	if got.Key != "shared-key" {
		t.Errorf("GetMostRecent().Key = %q, want %q", got.Key, "shared-key")
	}

	wantHit := goroutines * iterations
	if got.Hit != wantHit {
		t.Errorf("GetMostRecent().Hit = %d after %d concurrent increments, want %d",
			got.Hit, wantHit, wantHit)
	}
}

func TestStatisticsImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ statistics.Statistics = statistics.NewStatistics()
}
