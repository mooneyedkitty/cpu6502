package cpu6502

// ----------------------------------------------------------------------------
// inst_logical.go
// Logical Instructions
// AND, ORA, EOR
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

func (c *CPU) and(i *InstructionTableEntry) {
	c.accumulator = c.accumulator & c.load_by_addressing_mode(i.addressingMode)
	c.set_zero(c.accumulator)
	c.set_negative(c.accumulator)
	c.program_counter += uint16(i.bytes)
}

func (c *CPU) ora(i *InstructionTableEntry) {
	c.accumulator = c.accumulator | c.load_by_addressing_mode(i.addressingMode)
	c.set_zero(c.accumulator)
	c.set_negative(c.accumulator)
	c.program_counter += uint16(i.bytes)

}

func (c *CPU) eor(i *InstructionTableEntry) {
	c.accumulator = c.accumulator ^ c.load_by_addressing_mode(i.addressingMode)
	c.set_zero(c.accumulator)
	c.set_negative(c.accumulator)
	c.program_counter += uint16(i.bytes)

}
