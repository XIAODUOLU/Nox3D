package terminal

import (
	"github.com/gdamore/tcell/v2"
)

// Screen wraps tcell screen with double buffering
type Screen struct {
	tcellScreen tcell.Screen
	width       int
	height      int
	buffer      [][]Cell
}

// Cell represents a terminal cell
type Cell struct {
	Char  rune
	Color tcell.Color
}

// NewScreen creates a new terminal screen
func NewScreen() (*Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	s.EnableMouse()
	s.Clear()

	width, height := s.Size()

	screen := &Screen{
		tcellScreen: s,
		width:       width,
		height:      height,
		buffer:      make([][]Cell, height),
	}

	for i := range screen.buffer {
		screen.buffer[i] = make([]Cell, width)
	}

	return screen, nil
}

// Close closes the screen
func (s *Screen) Close() {
	s.tcellScreen.Fini()
}

// Size returns the screen dimensions
func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

// Clear clears the buffer
func (s *Screen) Clear() {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.buffer[y][x] = Cell{Char: ' ', Color: tcell.ColorBlack}
		}
	}
}

// SetPixel sets a pixel in the buffer
func (s *Screen) SetPixel(x, y int, char rune, color tcell.Color) {
	if x >= 0 && x < s.width && y >= 0 && y < s.height {
		s.buffer[y][x] = Cell{Char: char, Color: color}
	}
}

// Flush renders the buffer to the screen
func (s *Screen) Flush() {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.buffer[y][x]
			style := tcell.StyleDefault.Foreground(cell.Color).Background(tcell.ColorBlack)
			s.tcellScreen.SetContent(x, y, cell.Char, nil, style)
		}
	}
	s.tcellScreen.Show()
}

// PollEvent polls for events
func (s *Screen) PollEvent() tcell.Event {
	return s.tcellScreen.PollEvent()
}

// HasPendingEvent checks if there's a pending event
func (s *Screen) HasPendingEvent() bool {
	return s.tcellScreen.HasPendingEvent()
}
