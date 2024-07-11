package cpu6502

// ----------------------------------------------------------------------------
// inst_inc_dec.go
// Increment and Decrement Instructions
// INC, INX, INY, DEC, DEX, DEY
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

func (c *CPU) inc(i *InstructionTableEntry) {
	value := c.load_by_addressing_mode(i.addressingMode)
	c.store_by_addressing_mode(i.addressingMode, value+1)
	c.program_counter = c.program_counter + uint16(i.bytes)
	c.set_negative(value + 1)
	c.set_zero(value + 1)
}

func (c *CPU) inx(i *InstructionTableEntry) {
	c.x = c.x + 1
	c.program_counter = c.program_counter + uint16(i.bytes)
	c.set_negative(c.x)
	c.set_zero(c.x)
}

func (c *CPU) iny(i *InstructionTableEntry) {
	c.y = c.y + 1
	c.program_counter = c.program_counter + uint16(i.bytes)
	c.set_negative(c.y)
	c.set_zero(c.y)

}

func (c *CPU) dec(i *InstructionTableEntry) {
	value := c.load_by_addressing_mode(i.addressingMode)
	c.store_by_addressing_mode(i.addressingMode, value-1)
	c.program_counter = c.program_counter + uint16(i.bytes)
	c.set_negative(value - 1)
	c.set_zero(value - 1)

}

func (c *CPU) dex(i *InstructionTableEntry) {
	c.x = c.x - 1
	c.program_counter = c.program_counter + uint16(i.bytes)
	c.set_negative(c.x)
	c.set_zero(c.x)

}

func (c *CPU) dey(i *InstructionTableEntry) {
	c.y = c.y - 1
	c.program_counter = c.program_counter + uint16(i.bytes)
	c.set_negative(c.y)
	c.set_zero(c.y)

}
