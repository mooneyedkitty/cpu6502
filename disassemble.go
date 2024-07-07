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
// Options
// ----------------------------------------------------------------------------
// The caller can specify zero or more dump ranges in a range file. These
// ranges are used to tune the disassembly. If no ranges are provided, a single
// range of type CODE that covers all addresses is created.
// ----------------------------------------------------------------------------

type disassembleRangeType int

const (
	CODE disassembleRangeType = iota
	DATA
	BYTE
	WORD
	SKIP
)

type disassembleRange struct {
	startAddress uint16
	endAddress   uint16
	rangeType    disassembleRangeType
}

// ----------------------------------------------------------------------------
// Disassembly Options
// ----------------------------------------------------------------------------

type DisassembleOptions struct {
	FileName      string // Binary file to dump
	StartAddress  uint16 // Address to use as the base
	FileOffset    uint16 // Location to start at in the file
	RangeFileName string // Name of the optional range file
}

// ----------------------------------------------------------------------------
// File Reading / Parsing
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Read in the binary file into a byte array

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// ----------------------------------------------------------------------------
// Parse the provided range file into a slice of entries

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

		// We allow both comments and blank lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Assume this is a valid line; split it into fields.
		fields := strings.Fields(line)

		// Make sure we have three values and validate
		if len(fields) != 3 {
			return nil, fmt.Errorf("invalid line %d - expected three values", lineNumber)
		}
		startAddress, err := strconv.ParseInt(strings.TrimPrefix(fields[0], "0x"), 16, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid start address on line %d", lineNumber)
		}
		endAddress, err := strconv.ParseInt(strings.TrimPrefix(fields[0], "0x"), 16, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid end address on line %d", lineNumber)
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
		case "SKIP":
			code = SKIP
		default:
			return nil, fmt.Errorf("invalid type code on line %d", lineNumber)
		}

		// Append this data as a new entry
		ranges = append(ranges, disassembleRange{
			startAddress: uint16(startAddress),
			endAddress:   uint16(endAddress),
			rangeType:    code,
		})

	}

	return ranges, nil

}

// ----------------------------------------------------------------------------
//  Disassembly
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Dumps the current. Returns the number of bytes dumped and the string
// representation.

func disassemble_line(ranges []disassembleRange, effective_address uint16,
	data [3]uint8) (int, string) {

	var line strings.Builder
	fmt.Fprintf(&line, "%04X\t\t", effective_address)

	var bytes_dumped int = 0

	// scan all of the ranges and take appropriate action
	for _, rangeEntry := range ranges {

		// this range matches
		if effective_address >= rangeEntry.startAddress &&
			effective_address <= rangeEntry.endAddress {

			switch rangeEntry.rangeType {
			case CODE:
				// We are dumping code
				instruction := instructionTable[data[0]]

				if instruction.instruction == UNDEFINED {
					fmt.Fprintf(&line, "db\t$%02X", data[0])
					bytes_dumped = 1
				} else {
					// Print out the bytes for this instruction
					for i := range instruction.bytes {
						fmt.Fprintf(&line, "%02X ", data[i])
					}

					// Move to mnemonic column
					fmt.Fprint(&line, "\t")
					// if less than three bytes, we need another tab to
					// line things up
					if instruction.bytes < 3 {
						fmt.Fprint(&line, "\t")
					}

					// Print mnemonic column
					fmt.Fprintf(&line, "%s\t", instruction.mnemonic)

					// Print instruction data as appropriate
					switch instruction.addressingMode {
					case IMMEDIATE:
						fmt.Fprintf(&line, "#$%02X", data[1])
					case ABSOLUTE:
						addr := uint16(data[2])<<8 | uint16(data[1])
						fmt.Fprintf(&line, "$%04X", addr)
					case ZEROPAGE:
						fmt.Fprintf(&line, "$%02X", data[1])
					case ABSOLUTE_X:
						addr := uint16(data[2])<<8 | uint16(data[1])
						fmt.Fprintf(&line, "$%04X,X", addr)
					case ABSOLUTE_Y:
						addr := uint16(data[2])<<8 | uint16(data[1])
						fmt.Fprintf(&line, "$%04X,Y", addr)
					case ZEROPAGE_X:
						fmt.Fprintf(&line, "$%02X,X", data[1])
					case ZEROPAGE_Y:
						fmt.Fprintf(&line, "$%02X,Y", data[1])
					case INDIRECT:
						addr := uint16(data[2])<<8 | uint16(data[1])
						fmt.Fprintf(&line, "($%04X)", addr)
					case INDIRECT_X:
						addr := uint16(data[2])<<8 | uint16(data[1])
						fmt.Fprintf(&line, "($%04X,X)", addr)
					case INDIRECT_Y:
						addr := uint16(data[2])<<8 | uint16(data[1])
						fmt.Fprintf(&line, "($%04X),Y", addr)
					case RELATIVE:
						rel := int8(data[1])
						fmt.Fprintf(&line, "$%04X", int(effective_address)+int(rel)+instruction.bytes)
					}

					bytes_dumped = instruction.bytes
				}
			case DATA:
				fmt.Fprintf(&line, "db\t")
				for {
					fmt.Fprintf(&line, "%02X ", data[0])
					bytes_dumped++
					if effective_address+uint16(bytes_dumped) > rangeEntry.endAddress {
						break
					}
				}
			case BYTE:
				fmt.Fprintf(&line, "db\t%02X", data[0])
				bytes_dumped = 1
			case WORD:
				value := uint16(data[0])<<8 | uint16(data[1])
				fmt.Fprintf(&line, "dw\t%04X", value)
				bytes_dumped = 2
			case SKIP:
				bytes_dumped = 1
			}
		}

	}

	return bytes_dumped, line.String()

}

// ----------------------------------------------------------------------------
// Dissassembly from a Binary File

func Disassemble(options *DisassembleOptions) {

	data, err := readFile(options.FileName)
	if err != nil {
		panic(err)
	}
	var ranges []disassembleRange = nil

	if options.RangeFileName != "" {
		// Attempt to read the Range File
		ranges, err = parseRangeFile(options.RangeFileName)
		if err != nil {
			panic(err)
		}
	} else {
		// We have no range file, so create a generic CODE range
		ranges = make([]disassembleRange, 1)
		ranges[0] = disassembleRange{
			startAddress: 0,
			endAddress:   0xffff,
			rangeType:    CODE,
		}
	}

	for offset := int(options.FileOffset); offset < len(data); {
		var d [3]uint8
		if offset+2 <= 0xFFFF {
			d = [3]uint8{data[offset], data[offset+1], data[offset+2]}
		} else if offset+1 <= 0xFFFF {
			d = [3]uint8{data[offset], data[offset+1], data[0]}
		} else {
			d = [3]uint8{data[offset], data[0], data[1]}
		}
		bytes_dumped, line := disassemble_line(ranges,
			uint16(options.StartAddress)-(uint16(offset)-options.FileOffset), d)
		fmt.Println(line)
		offset += bytes_dumped
	}

}
