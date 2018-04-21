package solveuranimateur

type NumberIndex int
type NumberValue int

type numbersMap map[NumberIndex]NumberValue

func (values numbersMap) Indexes() []NumberIndex {
	indexes := make([]NumberIndex, len(values))
	for index := range indexes {
		indexes[index] = NumberIndex(index)
	}
	return indexes
}

func (values numbersMap) iteratePermutations(n int, callback func(permutation []NumberIndex)) {
	iterator := &numbersPermutationsIterator{
		permuation: make([]NumberIndex, n),
		callback:   callback,
	}
	iterator.recursion(values.Indexes(), iterator.permuation)
}

type numbersPermutationsIterator struct {
	permuation []NumberIndex

	callback func(permutation []NumberIndex)
}

func (iterator *numbersPermutationsIterator) recursion(partialIndexes, partialPermutation []NumberIndex) {
	if len(partialPermutation) == 0 {

		iterator.callback(iterator.permuation)
		return
	}
	for i := 0; i < len(partialIndexes); i++ {
		partialPermutation[0] = partialIndexes[i]
		iterator.swap(partialIndexes, i)
		iterator.recursion(partialIndexes[1:], partialPermutation[1:])
		iterator.swap(partialIndexes, i)
	}
}

func (iterator *numbersPermutationsIterator) swap(partialIndexes []NumberIndex, i int) {
	tmp := partialIndexes[i]
	partialIndexes[i] = partialIndexes[0]
	partialIndexes[0] = tmp
}
