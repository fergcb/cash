package stack

import (
	. "cash/word"
	"errors"
)

type Stack struct {
	Data []Word
}

func NewStack() *Stack {
	return &Stack{make([]Word, 0)}
}

func (s *Stack) Push(v Word) {
	s.Data = append(s.Data, v)
}

func (s *Stack) Pop() (Word, error) {
	n := len(s.Data)
	if n == 0 {
		return 0, errors.New("stack underflow")
	}

	res := s.Data[n-1]
	s.Data = s.Data[:n-1]
	return res, nil
}

func (s *Stack) Peek() (Word, error) {
	n := len(s.Data)
	if n == 0 {
		return 0, errors.New("cannot peek empty stack")
	}

	return s.Data[n-1], nil
}

func (s *Stack) AccessRandom(index int) (Word, error) {
	n := len(s.Data)
	if index >= n {
		return 0, errors.New("invalid stack access")
	}

	return s.Data[n-index-1], nil
}
