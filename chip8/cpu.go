package chip8

import (
	"encoding/binary"
	"log"
	"math/rand"
	"time"
)

// CPU represents our Chip8 Virtual machine
type CPU struct {
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

	// Key data channel
	Key         chan string
	Take        chan uint16
	PressedKeys Keys
}

// Cycle represents the CPU cycle of fetching and executing instructions.
func (m *CPU) Cycle() {
	for {
		code := m.Fetch()
		m.Decode(code)
		time.Sleep(time.Second / time.Duration(m.Frequency))
	}
}

// Fetch fetches the next instruction from the memory location from the PC.
func (m *CPU) Fetch() uint16 {
	instruction := m.Memory[m.PC : m.PC+2]

	m.PC += 2

	return binary.BigEndian.Uint16(instruction)
}

// Decode decodes and executes the provided opcode
func (m *CPU) Decode(code uint16) {
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
		if m.Registers[vx] == byte(nn) {
			m.PC += 2
		}
	case 0x4000:
		if m.Registers[vx] != byte(nn) {
			m.PC += 2
		}
	case 0x5000:
		switch code & 0x000F {
		case 0x0:
			if m.Registers[vx] == m.Registers[vy] {
				m.PC += 2
			}
		}
	case 0x6000:
		// 6XNN: Set
		m.Registers[vx] = byte(nn)
	case 0x7000:
		// 7XNN: Add
		m.Registers[vx] += byte(nn)
	case 0x8000:
		switch code & 0x000F {
		case 0x0:
			m.Registers[vx] = m.Registers[vy]
		case 0x1:
			m.Registers[vx] = m.Registers[vx] | m.Registers[vy]
		case 0x2:
			m.Registers[vx] = m.Registers[vx] & m.Registers[vy]
		case 0x3:
			m.Registers[vx] = m.Registers[vx] ^ m.Registers[vy]
		case 0x4:
			sum := m.Registers[vx] + m.Registers[vy]
			if sum > 0xFF {
				m.Registers[0xF] = 1
			}
			m.Registers[vx] = sum

		case 0x5:
			m.Registers[vx] = m.Registers[vx] - m.Registers[vy]
			if m.Registers[vx] > m.Registers[vy] {
				m.Registers[0xF] = 1
			} else {
				m.Registers[0xF] = 0
			}
		case 0x7:
			m.Registers[vx] = m.Registers[vy] - m.Registers[vx]
			if m.Registers[vx] > m.Registers[vy] {
				m.Registers[0xF] = 1
			} else {
				m.Registers[0xF] = 0
			}

		case 0x6:
			m.Registers[vx] = m.Registers[vy]
			lsb := m.Registers[vx] & 0x1
			m.Registers[vx] >>= 1
			m.Registers[0xF] = lsb

		case 0xE:
			m.Registers[vx] = m.Registers[vy]
			msb := m.Registers[vx] & 0x80
			m.Registers[vx] <<= 1
			m.Registers[0xF] = msb >> 7
		}
	case 0x9000:
		switch code & 0x000F {
		case 0x0:
			if m.Registers[vx] != m.Registers[vy] {
				m.PC += 2
			}
		}
	case 0xA000:
		m.Index = Address(nnn)
	case 0xB000:
		// BNNN: Jump with offset
		m.PC = Address(nnn) + Address(m.Registers[0])
	case 0xC000:
		rand.Seed(time.Now().Unix())
		rand := rand.Intn(int(nn))
		m.Registers[vx] = byte(rand) & byte(nn)
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
		switch code & 0x00FF {
		case 0x9E:
			key := uint16(m.Registers[vx])
			if m.PressedKeys.Contains(key) {
				m.PC += 2
			}
		case 0xA1:
			key := uint16(m.Registers[vx])
			if !m.PressedKeys.Contains(key) {
				m.PC += 2
			}
		}
	case 0xF000:
		switch code & 0x00FF {
		case 0x07:
			m.Registers[vx] = m.DelayTimer
		case 0x15:
			m.DelayTimer = m.Registers[vx]
		case 0x18:
			m.SoundTimer = m.Registers[vx]
		case 0x1E:
			m.Index += Address(m.Registers[vx])
			if m.Index&0x1000 > 0 {
				m.Registers[0xF] = 1
			}
		case 0x0A:
			// FX0A get key: blocks and waits for key input
			// decrement pc while the user has no tinputed any key
			// blocks ans waits forthe user to input some key
			key := <-m.Take
			m.Registers[vx] = byte(key)
		case 0x29:
			// FX29 :Font character
			// sets idex to address of character stored in VX
			char := m.Registers[vx] & 0x0F
			m.Index = Address(m.Memory[0x50+char])
		case 0x33:
			val := m.Registers[vx]
			r := val % 10
			val = val / 10
			c := val % 10
			val = val / 10
			l := val % 10

			m.Memory[m.Index] = l
			m.Memory[m.Index+1] = c
			m.Memory[m.Index+2] = r
		case 0x0055:
			for i := Address(0); i <= Address(vx); i++ {
				m.Memory[m.Index+i] = m.Registers[i]
			}
		case 0x0065:
			for i := Address(0); i <= Address(vx); i++ {
				m.Registers[i] = m.Memory[m.Index+i]
			}

		}
	default:
		log.Printf("Invalid instruction: %04x", code)

	}
}
