package cpu6502

import (
	_ "embed"
	"encoding/csv"
	"strconv"
	"strings"
)

// ----------------------------------------------------------------------------
// instruction_table.go
// 6502 Processor Instruction Table
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
// Instruction Source Table
// ----------------------------------------------------------------------------

//go:embed "optable.csv"
var src string

// ----------------------------------------------------------------------------
// Type Aliases
// ----------------------------------------------------------------------------

type Opcode uint8
type Instruction uint8
type AddressingMode uint8

// ----------------------------------------------------------------------------
// Instructions
// ----------------------------------------------------------------------------

const (
	ADC Instruction = iota
	AND
	ASL
	BCC
	BCS
	BEQ
	BIT
	BMI
	BNE
	BFL
	BRK
	BVC
	BVS
	CLC
	CLI
	CLV
	CMP
	CPX
	CPY
	DEC
	DEX
	DEY
	EOR
	INC
	INX
	INY
	JMP
	JSR
	LDA
	LDX
	LDY
	LSR
	NOP
	ORA
	PHA
	PHB
	PLA
	PLP
	ROL
	ROR
	RTI
	RTS
	SBC
	SEC
	SED
	SEI
	STA
	STX
	STY
	TAX
	TAY
	TSX
	TXA
	TYA
	UNDEFINED
)

var mnemonic_map = map[string]Instruction{
	"ADC": ADC,
	"AND": AND,
	"ASL": ASL,
	"BCC": BCC,
	"BCS": BCS,
	"BEQ": BEQ,
	"BIT": BIT,
	"BMI": BMI,
	"BNE": BNE,
	"BFL": BFL,
	"BRK": BRK,
	"BVC": BVC,
	"BVS": BVS,
	"CLC": CLC,
	"CLI": CLI,
	"CLV": CLV,
	"CMP": CMP,
	"CPX": CPX,
	"CPY": CPY,
	"DEC": DEC,
	"DEX": DEX,
	"DEY": DEY,
	"EOR": EOR,
	"INC": INC,
	"INX": INX,
	"INY": INY,
	"JMP": JMP,
	"JSR": JSR,
	"LDA": LDA,
	"LDX": LDX,
	"LDY": LDY,
	"LSR": LSR,
	"NOP": NOP,
	"ORA": ORA,
	"PHA": PHA,
	"PHB": PHB,
	"PLA": PLA,
	"PLP": PLP,
	"ROL": ROL,
	"ROR": ROR,
	"RTI": RTI,
	"RTS": RTS,
	"SBC": SBC,
	"SEC": SEC,
	"SED": SED,
	"SEI": SEI,
	"STA": STA,
	"STX": STX,
	"STY": STY,
	"TAX": TAX,
	"TAY": TAY,
	"TSX": TSX,
	"TXA": TXA,
	"TYA": TYA,
}

// ----------------------------------------------------------------------------
// Adressing Modes
// ----------------------------------------------------------------------------

const (
	IMPLIED AddressingMode = iota
	IMMEDIATE
	ABSOLUTE
	ZEROPAGE
	ABSOLUTE_X
	ABSOLUTE_Y
	ZEROPAGE_X
	ZEROPAGE_Y
	INDIRECT
	INDIRECT_X
	INDIRECT_Y
	RELATIVE
)

var addressing_mode_map = map[string]AddressingMode{
	"IMP":  IMPLIED,
	"IMM":  IMMEDIATE,
	"ABS":  ABSOLUTE,
	"ZP":   ZEROPAGE,
	"ABSX": ABSOLUTE_X,
	"ABSY": ABSOLUTE_Y,
	"ZPX":  ZEROPAGE_X,
	"ZPY":  ZEROPAGE_Y,
	"IND":  INDIRECT,
	"INDX": INDIRECT_X,
	"INDY": INDIRECT_Y,
	"REL":  RELATIVE,
}

// ----------------------------------------------------------------------------
// Instruction Table Entry
// ----------------------------------------------------------------------------

type InstructionTableEntry struct {
	opcode         Opcode
	instruction    Instruction
	mnemonic       string
	addressingMode AddressingMode
	bytes          int
	cycles         int
	flags          string
}

// ----------------------------------------------------------------------------
// Instruction Table
// ----------------------------------------------------------------------------

var instructionTable = buildInstructionTable()

func buildInstructionTable() *[256]InstructionTableEntry {

	reader := csv.NewReader(strings.NewReader(src))
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var instruction_table [256]InstructionTableEntry
	for i := range 256 {
		instruction_table[i] = InstructionTableEntry{
			opcode:         Opcode(i),
			instruction:    UNDEFINED,
			mnemonic:       "",
			addressingMode: AddressingMode(UNDEFINED),
			bytes:          0,
			cycles:         0,
			flags:          "",
		}
	}

	for _, record := range records {

		if record[0] == "opcode" {
			continue
		}

		opcode, err := strconv.ParseUint(record[0][2:], 16, 8)
		if err != nil {
			panic(err)
		}
		bytes, err := strconv.ParseUint(record[3], 10, 8)
		if err != nil {
			panic(err)
		}
		cycles, err := strconv.ParseUint(record[4], 10, 8)
		if err != nil {
			panic(err)
		}

		entry := InstructionTableEntry{
			opcode:         Opcode(opcode),
			instruction:    mnemonic_map[record[1]],
			mnemonic:       record[1],
			addressingMode: addressing_mode_map[record[2]],
			bytes:          int(bytes),
			cycles:         int(cycles),
		}

		instruction_table[opcode] = entry

	}

	return &instruction_table

}
