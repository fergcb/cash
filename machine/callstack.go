package machine

import (
	"cash/word"
	"errors"
)

type CallStack struct {
	Frames []StackFrame
}

func NewCallStack() *CallStack {
	return &CallStack{make([]StackFrame, 1)}
}

func (cs *CallStack) PushFrame(returnAddress word.Word) {
	cs.Frames = append(cs.Frames, *NewStackFrame(returnAddress))
}

func (cs *CallStack) PopFrame() (StackFrame, error) {
	n := len(cs.Frames)
	if n == 0 {
		return StackFrame{}, errors.New("call stack underflow")
	}

	res := cs.Frames[n-1]
	cs.Frames = cs.Frames[:n-1]
	return res, nil
}

func (cs *CallStack) CurrentFrame() (*StackFrame, error) {
	n := len(cs.Frames)
	if n == 0 {
		return &StackFrame{}, errors.New("the callstack is empty")
	}

	return &cs.Frames[n-1], nil
}

func (cs *CallStack) ParentFrame() (*StackFrame, error) {
	n := len(cs.Frames)
	if n < 2 {
		return &StackFrame{}, errors.New("there is no parent stack frame")
	}

	return &cs.Frames[n-2], nil
}
