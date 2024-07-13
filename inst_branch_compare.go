package cpu6502

import "fmt"

// ----------------------------------------------------------------------------
// inst_branch_comapre.go
// Branch and Compare Instructions
// JMP, BCC, BCS, BEQ, BNE, BMI, BPL, BVS, BVC, CMP, CPX, CPY, BIT
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

func (c *CPU) jmp(i *InstructionTableEntry) {
	switch i.addressingMode {
	case ABSOLUTE:
		c.program_counter = c.read_absolute_16()
	case INDIRECT:
		c.program_counter = c.read_indirect()
	default:
		panic(fmt.Errorf("invalid addressing mode: %d", i.addressingMode))
	}
}

func (c *CPU) bcc(i *InstructionTableEntry) {
	if !c.is_set(FLAG_CARRY) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) bcs(i *InstructionTableEntry) {
	if c.is_set(FLAG_CARRY) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) beq(i *InstructionTableEntry) {
	if c.is_set(FLAG_ZERO) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) bne(i *InstructionTableEntry) {
	if !c.is_set(FLAG_ZERO) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) bmi(i *InstructionTableEntry) {
	if c.is_set(FLAG_NEGATIVE) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) bpl(i *InstructionTableEntry) {
	if !c.is_set(FLAG_NEGATIVE) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) bvs(i *InstructionTableEntry) {
	if c.is_set(FLAG_OVERFLOW) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) bvc(i *InstructionTableEntry) {
	if !c.is_set(FLAG_OVERFLOW) {
		c.read_and_adjust_relative()
	} else {
		c.program_counter += uint16(i.bytes)
	}
}

func (c *CPU) cmp(i *InstructionTableEntry) {
	result := c.accumulator - c.load_by_addressing_mode(i.addressingMode)
	c.set_negative(result)
	c.set_zero(result)
	c.set_carry(result > c.accumulator)
	c.program_counter += uint16(i.bytes)
}

func (c *CPU) cpx(i *InstructionTableEntry) {
	result := c.x - c.load_by_addressing_mode(i.addressingMode)
	c.set_negative(result)
	c.set_zero(result)
	c.set_carry(result > c.x)
	c.program_counter += uint16(i.bytes)
}

func (c *CPU) cpy(i *InstructionTableEntry) {
	result := c.y - c.load_by_addressing_mode(i.addressingMode)
	c.set_negative(result)
	c.set_zero(result)
	c.set_carry(result > c.y)
	c.program_counter += uint16(i.bytes)
}

func (c *CPU) bit(i *InstructionTableEntry) {
	data := c.load_by_addressing_mode(i.addressingMode)
	c.set(FLAG_NEGATIVE, data&0b10000000 > 0)
	c.set(FLAG_OVERFLOW, data&0b01000000 > 0)
	c.set_zero(c.accumulator & data)
	c.program_counter += uint16(i.bytes)
}
