package services

import "fmt"

type FizzBuzz interface {
	Compute(multiplyFirst, multiplySecond, limit int, fizzStr, buzzStr string) []string
}

type fizzBuzzImpl struct{}

func NewFizzBuzz() FizzBuzz {
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

	return fmt.Sprintf("%d", current)
}

func (f fizzBuzzImpl) Compute(multiplyFirst, multiplySecond, limit int, fizzStr, buzzStr string) []string {
	result := make([]string, limit)

	for iter := 1; iter <= limit; iter++ {
		result[iter-1] = getFizzBuzzResult(iter, multiplyFirst, multiplySecond, fizzStr, buzzStr)
	}

	return result
}
