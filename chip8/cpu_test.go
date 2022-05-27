package chip8

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	machine := Machine{
		Memory: make([]byte, 4096),
		PC:     0x200,
	}
	testCases := []struct {
		desc string
		args struct {
			Instruction []byte
			Want        Address
		}
	}{
		{
			desc: "Zero value",
			args: struct {
				Instruction []byte
				Want        Address
			}{
				Instruction: []byte{0, 0},
				Want:        Address(0)},
		}, {
			desc: "something simple",
			args: struct {
				Instruction []byte
				Want        Address
			}{
				Instruction: []byte{0xAF, 0xFA},
				Want:        0xAFFA,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			machine := machine
			machine.LoadRom(tt.args.Instruction)
			code := machine.Fetch()
			assert.Equal(t, Address(0x202), machine.PC, "PC should increment")

			assert.Equal(t, uint16(tt.args.Want), code, "Instruction should be loaded")
		})
	}
}

func TestStack(t *testing.T) {
	machine := Machine{
		Stack: NewStack(),
	}
	for _, val := range machine.Stack.addr {
		assert.Zero(t, val, "Stack should start empty")
	}

	instruction := Address(0xFFCC)
	machine.Stack.Push(instruction)
	assert.Contains(t, machine.Stack.addr, instruction, "stack should contain added instruction")

	assert.Equal(t, len(machine.Stack.addr), 1)

	log.Println(machine.Stack)
	got := machine.Stack.Pop()
	assert.Equal(t, got, instruction, "should return added item")

	log.Println(machine.Stack)
	for _, val := range machine.Stack.addr {
		assert.Zero(t, val, "Stack should be now empty")
	}
}
