package solveuranimateur

import (
	"github.com/wkschwartz/pigosat"
)

func Solve() ([]Solution, Distribution) {

	p := newProblem(constClues, Numbers, Persons)

	// Compute all the solutions.
	solutions := p.solve()

	// Get the probability distribution.
	distribution := p.getDistribution(solutions)

	return solutions, distribution
}

type problem struct {
	clues []clue

	numbers numbersMap
	persons personsMap

	// Mapping of the persons indexes.
	personLocalToGlobal map[localPersonIndex]PersonIndex
	personGlobalToLocal map[PersonIndex]localPersonIndex
}

func newProblem(clues []clue, numbers numbersMap, persons personsMap) *problem {
	p := &problem{
		clues: clues,

		numbers: numbers,
		persons: persons,

		personLocalToGlobal: make(map[localPersonIndex]PersonIndex),
		personGlobalToLocal: make(map[PersonIndex]localPersonIndex),
	}
	p.mapPersons()
	return p
}

// Solves the problem and returns all the solutions.
func (p *problem) solve() []Solution {
	blocked := p.getBlockedAssociations()

	var formula pigosat.Formula
	formula = append(formula, p.getClausesBlockAssociations(blocked)...)
	formula = append(formula, p.getClausesPersonExactlyOneNumber()...)
	formula = append(formula, p.getClausesNumberAtMostOnePerson()...)

	solutions := p.getAllSolutions(formula)
	p.assertSolutions(solutions)

	return solutions
}

// Associations associates persons to their number value.
type associations map[PersonIndex]NumberIndex

// Computes the blocked associations given the clues.
func (p *problem) getBlockedAssociations() (blocked []associations) {
	for _, clue := range p.clues {
		pers := clue.concerns()
		p.numbers.iteratePermutations(len(pers), func(perm []NumberIndex) {
			vals := make(map[PersonIndex]NumberValue, len(pers))
			for i := 0; i < len(pers); i++ {
				vals[pers[i]] = p.numbers[perm[i]]
			}
			if !clue.valid(vals) {
				asso := make(associations, len(pers))
				for i := 0; i < len(pers); i++ {
					asso[pers[i]] = perm[i]
				}
				// The clue is not valid for this permutation, block it.
				blocked = append(blocked, asso)
			}
		})
	}
	return blocked
}

type localPersonIndex int

// Computes a linear mapping of the persons mentioned in at least one of the clues.
func (p *problem) mapPersons() {
	var localIndex localPersonIndex
	for _, clue := range p.clues {
		for _, globalIndex := range clue.concerns() {
			if _, ok := p.personGlobalToLocal[globalIndex]; !ok {
				// Add to map.
				p.personGlobalToLocal[globalIndex] = localIndex
				p.personLocalToGlobal[localIndex] = globalIndex
				localIndex++
			}
		}
	}
}

// Returns the clauses preventing the given associations.
func (p *problem) getClausesBlockAssociations(blocked []associations) (formula pigosat.Formula) {
	for _, associations := range blocked {
		clause := make(pigosat.Clause, 0, len(associations))
		for personIndex, numberIndex := range associations {
			clause = append(clause, -p.toLiteral(p.personGlobalToLocal[personIndex], numberIndex))
		}
		formula = append(formula, clause)
	}
	return formula
}

// Returns the clauses ensuring a person has exactly one number.
func (p *problem) getClausesPersonExactlyOneNumber() (formula pigosat.Formula) {
	for _, personIndex := range p.personGlobalToLocal {
		literals := make([]pigosat.Literal, len(p.numbers))
		for numberIndex, _ := range p.numbers {
			literals[numberIndex] = p.toLiteral(personIndex, numberIndex)
		}
		formula = append(formula, exactlyOneOf(literals)...)
	}
	return formula
}

// Returns the clauses ensuring a number belongs to at most one person.
func (p *problem) getClausesNumberAtMostOnePerson() (formula pigosat.Formula) {
	for numberIndex, _ := range p.numbers {
		literals := make([]pigosat.Literal, 0, len(p.personGlobalToLocal))
		for _, relevantIndex := range p.personGlobalToLocal {
			literals = append(literals, p.toLiteral(relevantIndex, numberIndex))
		}
		formula = append(formula, atMostOneOf(literals)...)
	}
	return formula
}

func (p *problem) toLiteral(person localPersonIndex, number NumberIndex) pigosat.Literal {
	return pigosat.Literal(1 + int(person)*len(p.numbers) + int(number))
}

func (p *problem) fromLiteral(literal pigosat.Literal) (relevant localPersonIndex, number NumberIndex) {
	relevant = localPersonIndex(int(literal-1) / len(p.numbers))
	number = NumberIndex(int(literal-1) % len(p.numbers))
	return
}

// Runs the SAT solver to compute the solutions given the clauses.
func (p *problem) getAllSolutions(clauses pigosat.Formula) (solutions []Solution) {
	sat, _ := pigosat.New(nil)
	defer sat.Delete()

	sat.Add(clauses)
	for rawSolution, status := sat.Solve(); status == pigosat.Satisfiable; rawSolution, status = sat.Solve() {
		sol := make(Solution)
		for literal, checked := range rawSolution {
			if checked {
				relevantIndex, numberIndex := p.fromLiteral(pigosat.Literal(literal))
				sol[p.personLocalToGlobal[relevantIndex]] = numberIndex
			}
		}
		solutions = append(solutions, sol)
		sat.BlockSolution(rawSolution)
	}
	return solutions
}

// Checks the solutions against the clues. Panic if it fails.
func (p *problem) assertSolutions(solutions []Solution) {
	vals := make(map[PersonIndex]NumberValue, len(p.personGlobalToLocal))
	for _, solution := range solutions {
		for personIndex, numberIndex := range solution {
			vals[personIndex] = p.numbers[numberIndex]
		}
		for _, clue := range p.clues {
			if !clue.valid(vals) {
				panic("invalid solution")
			}
		}
	}
}

type Solution map[PersonIndex]NumberIndex

// Compute the probability distribution over the provided solutions.
func (p problem) getDistribution(solutions []Solution) Distribution {
	distribution := make(Distribution)
	for _, solution := range solutions {
		for personIndex, numberIndex := range solution {
			if distribution[personIndex] == nil {
				distribution[personIndex] = make([]float64, len(p.numbers))
			}
			distribution[personIndex][numberIndex] += 1
		}
	}
	for personIndex, row := range distribution {
		for numberIndex, cell := range row {
			// Normalization over each person row.
			distribution[personIndex][numberIndex] = cell / float64(len(solutions))
		}
	}
	return distribution
}

// Return the formula which specifies that at most one of the provided literals must be true.
func atMostOneOf(literals []pigosat.Literal) pigosat.Formula {
	formula := pigosat.Formula{}
	iterateLiteralsCombinations(2, literals, func(combination []pigosat.Literal) {
		clause := make(pigosat.Clause, 2)
		for idx, literal := range combination {
			clause[idx] = -literal
		}
		formula = append(formula, clause)
	})
	return formula
}

// Return the formula which specifies that exactly one of the provided literals must be true.
func exactlyOneOf(literals []pigosat.Literal) pigosat.Formula {
	return append(atMostOneOf(literals), pigosat.Clause(literals))
}
