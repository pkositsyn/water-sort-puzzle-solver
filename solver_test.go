package watersortpuzzle_test

import (
	"testing"

	watersortpuzzle "github.com/pkositsyn/water-sort-puzzle-solver"
	"github.com/pkositsyn/water-sort-puzzle-solver/solvertest"
	"github.com/stretchr/testify/suite"
)

type AStarSolverSuite struct {
	solvertest.SolverSuite
}

func aStarFactoryMethod() watersortpuzzle.Solver {
	return watersortpuzzle.NewAStarSolver()
}

func (s *AStarSolverSuite) SetupSuite() {
	s.NewSolverFunc = aStarFactoryMethod
}

func TestAStarSolver(t *testing.T) {
	suite.Run(t, new(AStarSolverSuite))
}

func BenchmarkAStarSolver(b *testing.B) {
	solvertest.TemplateBenchmarkSolve(b, aStarFactoryMethod)
}

type DijkstraSolverSuite struct {
	solvertest.SolverSuite
}

func dijkstraSolverFactoryMethod() watersortpuzzle.Solver {
	return watersortpuzzle.NewDijkstraSolver()
}

func (s *DijkstraSolverSuite) SetupSuite() {
	s.NewSolverFunc = dijkstraSolverFactoryMethod
	s.MaxFlasks = 7
}

func TestDijkstraSolver(t *testing.T) {
	suite.Run(t, new(DijkstraSolverSuite))
}

func BenchmarkDijkstraSolver(b *testing.B) {
	solvertest.TemplateBenchmarkSolve(b, dijkstraSolverFactoryMethod)
}
