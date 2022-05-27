package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/smelton01/chip8/chip8"
)

const (
	Rows       = 32
	Cols       = 64
	MemorySize = 4096
	path       = "./assets/roms/IBM.ch8"
)

var Font = [][]byte{{0xF0, 0x90, 0x90, 0x90, 0xF0},
	{0x20, 0x60, 0x20, 0x20, 0x70},
	{0xF0, 0x10, 0xF0, 0x80, 0xF0},
	{0xF0, 0x10, 0xF0, 0x10, 0xF0},
	{0x90, 0x90, 0xF0, 0x10, 0x10},
	{0xF0, 0x80, 0xF0, 0x10, 0xF0},
	{0xF0, 0x80, 0xF0, 0x90, 0xF0},
	{0xF0, 0x10, 0x20, 0x40, 0x40},
	{0xF0, 0x90, 0xF0, 0x90, 0xF0},
	{0xF0, 0x90, 0xF0, 0x10, 0xF0},
	{0xF0, 0x90, 0xF0, 0x90, 0x90},
	{0xE0, 0x90, 0xE0, 0x90, 0xE0},
	{0xF0, 0x80, 0x80, 0x80, 0xF0},
	{0xE0, 0x90, 0x90, 0x90, 0xE0},
	{0xF0, 0x80, 0xF0, 0x80, 0xF0},
	{0xF0, 0x80, 0xF0, 0x80, 0x80},
}

func main() {
	// TODO add rom to flag
	// Log to a file. Useful in debugging since you can't really log to stdout.
	// Not required.
	logfilePath := os.Getenv("BUBBLETEA_LOG")
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatal(err)
		}
	}

	rom, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	display := chip8.Display{Rows: Rows, Cols: Cols, FPS: 30, Grid: make([]byte, Rows*Cols)}

	machine := chip8.Machine{
		Memory:    make([]byte, MemorySize),
		Stack:     chip8.NewStack(),
		Frequency: 7,
		PC:        0x200,
		Index:     0x00,
		Registers: make([]byte, 16),
		Font:      Font,
		Display:   display,
	}

	machine.LoadRom(rom)
	machine.LoadFont()
	// TODO add cancellation option
	go machine.StartTimers()

	model := chip8.Model{Chip8: &machine}
	go machine.Cycle()
	// Initialize our program
	p := tea.NewProgram(model)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
