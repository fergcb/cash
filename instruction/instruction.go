package instruction

import (
	. "cash/word"
)

type InstType Word

const (
	INST_NOP InstType = iota
	INST_HALT
	INST_DUMP
	INST_PRINT
	INST_PUSH
	INST_DUP
	INST_DUP2
	INST_SWAP
	INST_ADD
	INST_SUB
	INST_MUL
	INST_DIV
	INST_MOD
	INST_INC
	INST_DEC
	INST_BRA
	INST_BRE
	INST_BRT
	INST_BRZ
	INST_BRP
	INST_BRN
	INST_CALL
	INST_ARG
	INST_RETURN
)

var mnemonics = map[InstType]string{
	INST_NOP:    "NOP",
	INST_HALT:   "HALT",
	INST_PUSH:   "PUSH",
	INST_PRINT:  "PRINT",
	INST_DUP:    "DUP",
	INST_DUP2:   "DUP2",
	INST_SWAP:   "SWAP",
	INST_ADD:    "ADD",
	INST_SUB:    "SUB",
	INST_MUL:    "MUL",
	INST_DIV:    "DIV",
	INST_MOD:    "MOD",
	INST_INC:    "INC",
	INST_DEC:    "DEC",
	INST_DUMP:   "DUMP",
	INST_BRA:    "BRA",
	INST_BRE:    "BRE",
	INST_BRT:    "BRT",
	INST_BRZ:    "BRZ",
	INST_BRP:    "BRP",
	INST_BRN:    "BRN",
	INST_CALL:   "CALL",
	INST_ARG:    "ARG",
	INST_RETURN: "RETURN",
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
	Operand Word
}

func (instType InstType) New(operand Word) *Inst {
	return &Inst{instType, operand}
}

func FromMnemonic(mnemonic string, operand Word) *Inst {
	instType := typeFromMnemonic(mnemonic)
	return instType.New(operand)
}
