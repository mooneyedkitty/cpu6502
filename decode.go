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
// Call table
// ----------------------------------------------------------------------------

type InstructionHandler func(*InstructionTableEntry)

func build_call_table(cpu *CPU) map[Instruction]InstructionHandler {

	t := make(map[Instruction]InstructionHandler)

	t[LDA] = cpu.lda
	t[LDX] = cpu.ldx
	t[LDY] = cpu.ldy
	t[STA] = cpu.sta
	t[STX] = cpu.stx
	t[STY] = cpu.sty

	t[ADC] = cpu.adc
	t[SBC] = cpu.sbc

	return t
}

// ----------------------------------------------------------------------------
// Instruction Execution
// ----------------------------------------------------------------------------

func (cpu *CPU) ExecuteCycle() error {

	// Burn off any remaining cycles from the last instruction
	if cpu.remaining_cycles > 0 {
		cpu.remaining_cycles--
		return nil
	}

	// Read the opcode and load the instruction
	opcode := cpu.bus.Read(cpu.program_counter)
	instruction := instructionTable[opcode]
	if instruction.instruction == UNDEFINED {
		return fmt.Errorf("invalid opcode: %02X", opcode)
	}

	// Execute the instruction and set remaining cycles

	handler, found := cpu.handlers[instruction.instruction]
	if !found {
		return fmt.Errorf("unsupported instruction: %s", instruction.mnemonic)
	}
	handler(&instruction)
	cpu.remaining_cycles = instruction.cycles - 1

	return nil

}
