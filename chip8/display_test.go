package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplay(t *testing.T) {
	Rows, Cols := 32, 64
	x, y := 10, 10
	display := Display{Rows: Rows, Cols: Cols, Grid: make([]byte, Rows*Cols)}

	assert.Equal(t, Rows*Cols, len(display.Grid), "grid length should be rows * cols")

	for _, val := range display.Grid {
		assert.Zero(t, val, "grid values should be empty")
	}
	got := display.SetPixel(10, 10)
	assert.False(t, got, "setting off pixel should return false")

	assert.Equal(t, display.Grid[x+y*Cols], byte(1), "set pixel should still be on")

	got = display.SetPixel(10, 10)
	assert.True(t, got, "resetting same pixel should return true")

	for _, val := range display.Grid {
		assert.Zero(t, val, "grid values should be empty again")
	}
}
