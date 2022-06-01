package chip8

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	machine := CPU{
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
