package solveuranimateur

import (
	"github.com/wkschwartz/pigosat"
)

func iterateLiteralsCombinations(n int, values []pigosat.Literal, callback func(combination []pigosat.Literal)) {
	iterator := &literalsCombinationsIterator{
		values:      values,
		callback:    callback,
		combination: make([]pigosat.Literal, n),
	}
	iterator.recursion(iterator.values, iterator.combination)
}

type literalsCombinationsIterator struct {
	values, combination []pigosat.Literal

	callback func(combination []pigosat.Literal)
}

func (iterator *literalsCombinationsIterator) recursion(partialValues, partialCombination []pigosat.Literal) {
	if len(partialCombination) == 0 {
		iterator.callback(iterator.combination)
		return
	}
	for i := 0; i < len(partialValues); i++ {
		partialCombination[0] = partialValues[i]
		iterator.recursion(partialValues[i+1:], partialCombination[1:])
	}
}
