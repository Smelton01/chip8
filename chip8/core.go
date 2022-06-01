package chip8

import (
	"strings"
	"sync"
	"time"
)

type Address uint16

type Stack struct {
	addr []Address
	mu   sync.Mutex
}

// Key contains keys and the timestamp they were pressed.
type Key struct {
	Timestamp time.Time
	Value     uint16
}

// Keys holds the currently pressed keys.
type Keys struct {
	mu  sync.Mutex
	set map[Key]struct{}
}

// NewKeys returns an empty set of keys.
func NewKeys() Keys {
	return Keys{set: make(map[Key]struct{})}
}

// Add appends a key to the set of pressed keys.
func (k *Keys) Add(keys ...uint16) {
	k.mu.Lock()
	defer k.mu.Unlock()
	for _, key := range keys {
		cur := Key{
			Timestamp: time.Now(),
			Value:     key,
		}
		k.set[cur] = struct{}{}
	}
}

// ListenForKeys keeps track of key presses and forwards depressed key to CPU
// if get key instruction has been called
func (c *CPU) ListenForKeys(take chan uint16) {
	for {
		var key uint16
		msg := <-c.Key
		msg = strings.ToUpper(msg)
		switch msg {
		case "1":
			key = 0x1
		case "2":
			key = 0x2
		case "3":
			key = 0x3
		case "4":
			key = 0xC
		case "Q":
			key = 0x4
		case "W":
			key = 0x5
		case "E":
			key = 0x6
		case "R":
			key = 0xD
		case "A":
			key = 0x7
		case "S":
			key = 0x8
		case "D":
			key = 0x9
		case "F":
			key = 0xE
		case "Z":
			key = 0xA
		case "X":
			key = 0x0
		case "C":
			key = 0xB
		case "V":
			key = 0xF
		case "ENTER":
			key = 0x99
		default:
			// only accept valid key presses
			continue
		}
		select {
		case take <- key:
			// send key if chanel is open
		default:
			c.PressedKeys.Add(key)
		}
	}
}

// Contains check whether the key set contains a
// specified key and deletes any key which are timed out.
func (k *Keys) Contains(key uint16) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	for cur := range k.set {
		if time.Since(cur.Timestamp) > time.Millisecond*50 {
			delete(k.set, cur)
		} else if cur.Value == key {
			return true
		}
	}
	return false
}

// NewStack returns an empty stack
func NewStack() *Stack {
	return &Stack{
		addr: make([]Address, 0, 16),
		mu:   sync.Mutex{},
	}
}

// LoadRom loads the application data to memory
func (m *CPU) LoadRom(rom []byte) {
	for i := 0; i < len(rom); i++ {
		m.Memory[0x200+i] = rom[i]
	}
}

// LoadFont loads the font for drawing sprites in memory
// from address 0x050 to 0x09F
func (m *CPU) LoadFont() {
	for y, char := range m.Font {
		for x, snippet := range char {
			index := 0x50 + x + y*5
			m.Memory[index] = snippet
		}
	}
}

// StartTimers starts the delay and sound timers.
func (m *CPU) StartTimers() {
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

// Push pushes an address onto the stack.
func (s *Stack) Push(addr Address) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.addr = append(s.addr, addr)
}

// Pop removes an address from the top of the stack.
func (s *Stack) Pop() Address {
	s.mu.Lock()
	defer s.mu.Unlock()

	tail := len(s.addr) - 1
	elem := s.addr[tail]
	s.addr = s.addr[:tail]
	return elem
}
