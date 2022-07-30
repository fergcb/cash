package machine

import (
	"cash/word"
)

type StackFrame struct {
	Stack         Stack
	ReturnAddress word.Word
	Symbols       map[string]word.Word
}

func NewStackFrame(returnAddress word.Word) *StackFrame {
	return &StackFrame{
		Stack:         Stack{},
		ReturnAddress: returnAddress,
		Symbols:       map[string]word.Word{},
	}
}
