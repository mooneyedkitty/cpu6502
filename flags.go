package cpu6502

// ----------------------------------------------------------------------------
// flags.go
// Flag (processor status) operation
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
// masks

const (
	FLAG_NEGATIVE = 0b10000000
	FLAG_OVERFLOW = 0b01000000
	FLAG_BRK      = 0b00010000
	FLAG_DECIMAL  = 0b00001000
	FLAG_IRQ      = 0b00000100
	FLAG_ZERO     = 0b00000010
	FLAG_CARRY    = 0b00000001
	MASK_NEGATIVE = 0b01111111
	MASK_OVERFLOW = 0b10111111
	MASK_BRK      = 0b11101111
	MASK_DECIMAL  = 0b11110111
	MASK_IRQ      = 0b11111011
	MASK_ZERO     = 0b11111101
	MASK_CARRY    = 0b11111110
)

// ----------------------------------------------------------------------------
// Flag Setting

func (c *CPU) set_negative(data uint8) {
	if data>>7 == 1 {
		c.processor_status = c.processor_status | FLAG_NEGATIVE
	} else {
		c.processor_status = c.processor_status & MASK_NEGATIVE
	}
}

func (c *CPU) set_zero(data uint8) {
	if data == 0 {
		c.processor_status = c.processor_status | FLAG_ZERO
	} else {
		c.processor_status = c.processor_status & MASK_ZERO
	}
}

func (c *CPU) set_carry(on bool) {
	if on {
		c.processor_status = c.processor_status | FLAG_CARRY
	} else {
		c.processor_status = c.processor_status & MASK_CARRY
	}
}

func (c *CPU) set_overflow(on bool) {
	if on {
		c.processor_status = c.processor_status | FLAG_OVERFLOW
	} else {
		c.processor_status = c.processor_status & MASK_OVERFLOW
	}
}

func (c *CPU) set(flag uint8, on bool) {
	if on {
		c.processor_status = c.processor_status | flag
	} else {
		c.processor_status = c.processor_status & ^flag
	}
}

// ----------------------------------------------------------------------------
// Flag Getting

func (c *CPU) get_carry() uint8 {
	if c.processor_status&FLAG_CARRY > 0 {
		return 1
	}
	return 0
}

func (c *CPU) is_set(flag uint8) bool {
	return c.processor_status&flag > 0
}
