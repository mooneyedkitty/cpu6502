package cpu6502

import (
	"fmt"
	"io"
	"os"
)

// ----------------------------------------------------------------------------
// disassemble.go
// 6502 Disassembler
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

type DisassembleOptions struct {
	FileName     string
	StartAddress int
	FileOffset   int
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func Disassemble(options *DisassembleOptions) {

	data, err := readFile(options.FileName)
	if err != nil {
		panic(err)
	}

	address := 0
	for offset := options.FileOffset; offset < len(data); {

		fmt.Printf("%04X\t\t", address+options.StartAddress)
		instruction := instructionTable[data[offset]]
		if instruction.instruction == UNDEFINED {
			fmt.Printf("db\t$%02X", data[offset])
			offset++
			address++
		} else {
			for i := range instruction.bytes {
				fmt.Printf("%02X ", data[offset+i])
			}
			fmt.Printf("\t")
			if instruction.bytes < 3 {
				fmt.Printf("\t")
			}
			fmt.Printf("%s\t", instruction.mnemonic)
			switch instruction.addressingMode {
			case IMMEDIATE:
				fmt.Printf("#$%02X", data[offset+1])
			case ABSOLUTE:
				addr := uint16(data[offset+2])<<8 | uint16(data[offset+1])
				fmt.Printf("$%04X", addr)
			case ZEROPAGE:
				fmt.Printf("$%02X", data[offset+1])
			case ABSOLUTE_X:
				addr := uint16(data[offset+2])<<8 | uint16(data[offset+1])
				fmt.Printf("$%04X,X", addr)
			case ABSOLUTE_Y:
				addr := uint16(data[offset+2])<<8 | uint16(data[offset+1])
				fmt.Printf("$%04X,Y", addr)
			case ZEROPAGE_X:
				fmt.Printf("$%02X,X", data[offset+1])
			case ZEROPAGE_Y:
				fmt.Printf("$%02X,Y", data[offset+1])
			case INDIRECT:
				addr := uint16(data[offset+2])<<8 | uint16(data[offset+1])
				fmt.Printf("($%04X)", addr)
			case INDIRECT_X:
				addr := uint16(data[offset+2])<<8 | uint16(data[offset+1])
				fmt.Printf("($%04X,X)", addr)
			case INDIRECT_Y:
				addr := uint16(data[offset+2])<<8 | uint16(data[offset+1])
				fmt.Printf("($%04X),Y", addr)
			case RELATIVE:
				rel := int8(data[offset+1])
				fmt.Printf("$%04X", address+int(rel)+instruction.bytes+options.StartAddress)
			}

			offset += instruction.bytes
			address += instruction.bytes
		}
		fmt.Println()
	}

}
