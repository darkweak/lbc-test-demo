package fizzbuzz

import (
	"strconv"
)

type FizzBuzz interface {
	Compute(multiplyFirst, multiplySecond, limit int, fizzStr, buzzStr string) []string
}

var _ FizzBuzz = (*fizzBuzzImpl)(nil)

type fizzBuzzImpl struct{}

func NewFizzBuzz() *fizzBuzzImpl {
	return new(fizzBuzzImpl)
}

func getFizzBuzzResult(current, first, second int, fizz, buzz string) string {
	if (current % first) == 0 {
		if (current % second) == 0 {
			return fizz + buzz
		}

		return fizz
	}

	if (current % second) == 0 {
		return buzz
	}

	return strconv.Itoa(current)
}

func (f fizzBuzzImpl) Compute(multiplyFirst, multiplySecond, limit int, fizzStr, buzzStr string) []string {
	result := make([]string, limit)

	for iter := 1; iter <= limit; iter++ {
		result[iter-1] = getFizzBuzzResult(iter, multiplyFirst, multiplySecond, fizzStr, buzzStr)
	}

	return result
}
