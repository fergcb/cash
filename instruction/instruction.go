package instruction

import (
	"cash/word"
)

type InstType word.Word

const (
	NOP InstType = iota
	HALT
	DUMP
	PRINT
	PUSH
	VOID
	DUP
	DUP2
	SWAP
	ADD
	SUB
	MUL
	DIV
	MOD
	INC
	DEC
	BRA
	BRE
	BRT
	BRZ
	BRP
	BRN
	CALL
	ARG
	RETURN
	ALLOC
	DEALLOC
	WRITE
	READ
)

var mnemonics = map[InstType]string{
	NOP:     "NOP",
	HALT:    "HALT",
	PUSH:    "PUSH",
	VOID:    "VOID",
	PRINT:   "PRINT",
	DUP:     "DUP",
	DUP2:    "DUP2",
	SWAP:    "SWAP",
	ADD:     "ADD",
	SUB:     "SUB",
	MUL:     "MUL",
	DIV:     "DIV",
	MOD:     "MOD",
	INC:     "INC",
	DEC:     "DEC",
	DUMP:    "DUMP",
	BRA:     "BRA",
	BRE:     "BRE",
	BRT:     "BRT",
	BRZ:     "BRZ",
	BRP:     "BRP",
	BRN:     "BRN",
	CALL:    "CALL",
	ARG:     "ARG",
	RETURN:  "RETURN",
	ALLOC:   "ALLOC",
	DEALLOC: "DEALLOC",
	WRITE:   "WRITE",
	READ:    "READ",
}

var opcodes = MakeOpcodesMap(mnemonics)

func MakeOpcodesMap(mnemonics map[InstType]string) map[string]InstType {
	opcodes := make(map[string]InstType)
	for opcode, mnemonic := range mnemonics {
		opcodes[mnemonic] = opcode
	}
	return opcodes
}

func (instType InstType) Name() string {
	if val, ok := mnemonics[instType]; ok {
		return val
	}

	return "INVALID_INSTRUCTION"
}

func typeFromMnemonic(mnemonic string) InstType {
	if val, ok := opcodes[mnemonic]; ok {
		return val
	}

	return 0
}

type Inst struct {
	Type    InstType
	Operand word.Word
}

func (instType InstType) New(operand word.Word) *Inst {
	return &Inst{instType, operand}
}

func FromMnemonic(mnemonic string, operand word.Word) *Inst {
	instType := typeFromMnemonic(mnemonic)
	return instType.New(operand)
}
