package watersortpuzzle

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// State of the game field. It is represented as ordered array of flasks.
type State []Flask

// IsTerminal state. In this state game ends.
func (s State) IsTerminal() bool {
	for _, f := range s {
		if !f.IsFinished() {
			return false
		}
	}
	return true
}

// Heuristic is a monotonic lower estimate of number of steps to reach terminal state.
// Monotonic means h(currentState) >= h(currentState with one step forward).
func (s State) Heuristic() int {
	var heuristic int

	bottomColorsCount := make(map[Color]int)
	for _, f := range s {
		if f.IsEmpty() {
			continue
		}

		// We will move at least all but the lowest color
		heuristic += f.ColorTowers() - 1
		bottomColorsCount[f.BottomColor()]++
	}

	for _, bottomColorCnt := range bottomColorsCount {
		// If several colors at bottom of flasks, then need to move at least all but one of them.
		heuristic += bottomColorCnt - 1
	}
	return heuristic
}

type stateFlaskInfo struct {
	index           int
	leftCapacity    int
	lastTowerHeight int
}

type stepChoice struct {
	flaskInfoByCapacity []stateFlaskInfo
	flaskInfoByHeight   []stateFlaskInfo
}

func (c *stepChoice) add(i stateFlaskInfo) {
	c.flaskInfoByHeight = append(c.flaskInfoByHeight, i)
	c.flaskInfoByCapacity = append(c.flaskInfoByCapacity, i)
}

func (c *stepChoice) prepare() {
	sort.Slice(c.flaskInfoByCapacity, func(i, j int) bool {
		return c.flaskInfoByCapacity[i].leftCapacity < c.flaskInfoByCapacity[j].leftCapacity
	})

	sort.Slice(c.flaskInfoByHeight, func(i, j int) bool {
		return c.flaskInfoByHeight[i].lastTowerHeight < c.flaskInfoByHeight[j].lastTowerHeight
	})
}

func (c *stepChoice) steps() []Step {
	c.prepare()

	var steps []Step

	var itCap int
	for _, flask := range c.flaskInfoByHeight {
		for itCap < len(c.flaskInfoByCapacity) && c.flaskInfoByCapacity[itCap].leftCapacity < flask.lastTowerHeight {
			itCap++
		}

		for i := itCap; i < len(c.flaskInfoByCapacity); i++ {
			from := flask.index
			to := c.flaskInfoByCapacity[i].index
			if from != to {
				steps = append(steps, Step{From: from, To: to})
			}
		}
	}
	return steps
}

// returns (map: last tower color -> flasks info, non-empty flasks indexes, empty flasks indexes)
func (s State) collectFlasksInfo() (map[Color]stepChoice, []int, []int) {
	mp := make(map[Color]stepChoice)
	var nonEmptyFlasks []int
	var emptyFlasks []int

	for i, flask := range s {
		if flask.IsEmpty() {
			emptyFlasks = append(emptyFlasks, i)
			continue
		}
		nonEmptyFlasks = append(nonEmptyFlasks, i)

		topColor, height := flask.Top()

		choice := mp[topColor]
		choice.add(stateFlaskInfo{index: i, leftCapacity: flask.Left(), lastTowerHeight: height})
		mp[topColor] = choice
	}
	return mp, nonEmptyFlasks, emptyFlasks
}

func (s State) generateStatesFromSteps(steps []Step) []State {
	var newStates []State
	for _, step := range steps {
		newState, err := s.Step(step)
		if err != nil {
			panic("logic error: cannot pour in generate steps")
		}

		newStates = append(newStates, newState)
	}
	return newStates
}

func (s State) getNonEmptyFlasksSteps(mp map[Color]stepChoice) []State {
	var newStates []State
	for _, choice := range mp {
		steps := choice.steps()
		newStates = append(newStates, s.generateStatesFromSteps(steps)...)
	}
	return newStates
}

func (s State) getEmptyFlaskSteps(nonEmptyFlasks, emptyFlasks []int) []State {
	var steps []Step
	for _, nonEmptyIdx := range nonEmptyFlasks {
		for _, emptyIdx := range emptyFlasks {
			steps = append(steps, Step{From: nonEmptyIdx, To: emptyIdx})
		}
	}
	return s.generateStatesFromSteps(steps)
}

// ReachableStates from current one in one step.
func (s State) ReachableStates() []State {
	mp, nonEmptyFlasks, emptyFlasks := s.collectFlasksInfo()
	return append(s.getNonEmptyFlasksSteps(mp), s.getEmptyFlaskSteps(nonEmptyFlasks, emptyFlasks)...)
}

// Copy state for modification.
func (s State) Copy() State {
	return append([]Flask(nil), s...)
}

// String is a unique string representation of game state.
func (s State) String() string {
	var builder strings.Builder
	for i, f := range s {
		builder.WriteString(f.String())
		if i != len(s)-1 {
			builder.WriteRune(rune(invalidColor))
		}
	}
	return builder.String()
}

// EquivalentString is a string representation of game state.
// This differs from String method by omitting information about order.
// Essentially for solving order doesn't matter.
func (s State) EquivalentString() string {
	var allStrings []string
	for _, f := range s {
		allStrings = append(allStrings, f.String())
	}
	sort.Strings(allStrings)
	return strings.Join(allStrings, string(invalidColor))
}

// FromString fills the state from String representation of the board.
func (s *State) FromString(str string) error {
	flasksStrs := strings.Split(str, string(invalidColor))
	*s = make(State, len(flasksStrs))
	for i, fStr := range flasksStrs {
		if err := (*s)[i].FromString(fStr); err != nil {
			return fmt.Errorf("cannot initialize flask from string: %w", err)
		}
	}
	return nil
}

// Step returns a new state, which is created via applying given step to current State.
func (s State) Step(step Step) (State, error) {
	newState := s.Copy()
	topColor, height := newState[step.From].PopTop()
	if err := newState[step.To].Pour(topColor, height); err != nil {
		return State{}, fmt.Errorf("failed to pour: %w", err)
	}
	return newState, nil
}

// GetStepTo returns the step, which connects current state and its child.
func (s State) GetStepTo(child State) (Step, error) {
	var step Step
	for i := 0; i < len(s); i++ {
		if s[i].Size() < child[i].Size() {
			step.To = i
			continue
		}
		if s[i].Size() > child[i].Size() {
			step.From = i
		}
	}

	newState, err := s.Step(step)
	if err != nil {
		return Step{}, fmt.Errorf("invalid child: %w", err)
	}
	if len(child) != len(newState) {
		return Step{}, fmt.Errorf("invalid child: has %d flasks, but parent has %d", len(newState), len(s))
	}
	for i := range child {
		if child[i] != newState[i] {
			return Step{}, errors.New("invalid child: parent + step and child differ")
		}
	}
	return step, nil
}
