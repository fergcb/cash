package instruction

import (
	. "cash/word"
)

type InstType Word

const (
	INST_NOP InstType = iota
	INST_HALT
	INST_DUMP
	INST_PUSH
	INST_DUP
	INST_DUP2
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
)

func (inst_type InstType) Name() string {
	var i interface{} = inst_type
	it := i.(InstType)

	switch it {
	case INST_NOP:
		return "NOP"
	case INST_HALT:
		return "HALT"
	case INST_PUSH:
		return "PUSH"
	case INST_DUP:
		return "DUP"
	case INST_DUP2:
		return "DUP2"
	case INST_ADD:
		return "ADD"
	case INST_SUB:
		return "SUB"
	case INST_MUL:
		return "MUL"
	case INST_DIV:
		return "DIV"
	case INST_MOD:
		return "MOD"
	case INST_INC:
		return "INC"
	case INST_DEC:
		return "DEC"
	case INST_DUMP:
		return "DUMP"
	case INST_BRA:
		return "BRA"
	case INST_BRE:
		return "BRE"
	case INST_BRT:
		return "BRT"
	case INST_BRZ:
		return "BRZ"
	case INST_BRP:
		return "BRP"
	case INST_BRN:
		return "BRN"
	}

	return "INVALID_INSTRUCTION"
}

type Inst struct {
	Type     InstType
	Operands []Word
}

func (instType InstType) New(operands ...Word) *Inst {
	return &Inst{instType, operands}
}
