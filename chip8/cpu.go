package chip8

import (
	"encoding/binary"
	"sync"
	"time"
)

// TODO make this interface time and move it to tui package
// type Display interface {
// 	SetPixel(x, y int)
// }

type Address uint16

// TODO maybe limit stack
type Stack struct {
	addr []Address
	m    sync.Mutex
}

// Machine represents our Chip8 Virtual machine
type Machine struct {
	// Memory represents the available memory on our machine.
	Memory []byte

	// Stack holds 16-bit addresses and is used to call subroutines/functions and return from them.
	Stack *Stack

	// DelayTimer is decremented at a rate of 60 Hz while above zero.
	DelayTimer byte

	// DelayTimer is decremented at a rate of 60 Hz while above zero.
	SoundTimer byte

	Frequency uint32

	// PC is the program counter that points at the current instruction in memory.
	PC Address

	// Index register called “I” which is used to point at locations in memory.
	Index Address

	// 16 one byte variable registers V0 to VF.
	Registers []byte

	// Display represents our display to render sprites.
	Display Display

	// Font used to load our font.
	Font [][]byte
}

// Fetch fetches the next instruction from the memory location from the PC.
func (m *Machine) Fetch() uint16 {
	// get addr from PC
	instruction := m.Memory[m.PC : m.PC+2]

	m.PC += 2

	return binary.BigEndian.Uint16(instruction)
}

// LoadFont loads the font for drawing sprites in memory
// from address 0x050 to 0x09F
func (m *Machine) LoadFont() {
	for y, char := range m.Font {
		for x, snippet := range char {
			index := 0x50 + x + y*5
			m.Memory[index] = snippet
		}
	}
}

// Cycle represents the CPU cycle of fetching and executing instructions.
func (m *Machine) Cycle() {
	code := m.Fetch()
	m.Decode(code)
}

// Decode decodes and executes the provided opcode
func (m *Machine) Decode(code uint16) {
	// Extract important values
	vx, vy := code&0x0F00>>(4*2), code&0x00F0>>(4*1)
	n := code & 0x000F
	nn := code & 0x00FF
	nnn := code & 0x0FFF

	switch code & 0xF000 {
	case 0x0000:
		switch code & 0x0FFF {
		case 0x00E0:
			// 00E0 (clear screen)
			m.Display.Clear()
		case 0x00EE:
			// 0x00EE return from subroutine
			m.PC = m.Stack.Pop()
		}
	case 0x1000:
		// 1NNN (jump)
		m.PC = Address(nnn)
	case 0x2000:
		// 0x2000 call subroutine at nnn
		m.Stack.Push(m.PC)
		m.PC = Address(nnn)
	case 0x3000:
		if m.Memory[vx] == byte(nn) {
			m.PC += 2
		}
	case 0x4000:
		if m.Memory[vx] != byte(nn) {
			m.PC += 2
		}
	case 0x5000:
		switch code & 0x000F {
		case 0x0:
			if m.Memory[vx] == m.Memory[vy] {
				m.PC += 2
			}
		}
	case 0x6000:
		// 6XNN: Set
		// log.Println(nn, vx)
		m.Registers[vx] = byte(nn)
	case 0x7000:
		// 7XNN: Add
		m.Registers[vx] += byte(nn)
	case 0x8000:
	case 0x9000:
		switch code & 0x000F {
		case 0x0:
			if m.Memory[vx] != m.Memory[vy] {
				m.PC += 2
			}
		}
	case 0xA000:
		m.Index = Address(nnn)
		// panic(m.Memory[m.Index:])
	case 0xB000:
	case 0xC000:
	case 0xD000:
		// DXYN: Display
		x := m.Registers[vx]
		y := m.Registers[vy]

		x = x & 63
		y = y & 31

		m.Registers[len(m.Registers)-1] = 0

		for row := 0; row < int(n); row++ {
			index := m.Index
			sprite := m.Memory[index+Address(row)]

			width := 8

			for col := 0; col < width; col++ {
				// check if first bit is one
				if sprite&0x80 > 0 {
					if m.Display.SetPixel(int(x)+col, int(y)+row) {
						m.Registers[0xF] = 1
					}
				}
				sprite <<= 1
			}
		}
	case 0xE000:
	case 0xF000:
	}
}

// LoadRom loads the application dataa to memory
func (m *Machine) LoadRom(rom []byte) {
	for i := 0; i < len(rom); i++ {
		m.Memory[0x200+i] = rom[i]
	}
}

// StartTimers starts the delay and sound timers.
func (m *Machine) StartTimers() {
	for {
		// beep if above zero
		if m.DelayTimer > 0 {
			m.DelayTimer--
		}
		if m.DelayTimer > 0 {
			m.DelayTimer--
		}
		time.Sleep(time.Second / 60)
	}
}

// NewStack returns an empty stack
func NewStack() *Stack {
	return &Stack{
		addr: make([]Address, 0, 16),
		m:    sync.Mutex{},
	}
}

// Push pushes an address onto the stack.
func (s *Stack) Push(addr Address) {
	s.m.Lock()
	defer s.m.Unlock()

	s.addr = append(s.addr, addr)
}

// Pop removes an address from the top of the stack.
func (s *Stack) Pop() Address {
	s.m.Lock()
	defer s.m.Unlock()

	tail := len(s.addr) - 1
	elem := s.addr[tail]
	s.addr = s.addr[:tail]
	return elem
}
