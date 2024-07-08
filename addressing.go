package cpu6502

import "fmt"

// ----------------------------------------------------------------------------
// addressing.go
// Reading and writing based on addressing mode
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
// Immediate

func (c *CPU) read_immediate() uint8 {
	return c.bus.Read(c.program_counter + 1)
}

// ----------------------------------------------------------------------------
// Absolute

func (c *CPU) read_absolute() uint8 {
	addr := uint16(c.bus.Read(c.program_counter+2))<<8 | uint16(c.bus.Read(c.program_counter+1))
	return c.bus.Read(addr)
}

func (c *CPU) write_absolute(data uint8) {
	addr := uint16(c.bus.Read(c.program_counter+2))<<8 | uint16(c.bus.Read(c.program_counter+1))
	c.bus.Write(addr, data)
}

// ----------------------------------------------------------------------------
// Zero Page

func (c *CPU) read_zero_page() uint8 {
	return c.bus.Read(uint16(c.bus.Read(c.program_counter + 1)))
}

func (c *CPU) write_zero_page(data uint8) {
	c.bus.Write(uint16(c.bus.Read(c.program_counter+1)), data)
}

// ----------------------------------------------------------------------------
// Absolute X

func (c *CPU) read_absolute_x() uint8 {
	addr := uint16(c.bus.Read(c.program_counter+2))<<8 | uint16(c.bus.Read(c.program_counter+1))
	return c.bus.Read(addr + uint16(c.x))
}

func (c *CPU) write_absolute_x(data uint8) {
	addr := uint16(c.bus.Read(c.program_counter+2))<<8 | uint16(c.bus.Read(c.program_counter+1))
	c.bus.Write(addr+uint16(c.x), data)
}

// ----------------------------------------------------------------------------
// Absolute Y

func (c *CPU) read_absolute_y() uint8 {
	addr := uint16(c.bus.Read(c.program_counter+2))<<8 | uint16(c.bus.Read(c.program_counter+1))
	return c.bus.Read(addr + uint16(c.y))
}

func (c *CPU) write_absolute_y(data uint8) {
	addr := uint16(c.bus.Read(c.program_counter+2))<<8 | uint16(c.bus.Read(c.program_counter+1))
	c.bus.Write(addr+uint16(c.y), data)
}

// ----------------------------------------------------------------------------
// Zero Page X

func (c *CPU) read_zero_page_x() uint8 {
	return c.bus.Read(uint16(c.bus.Read(c.program_counter+1) + c.x))
}

func (c *CPU) write_zero_page_x(data uint8) {
	c.bus.Write(uint16(c.bus.Read(c.program_counter+1)+c.x), data)
}

// ----------------------------------------------------------------------------
// Zero Page Y

func (c *CPU) read_zero_page_y() uint8 {
	return c.bus.Read(uint16(c.bus.Read(c.program_counter+1) + c.y))
}

func (c *CPU) write_zero_page_y(data uint8) {
	c.bus.Write(uint16(c.bus.Read(c.program_counter+1)+c.y), data)
}

// ----------------------------------------------------------------------------
// Indirect

func (c *CPU) read_indirect() uint16 {

	low := c.bus.Read(c.program_counter + 2)
	high := c.bus.Read(c.program_counter + 1)

	return uint16(high)<<8 | uint16(low)

}

// ----------------------------------------------------------------------------
// Indirect X

func (c *CPU) read_indirect_x() uint8 {

	low := c.bus.Read(uint16(c.bus.Read(c.program_counter + 2)))
	high := c.bus.Read(uint16(c.bus.Read(c.program_counter + 1)))
	addr := uint16(high)<<8 | uint16(low)
	return c.bus.Read(addr + uint16(c.x))

}

func (c *CPU) write_indirect_x(data uint8) {
	low := c.bus.Read(uint16(c.bus.Read(c.program_counter + 2)))
	high := c.bus.Read(uint16(c.bus.Read(c.program_counter + 1)))
	addr := uint16(high)<<8 | uint16(low)
	c.bus.Write(addr+uint16(c.x), data)
}

// ----------------------------------------------------------------------------
// Indirect Y

func (c *CPU) read_indirect_y() uint8 {

	low := c.bus.Read(uint16(c.bus.Read(c.program_counter + 2)))
	high := c.bus.Read(uint16(c.bus.Read(c.program_counter + 1)))
	addr := uint16(high)<<8 | uint16(low)
	return c.bus.Read(addr) + c.y

}

func (c *CPU) write_indirect_y(data uint8) {
	low := c.bus.Read(uint16(c.bus.Read(c.program_counter + 2)))
	high := c.bus.Read(uint16(c.bus.Read(c.program_counter + 1)))
	addr := uint16(high)<<8 | uint16(low)
	c.bus.Write(addr, data+c.y)
}

// ----------------------------------------------------------------------------
// Switch and Load
// ----------------------------------------------------------------------------

func (c *CPU) load_by_addressing_mode(addressing_mode AddressingMode) uint8 {

	switch addressing_mode {
	case IMMEDIATE:
		return c.read_immediate()
	case ABSOLUTE:
		return c.read_absolute()
	case ABSOLUTE_X:
		return c.read_absolute_x()
	case ABSOLUTE_Y:
		return c.read_absolute_y()
	case ZEROPAGE:
		return c.read_zero_page()
	case ZEROPAGE_X:
		return c.read_zero_page_x()
	case ZEROPAGE_Y:
		return c.read_zero_page_y()
	case INDIRECT_X:
		return c.read_indirect_x()
	case INDIRECT_Y:
		return c.read_indirect_y()
	default:
		panic(fmt.Errorf("invalid addressing mode for load: %02x", addressing_mode))
	}

}

func (c *CPU) store_by_addressing_mode(addressing_mode AddressingMode, data uint8) {
	switch addressing_mode {
	case ABSOLUTE:
		c.write_absolute(data)
	case ABSOLUTE_X:
		c.write_absolute_x(data)
	case ABSOLUTE_Y:
		c.write_absolute_y(data)
	case ZEROPAGE:
		c.write_zero_page(data)
	case ZEROPAGE_X:
		c.write_zero_page_x(data)
	case ZEROPAGE_Y:
		c.write_zero_page_y(data)
	case INDIRECT_X:
		c.write_indirect_x(data)
	case INDIRECT_Y:
		c.write_indirect_y(data)
	default:
		panic(fmt.Errorf("invalid addressing mode for store: %02x", addressing_mode))
	}
}
