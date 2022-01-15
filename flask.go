package watersortpuzzle

import (
	"errors"
	"fmt"
	"unsafe"
)

// Color of water piece.
// Actually it doesn't matter what exact colors there are.
// Just need to be able to compare them by ==.
// Zero-value rune is reserved for absence of Color.
type Color rune

const (
	colorNone Color = 0

	// invalidColor is an invalid Color
	invalidColor Color = ';'

	// waterPiecesPerFlask must be > 0.
	waterPiecesPerFlask = 4
)

type Flask [waterPiecesPerFlask]Color

func (f *Flask) Size() int {
	for i, c := range f {
		if c == colorNone {
			return i
		}
	}
	return len(f)
}

func (f *Flask) Left() int {
	return len(f) - f.Size()
}

func (f *Flask) IsFull() bool {
	return f[len(f)-1] != colorNone
}

func (f *Flask) IsEmpty() bool {
	return f[0] == colorNone
}

// Pour adds a tower of a given color to the flask.
func (f *Flask) Pour(c Color, height int) error {
	if f.Left() < height {
		return errors.New("cannot pour full flask")
	}
	if c == colorNone {
		return errors.New("cannot pour none Color")
	}

	size := f.Size()
	for i := 0; i < height; i++ {
		f[size+i] = c
	}
	return nil
}

func (f *Flask) IsFinished() bool {
	if f.IsEmpty() {
		return true
	}

	if !f.IsFull() {
		return false
	}

	for i := 1; i < len(f); i++ {
		if f[i] != f[i-1] {
			return false
		}
	}
	return true
}

// ColorTowers returns number of color towers in flask.
func (f *Flask) ColorTowers() int {
	if f.IsEmpty() {
		return 0
	}

	towers := 1
	for i := 1; i < f.Size(); i++ {
		if f[i] != f[i-1] {
			towers++
		}
	}
	return towers
}

// BottomColor of flask. For empty flask returns colorNone.
func (f *Flask) BottomColor() Color {
	return f[0]
}

// Top returns last tower stats.
func (f *Flask) Top() (clr Color, height int) {
	for i := f.Size() - 1; i >= 0; i-- {
		if f[i] == clr {
			height++
			continue

		}
		if height != 0 {
			return
		}
		clr = f[i]
		height = 1
	}
	return
}

// PopTop pops the last tower of one color.
func (f *Flask) PopTop() (clr Color, height int) {
	clr, height = f.Top()
	size := f.Size()

	for i := 0; i < height; i++ {
		f[size-1-i] = colorNone
	}
	return
}

// String representation of the flask.
func (f *Flask) String() string {
	flaskSlice := f[:]
	runes := *(*[]rune)(unsafe.Pointer(&flaskSlice))
	return string(runes[:f.Size()])
}

// FromString initializes flask from string.
// String mustn't contain ';' and empty rune.
func (f *Flask) FromString(s string) error {
	if len(s) > waterPiecesPerFlask {
		return fmt.Errorf("cannot initialize flask of capacity %d from %q", waterPiecesPerFlask, s)
	}

	for i, r := range s {
		clr := Color(r)
		if clr == invalidColor || clr == colorNone {
			return errors.New("invalid color provided")
		}
		f[i] = clr
	}
	return nil
}
