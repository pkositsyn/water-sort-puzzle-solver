package watersortpuzzle

import (
	"container/heap"
	"errors"
	"math"
)

type Step struct {
	From, To int
}

type Solver interface {
	Solve(initialState State) ([]Step, error)
}

type SolverWithStats interface {
	Solver
	Stats() Stats
}

var ErrNotExist = errors.New("solution doesn't exist")

type AStarSolver struct {
	heap      *distanceHeap
	parents   map[string]aStarParent
	heapElems map[string]*distanceHeapElem
	heuristic func(State) int
	stats     Stats
}

type Stats struct {
	Steps int
}

var _ Solver = (*AStarSolver)(nil)

func NewAStarSolver(opts ...AStarOption) *AStarSolver {
	solver := &AStarSolver{
		heap:      newDistanceHeap(),
		parents:   make(map[string]aStarParent),
		heapElems: make(map[string]*distanceHeapElem),
		heuristic: func(s State) int { return s.Heuristic() },
	}

	for _, opt := range opts {
		opt(solver)
	}
	return solver
}

type aStarParent struct {
	parents  []State
	distance int
}

type AStarOption func(solver *AStarSolver)

func AStarWithHeuristic(heuristic func(State) int) AStarOption {
	return func(solver *AStarSolver) {
		solver.heuristic = heuristic
	}
}

func NewDijkstraSolver() *AStarSolver {
	return NewAStarSolver(AStarWithHeuristic(func(state State) int {
		return 0
	}))
}

func (s *AStarSolver) Solve(initialState State) ([]Step, error) {
	newHeapElem := &distanceHeapElem{
		distance: s.heuristic(initialState),
		elem:     initialState,
	}
	stateStr := initialState.EquivalentString()
	s.heapElems[stateStr] = newHeapElem
	heap.Push(s.heap, newHeapElem)
	s.parents[stateStr] = aStarParent{parents: nil, distance: 0}

	for s.heap.Len() > 0 {
		s.stats.Steps++
		vertex := heap.Pop(s.heap).(*distanceHeapElem)
		state := vertex.elem
		delete(s.heapElems, state.EquivalentString())

		if state.IsTerminal() {
			return s.collectPathTo(state), nil
		}

		for _, newState := range state.ReachableStates() {
			stateStr = newState.EquivalentString()
			newRealDistance := vertex.realDistance + 1
			newDistance := newRealDistance + s.heuristic(newState)

			if parents, ok := s.parents[stateStr]; ok {
				if heapElem, ok := s.heapElems[stateStr]; ok {
					if newRealDistance < heapElem.realDistance {
						s.parents[stateStr] = aStarParent{parents: []State{state}, distance: newRealDistance}
						heapElem.distance = newDistance
						heapElem.elem = newState
						heapElem.realDistance = newRealDistance
						s.heap.Fix(heapElem)
						continue
					}
				}
				if newRealDistance == parents.distance {
					parents.parents = append(parents.parents, state)
					s.parents[stateStr] = parents
				}
				continue
			}
			s.parents[stateStr] = aStarParent{parents: []State{state}, distance: newRealDistance}

			if s.heuristic(newState) > s.heuristic(state) {
				panic("heuristic is not monotonous")
			}

			newHeapElem = &distanceHeapElem{
				distance:     newDistance,
				elem:         newState,
				realDistance: newRealDistance,
			}
			s.heapElems[stateStr] = newHeapElem
			heap.Push(s.heap, newHeapElem)
		}
	}
	return nil, ErrNotExist
}

func (s *AStarSolver) collectPathTo(state State) []Step {
	var steps []Step
	for {
		parents := s.parents[state.EquivalentString()]
		if parents.parents == nil {
			for i := 0; i < len(steps)/2; i++ {
				steps[i], steps[len(steps)-1-i] = steps[len(steps)-1-i], steps[i]
			}
			return steps
		}

		var step Step
		var realParent State
		for _, parent := range parents.parents {
			gotStep, err := parent.GetStepTo(state)
			if err == nil {
				step = gotStep
				realParent = parent
				break
			}
		}
		var uninitializedStep Step
		if step == uninitializedStep {
			panic("logic error: cannot find previous step for state")
		}

		steps = append(steps, step)
		state = realParent
	}
}

func (s *AStarSolver) Stats() Stats {
	return s.stats
}

type IDAStarSolver struct {
	heuristic    func(State) int
	path         []State
	pathVertices map[string]struct{}
	stats        Stats
}

var _ Solver = (*IDAStarSolver)(nil)

func NewIDAStarSolver(opts ...IDAStarOption) *IDAStarSolver {
	solver := &IDAStarSolver{
		heuristic:    func(s State) int { return s.Heuristic() },
		pathVertices: make(map[string]struct{}),
	}

	for _, opt := range opts {
		opt(solver)
	}
	return solver
}

type IDAStarOption func(solver *IDAStarSolver)

func IDAStarWithHeuristic(heuristic func(State) int) IDAStarOption {
	return func(solver *IDAStarSolver) {
		solver.heuristic = heuristic
	}
}

func (s *IDAStarSolver) Solve(initialState State) ([]Step, error) {
	s.path = []State{initialState}
	s.pathVertices[initialState.String()] = struct{}{}

	minDistance := s.heuristic(initialState)
	var found bool
	for {
		minDistance, found = s.iterate(initialState, minDistance)
		if found {
			return s.composePath(), nil
		}
		if minDistance == math.MaxInt {
			return nil, ErrNotExist
		}
	}
}

func (s *IDAStarSolver) iterate(state State, minDistance int) (newMinDistance int, found bool) {
	s.stats.Steps++
	newDistance := len(s.path) + s.heuristic(state)
	if newDistance > minDistance {
		return newDistance, false
	}

	if state.IsTerminal() {
		return 0, true
	}

	newMinDistance = math.MaxInt
	for _, newState := range state.ReachableStates() {
		newStateStr := newState.String()
		if _, ok := s.pathVertices[newStateStr]; ok {
			continue
		}
		s.pathVertices[newStateStr] = struct{}{}
		s.path = append(s.path, newState)

		newReachableDistance, found := s.iterate(newState, minDistance)
		if found {
			return 0, true
		}

		if newReachableDistance < newMinDistance {
			newMinDistance = newReachableDistance
		}

		delete(s.pathVertices, newStateStr)
		s.path = s.path[:len(s.path)-1]
	}
	return newMinDistance, false
}

func (s *IDAStarSolver) Stats() Stats {
	return s.stats
}

func (s *IDAStarSolver) composePath() []Step {
	var steps []Step
	for i := 1; i < len(s.path); i++ {
		newStep, err := s.path[i-1].GetStepTo(s.path[i])
		if err != nil {
			// Should never happen.
			panic(err)
		}

		steps = append(steps, newStep)
	}
	return steps
}
