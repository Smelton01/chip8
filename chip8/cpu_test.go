package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPU(t *testing.T) {
	machine := CPU{
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
			machine.PC = 0x200
			machine.LoadRom(tt.args.Instruction)
			code := machine.Fetch()
			assert.Equal(t, Address(0x202), machine.PC, "PC should increment")

			assert.Equal(t, uint16(tt.args.Want), code, "Instruction should be loaded")
		})
	}
}
