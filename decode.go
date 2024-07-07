package cpu6502

import "fmt"

// ----------------------------------------------------------------------------
// decode.go
// Instruction execution mainline
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
// Instruction Execution
// ----------------------------------------------------------------------------

func (cpu *CPU) ExecuteCycle() error {

	if cpu.remaining_cycles > 0 {
		cpu.remaining_cycles--
		return nil
	}

	opcode := cpu.bus.Read(cpu.program_counter)
	instruction := instructionTable[opcode]
	if instruction.instruction == UNDEFINED {
		return fmt.Errorf("invalid opcode: %02X", opcode)
	}

	ranges := make([]disassembleRange, 1)
	ranges[0] = disassembleRange{
		startAddress: 0,
		endAddress:   0xffff,
		rangeType:    CODE,
	}

	instruction_data := [3]uint8{
		cpu.bus.Read(cpu.program_counter),
		cpu.bus.Read(cpu.program_counter + 1),
		cpu.bus.Read(cpu.program_counter + 2),
	}

	if cpu.Logging {
		disassemble_line(ranges, cpu.program_counter, instruction_data)
	}

	return nil

}
