package cpu6502

import (
	"fmt"
	"testing"
)

//go:embed instruction_test.bin
var instructionTest []byte

// ----------------------------------------------------------------------------
// cpu_test.go
// Tests the 6502 emulator
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
// Memory with the Bus Implementation
// ----------------------------------------------------------------------------

type Memory struct {
	data [0xFFFF]uint8
}

func (m *Memory) Read(addr uint16) uint8 {
	return m.data[addr]
}

func (m *Memory) Write(addr uint16, data uint8) {
	m.data[addr] = data
}

// ----------------------------------------------------------------------------
// Instruction Set Test
// ----------------------------------------------------------------------------

func TestCPUInstructions(t *testing.T) {

	memory := Memory{}
	copy(memory.data[:], instructionTest)
	cpu := NewCPU(&memory)
	fmt.Printf("cpu: %v\n", cpu)

}
