package cpu6502

// ----------------------------------------------------------------------------
// inst_arithmetic.go
// Arithmetic Instructions
// ADC, SBC
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

func (c *CPU) adc(i *InstructionTableEntry) {

	operand := c.load_by_addressing_mode(i.addressingMode)
	carry := c.get_carry()
	trial := uint(operand) + uint(c.accumulator) + uint(carry)
	c.accumulator = c.accumulator + operand + carry
	c.set_carry(trial > 0xFFFF)
	c.set_overflow(int8(c.accumulator) > 127 || int8(c.accumulator) < -128)
	c.set_negative(c.accumulator)
	c.set_zero(c.accumulator)

}

func (c *CPU) sbc(i *InstructionTableEntry) {

	operand := c.load_by_addressing_mode(i.addressingMode)
	carry := c.get_carry()
	trial := int(c.accumulator) - int(operand) - int(carry)
	c.accumulator = c.accumulator - operand - carry
	c.set_carry(trial < 0)
	c.set_overflow(int8(c.accumulator) > 127 || int8(c.accumulator) < -128)
	c.set_negative(c.accumulator)
	c.set_zero(c.accumulator)
}
