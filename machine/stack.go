package machine

import (
	"cash/word"
	"errors"
)

type Stack struct {
	Data []word.Word
}

func NewStack() *Stack {
	return &Stack{make([]word.Word, 0)}
}

func (s *Stack) Push(v word.Word) {
	s.Data = append(s.Data, v)
}

func (s *Stack) Pop() (word.Word, error) {
	n := len(s.Data)
	if n == 0 {
		return 0, errors.New("stack underflow")
	}

	res := s.Data[n-1]
	s.Data = s.Data[:n-1]
	return res, nil
}

func (s *Stack) Peek() (word.Word, error) {
	n := len(s.Data)
	if n == 0 {
		return 0, errors.New("cannot peek empty stack")
	}

	return s.Data[n-1], nil
}

func (s *Stack) AccessRandom(index int) (word.Word, error) {
	n := len(s.Data)
	if index >= n {
		return 0, errors.New("invalid stack access")
	}

	return s.Data[n-index-1], nil
}
