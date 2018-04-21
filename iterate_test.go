package solveuranimateur

import (
	"fmt"
	"testing"

	"github.com/wkschwartz/pigosat"
)

func TestIterateLiteralsCombinations(t *testing.T) {

	iterateLiteralsCombinations(3, []pigosat.Literal{0, 1, 2, 3, 4, 5}, func(combination []pigosat.Literal) {
		fmt.Println(combination)
	})
}
