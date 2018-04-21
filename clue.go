package solveuranimateur

type clue interface {

	// Concerns returns the indexes of the persons concerned by this clue.
	concerns() []PersonIndex

	// Valid returns whether the provided association is legal given the clue.
	valid(map[PersonIndex]NumberValue) bool
}

// Basic clue implementation.
type customClue struct {
	c []PersonIndex
	v func(n map[PersonIndex]NumberValue) bool
}

func (clue *customClue) concerns() []PersonIndex {
	return clue.c
}

func (clue *customClue) valid(n map[PersonIndex]NumberValue) bool {
	return clue.v(n)
}

// The person has this number
type associationClue struct {
	person PersonIndex
	number NumberValue
}

func (a *associationClue) concerns() []PersonIndex {
	return []PersonIndex{a.person}
}

func (a *associationClue) valid(n map[PersonIndex]NumberValue) bool {
	return a.number == n[a.person]
}

// The person does not have this number.
type nonAssociationClue struct {
	person PersonIndex
	number NumberValue
}

func (a *nonAssociationClue) concerns() []PersonIndex {
	return []PersonIndex{a.person}
}

func (a *nonAssociationClue) valid(n map[PersonIndex]NumberValue) bool {
	return a.number != n[a.person]
}
