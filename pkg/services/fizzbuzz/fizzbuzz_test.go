package fizzbuzz_test

import (
	"leboncoin/pkg/services/fizzbuzz"
	"testing"
)

const (
	fizzStr    = "fizz"
	buzzStr    = "buzz"
	fizzStrCap = "Fizz"
	buzzStrCap = "Buzz"
)

func TestFizzBuzzComputeClassic(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(3, 5, 15, fizzStr, buzzStr)

	want := []string{
		"1", "2", fizzStr, "4", buzzStr,
		fizzStr, "7", "8", fizzStr, buzzStr,
		"11", fizzStr, "13", "14", fizzStr + buzzStr,
	}

	assertStringSliceEqual(t, got, want)
}

func TestFizzBuzzComputeLimitOne(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	t.Run("not divisible", func(t *testing.T) {
		t.Parallel()

		got := svc.Compute(3, 5, 1, fizzStr, buzzStr)
		assertStringSliceEqual(t, got, []string{"1"})
	})

	t.Run("divisible by first", func(t *testing.T) {
		t.Parallel()

		got := svc.Compute(1, 5, 1, fizzStr, buzzStr)
		assertStringSliceEqual(t, got, []string{fizzStr})
	})

	t.Run("divisible by both", func(t *testing.T) {
		t.Parallel()

		got := svc.Compute(1, 1, 1, "foo", "bar")
		assertStringSliceEqual(t, got, []string{"foobar"})
	})
}

func TestFizzBuzzComputeEmptyStrings(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(3, 3, 3, "", "")
	assertStringSliceEqual(t, got, []string{"1", "2", ""})
}

func TestFizzBuzzComputeCustomMultipliers(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(2, 7, 14, fizzStrCap, buzzStrCap)

	want := []string{
		"1", fizzStrCap, "3", fizzStrCap, "5", fizzStrCap, buzzStrCap,
		fizzStrCap, "9", fizzStrCap, "11", fizzStrCap, "13", fizzStrCap + buzzStrCap,
	}

	assertStringSliceEqual(t, got, want)
}

func TestFizzBuzzComputeResultLength(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	for _, limit := range []int{1, 5, 10, 100} {
		got := svc.Compute(3, 5, limit, fizzStr, buzzStr)
		if len(got) != limit {
			t.Errorf("Compute(limit=%d) returned %d elements, want %d", limit, len(got), limit)
		}
	}
}

func TestFizzBuzzComputeDivisorOne(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(1, 100, 5, fizzStr, buzzStr)

	for index, value := range got {
		if value != fizzStr {
			t.Errorf("Compute()[%d] = %q, want %q (divisor 1 should match every element)", index, value, fizzStr)
		}
	}
}

func TestFizzBuzzComputeDivisibleByFirstOnly(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(3, 5, 3, fizzStr, buzzStr)

	if got[2] != fizzStr {
		t.Errorf("Compute()[2] = %q, want %q", got[2], fizzStr)
	}
}

func TestFizzBuzzComputeDivisibleBySecondOnly(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(3, 5, 5, fizzStr, buzzStr)

	if got[4] != buzzStr {
		t.Errorf("Compute()[4] = %q, want %q", got[4], buzzStr)
	}
}

func TestFizzBuzzComputeDivisibleByBoth(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(3, 5, 15, fizzStr, buzzStr)

	if got[14] != fizzStr+buzzStr {
		t.Errorf("Compute()[14] = %q, want %q", got[14], fizzStr+buzzStr)
	}
}

func TestFizzBuzzComputeLargeNumber(t *testing.T) {
	t.Parallel()

	svc := fizzbuzz.NewFizzBuzz()

	got := svc.Compute(3, 5, 97, fizzStr, buzzStr)

	if got[96] != "97" {
		t.Errorf("Compute()[96] = %q, want %q", got[96], "97")
	}
}

func TestFizzBuzzImplementsInterface(t *testing.T) {
	t.Parallel()

	var _ fizzbuzz.FizzBuzz = fizzbuzz.NewFizzBuzz()
}

func assertStringSliceEqual(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("slice length %d, want %d; got=%v want=%v", len(got), len(want), got, want)
	}

	for index := range want {
		if got[index] != want[index] {
			t.Errorf("[%d] = %q, want %q", index, got[index], want[index])
		}
	}
}
