package cpu6502

// ----------------------------------------------------------------------------
// cpu.go
// Public API for the 6502 CPU Emulator
// ----------------------------------------------------------------------------
// Copyright (c) 2024 Robert L. Snyder <rob@mooneyedkitty.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Interfaces
// ----------------------------------------------------------------------------

type Bus interface {
	Read(uint16) uint8
	Write(uint16, uint8)
}

// ----------------------------------------------------------------------------
// Structures
// ----------------------------------------------------------------------------

type CPU struct {
	accumulator      uint8
	x                uint8
	y                uint8
	stack_pointer    uint8
	processor_status uint8
	program_counter  uint16
	bus              Bus
	remaining_cycles int
	handlers         map[Instruction]InstructionHandler
}

// ----------------------------------------------------------------------------
// Initialization
// ----------------------------------------------------------------------------

func NewCPU(bus Bus) *CPU {
	cpu := CPU{
		accumulator:      0,
		x:                0,
		y:                0,
		stack_pointer:    0xFF,
		processor_status: 0,
		program_counter:  0,
		bus:              bus,
		remaining_cycles: 0,
	}

	cpu.handlers = build_call_table(&cpu)
	startAddress := uint16(bus.Read(0x0ffc)) | uint16(bus.Read(0x0ffd))<<8
	cpu.program_counter = startAddress

	return &cpu
}
