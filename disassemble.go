package cpu6502

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

type disassembleRangeType int

const (
	CODE disassembleRangeType = iota
	DATA
	BYTE
	WORD
	STRING
	SKIP
)

// ----------------------------------------------------------------------------
// Type Definitions
// ----------------------------------------------------------------------------

type DisassembleOptions struct {
	FileName      string
	StartAddress  int
	FileOffset    int
	RangeFileName string
}

type disassembleRange struct {
	startAddress int
	endAddress   int
	rangeType    disassembleRangeType
}

// ----------------------------------------------------------------------------
// File Reading / Parsing
// ----------------------------------------------------------------------------

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func parseRangeFile(filename string) ([]disassembleRange, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ranges []disassembleRange
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("invalid line %d", lineNumber)
		}
		startAddress, err := strconv.ParseInt(strings.TrimPrefix(fields[0], "0x"), 16, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid line %d", lineNumber)
		}
		endAddress, err := strconv.ParseInt(strings.TrimPrefix(fields[0], "0x"), 16, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid line %d", lineNumber)
		}

		var code disassembleRangeType
		switch fields[3] {
		case "CODE":
			code = CODE
		case "DATA":
			code = DATA
		case "BYTE":
			code = BYTE
		case "WORD":
			code = WORD
		case "STRING":
			code = STRING
		case "SKIP":
			code = SKIP
		default:
			return nil, fmt.Errorf("invalid line %d", lineNumber)
		}

		ranges = append(ranges, disassembleRange{
			startAddress: int(startAddress),
			endAddress:   int(endAddress),
			rangeType:    code,
		})

	}

	return ranges, nil

}

// ----------------------------------------------------------------------------
// Mainline Disassembly
// ----------------------------------------------------------------------------

func Disassemble(options *DisassembleOptions) {

	data, err := readFile(options.FileName)
	if err != nil {
		panic(err)
	}

	var ranges []disassembleRange = nil

	if options.RangeFileName != "" {
		ranges, err = parseRangeFile(options.RangeFileName)
		if err != nil {
			panic(err)
		}
	} else {
		ranges = make([]disassembleRange, 1)
		ranges[0] = disassembleRange{
			startAddress: 0,
			endAddress:   0xffff,
			rangeType:    CODE,
		}
	}

	address := 0
	for offset := options.FileOffset; offset < len(data); {

		fmt.Printf("%04X\t\t", address+options.StartAddress)

		for _, rangeEntry := range ranges {

			if address+options.StartAddress >= rangeEntry.startAddress &&
				address+options.StartAddress <= rangeEntry.endAddress {

				switch rangeEntry.rangeType {
				case CODE:
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
				case DATA:
					fmt.Printf("db\t")
					for {
						fmt.Printf("%02X ", data[offset])
						offset++
						address++
						if address+options.StartAddress > rangeEntry.endAddress {
							break
						}
					}
				case BYTE:
					fmt.Printf("db\t%02X", data[offset])
					offset++
					address++
				case WORD:
					value := uint16(data[offset])<<8 | uint16(data[offset+1])
					fmt.Printf("dw\t%04X", value)
					offset += 2
					address += 2
				case STRING:
					fmt.Printf("ds\t")
					for value := data[offset]; value != 0; offset++ {
						fmt.Printf("%c", value)
						address++
					}
				case SKIP:
				}
			}

		}

		fmt.Println()
	}

}
