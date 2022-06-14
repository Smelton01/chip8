package chip8

import (
	"bytes"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Chip8 *CPU
}

type Display struct {
	Cols int
	Rows int
	M    int
	Grid []byte
	FPS  int
}

// SetPixel toggles thi value of the pixel at
// x,y and return true if a set pixel was unset.
func (d *Display) SetPixel(x, y int) bool {
	if x > d.Cols {
		x -= d.Cols
	}
	if x < 0 {
		x += d.Cols
	}
	if y > d.Rows {
		y -= d.Rows
	}
	if y < 0 {
		y += d.Rows
	}

	index := x + y*d.Cols
	d.Grid[index] ^= 1

	return d.Grid[index] != 1
}

// Init starts the timer.
func (m Model) Init() tea.Cmd {
	return m.tick
}

// Update is called when messages are received.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		m.Chip8.Key <- msg.String()
	case tickMsg:
		return m, m.tick
	}
	return m, nil
}

// Clear resets all pixels to off
func (d *Display) Clear() {
	d.Grid = make([]byte, d.Rows*d.Cols)
}

// View returns a string to be rendered to the terminal.
func (m Model) View() string {
	display := m.Chip8.Display
	head := "Hi. First iteration of our Chip8 machine!\n"

	var body bytes.Buffer

	body.Write(bytes.Repeat([]byte{'-'}, m.Chip8.Display.Cols+2))
	body.WriteRune('\n')
	for y := 0; y < display.Rows; y++ {
		body.WriteByte('|')
		for x := 0; x < display.Cols; x++ {
			cur := ' '
			if display.Grid[x+y*display.Cols] == 1 {
				cur = '█'
			}
			body.WriteRune(cur)
		}
		body.WriteString("|\n")
	}
	body.Write(bytes.Repeat([]byte{'-'}, m.Chip8.Display.Cols+2))
	return head + body.String()
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func (m *Model) tick() tea.Msg {
	time.Sleep(time.Second / time.Duration(m.Chip8.Display.FPS))
	return tickMsg{}
}
