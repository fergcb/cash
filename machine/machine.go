package machine

import (
	e "cash/error"
	. "cash/instruction"
	. "cash/word"
	"errors"
	"fmt"
	"os"
)

type Bits uint8

const (
	FLAG_ZERO Bits = 1 << iota
	FLAG_NEG
)

type Machine struct {
	program      []Inst
	program_size Word
	ip           Word
	callStack    CallStack
	halted       bool
}

func NewMachine() *Machine {
	return &Machine{
		program:      []Inst{},
		program_size: 0,
		ip:           0,
		callStack:    *NewCallStack(),
		halted:       false,
	}
}

func (m *Machine) LoadProgram(program []Inst) {
	m.program = program
	m.program_size = Word(len(program))
}

func peekN(stack *Stack, count int) []Word {
	values := []Word{}
	for i := 0; i < count; i++ {
		value, err := stack.AccessRandom(i)
		e.Check(err)
		values = append([]Word{value}, values...)
	}
	return values
}

func popN(stack *Stack, count int) []Word {
	values := []Word{}
	for i := 0; i < count; i++ {
		value, err := stack.Pop()
		e.Check(err)
		values = append([]Word{value}, values...)
	}
	return values
}

func (m *Machine) Execute(inst Inst) error {

	frame, err := m.callStack.CurrentFrame()
	if err != nil {
		return err
	}

	stack := &frame.Stack

	switch inst.Type {
	case NOP:
		m.ip += 1
		return nil

	case HALT:
		m.halted = true
		return nil

	case DUMP:
		fmt.Printf("Stack: %v\nip: %d\n", stack.Data, m.ip)
		m.ip += 1
		return nil

	case PRINT:
		value, err := stack.Pop()
		e.Check(err)
		fmt.Print(value)
		m.ip += 1
		return nil

	case PUSH:
		value := inst.Operand
		stack.Push(value)
		m.ip += 1
		return nil

	case DUP:
		value, err := stack.Peek()
		if err != nil {
			return err
		}
		stack.Push(value)
		m.ip += 1
		return nil

	case DUP2:
		values := peekN(stack, 2)
		stack.Push(values[0])
		stack.Push(values[1])
		m.ip += 1
		return nil

	case SWAP:
		values := popN(stack, 2)
		stack.Push(values[1])
		stack.Push(values[0])
		m.ip += 1
		return nil

	case ADD:
		values := popN(stack, 2)
		res := values[0] + values[1]
		stack.Push(res)
		m.ip += 1
		return nil

	case SUB:
		values := popN(stack, 2)
		stack.Push(values[0] - values[1])
		m.ip += 1
		return nil

	case MUL:
		values := popN(stack, 2)
		stack.Push(values[0] * values[1])
		m.ip += 1
		return nil

	case DIV:
		values := popN(stack, 2)
		stack.Push(values[0] / values[1])
		m.ip += 1
		return nil

	case MOD:
		values := popN(stack, 2)
		stack.Push(values[0] % values[1])
		m.ip += 1
		return nil

	case INC:
		value, err := stack.Pop()
		if err != nil {
			return err
		}
		stack.Push(value + 1)
		m.ip += 1
		return nil

	case DEC:
		value, err := stack.Pop()
		e.Check(err)
		stack.Push(value - 1)
		m.ip += 1
		return nil

	case BRA:
		m.ip = inst.Operand
		return nil

	case BRE:
		values := popN(stack, 2)
		if values[0] == values[1] {
			m.ip = inst.Operand
		} else {
			m.ip += 1
		}
		return nil

	case BRT:
		value, err := stack.Pop()
		e.Check(err)
		if value != 0 {
			m.ip = inst.Operand
		} else {
			m.ip += 1
		}
		return nil

	case BRZ:
		value, err := stack.Pop()
		e.Check(err)
		if value == 0 {
			m.ip = inst.Operand
		} else {
			m.ip += 1
		}
		return nil

	case BRP:
		value, err := stack.Pop()
		e.Check(err)
		if value > 0 {
			m.ip = inst.Operand
		} else {
			m.ip += 1
		}
		return nil

	case BRN:
		value, err := stack.Pop()
		e.Check(err)
		if value < 0 {
			m.ip = inst.Operand
		} else {
			m.ip += 1
		}
		return nil

	case CALL:
		m.callStack.PushFrame(m.ip)
		m.ip = inst.Operand
		return nil

	case ARG:
		parent, err := m.callStack.ParentFrame()
		e.Check(err)
		index := int(inst.Operand)
		value, err := parent.Stack.AccessRandom(index)
		e.Check(err)
		stack.Push(value)
		m.ip += 1
		return nil

	case RETURN:
		parent, err := m.callStack.ParentFrame()
		e.Check(err)

		argCount := int(inst.Operand)
		for i := 0; i < argCount; i++ {
			parent.Stack.Pop()
		}

		returnValue, err := stack.Peek()
		if err == nil {
			parent.Stack.Push(returnValue)
		}

		m.callStack.PopFrame()
		m.ip = frame.ReturnAddress + 1
		return nil
	}

	return errors.New("invalid opcode")
}

func (m *Machine) Run() int {
	for !m.halted {
		if m.ip >= Word(m.program_size) {
			fmt.Fprintf(os.Stderr, "invalid instruction access @%d\n", m.ip)
			return 1
		}

		inst := m.program[m.ip]
		err := m.Execute(inst)
		if err != nil {
			msg := fmt.Sprintf("error executing instruction %s at ip=%d:\n  %s", inst.Type.Name(), m.ip, err.Error())
			fmt.Fprintln(os.Stderr, msg)
			return 1
		}
	}

	return 0
}
