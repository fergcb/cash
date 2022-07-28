package machine

import (
	. "cash/instruction"
	. "cash/stack"
	. "cash/word"
	"errors"
	"fmt"
	"os"
)

type Bits uint8

type Flags struct {
	bits Bits
}

func (f *Flags) Set(flag Bits) {
	f.bits |= flag
}

func (f *Flags) Clear(flag Bits) {
	f.bits &^= flag
}

func (f *Flags) Toggle(flag Bits) {
	f.bits ^= flag
}

func (f *Flags) Has(flag Bits) bool {
	return f.bits&flag != 0
}

const (
	FLAG_ZERO Bits = 1 << iota
	FLAG_NEG
)

type Machine struct {
	program      []Inst
	program_size Word
	ip           Word
	stack        Stack
	flags        Flags
	halted       bool
}

func NewMachine() *Machine {
	return &Machine{
		program:      []Inst{},
		program_size: 0,
		ip:           0,
		stack:        *NewStack(),
		flags:        Flags{},
		halted:       false,
	}
}

func (m *Machine) LoadProgram(program []Inst) {
	m.program = program
	m.program_size = Word(len(program))
}

func peekN(stack *Stack, count int) ([]Word, error) {
	values := []Word{}
	for i := 0; i < count; i++ {
		value, err := stack.AccessRandom(i)
		if err != nil {
			return nil, err
		}
		values = append([]Word{value}, values...)
	}
	return values, nil
}

func popN(stack *Stack, count int) ([]Word, error) {
	values := []Word{}
	for i := 0; i < count; i++ {
		value, err := stack.Pop()
		if err != nil {
			return nil, err
		}
		values = append([]Word{value}, values...)
	}
	return values, nil
}

func (m *Machine) Execute(inst Inst) error {
	switch inst.Type {

	case INST_NOP:
		m.ip += 1
		return nil

	case INST_HALT:
		m.halted = true
		return nil

	case INST_DUMP:
		fmt.Printf("Stack: %v\nip: %d\n", m.stack.Data, m.ip)
		m.ip += 1
		return nil

	case INST_PUSH:
		value := inst.Operands[0]
		m.stack.Push(value)
		m.ip += 1
		return nil

	case INST_DUP:
		value, err := m.stack.Peek()
		if err != nil {
			return err
		}
		m.stack.Push(value)
		m.ip += 1
		return nil

	case INST_DUP2:
		values, err := peekN(&m.stack, 2)
		if err != nil {
			return err
		}
		m.stack.Push(values[0])
		m.stack.Push(values[1])
		m.ip += 1
		return nil

	case INST_ADD:
		values, err := popN(&m.stack, 2)
		if err != nil {
			return err
		}
		res := values[0] + values[1]
		m.stack.Push(res)
		m.ip += 1
		return nil

	case INST_SUB:
		values, err := popN(&m.stack, 2)
		if err != nil {
			return err
		}
		m.stack.Push(values[0] - values[1])
		m.ip += 1
		return nil

	case INST_MUL:
		values, err := popN(&m.stack, 2)
		if err != nil {
			return err
		}
		m.stack.Push(values[0] * values[1])
		m.ip += 1
		return nil

	case INST_DIV:
		values, err := popN(&m.stack, 2)
		if err != nil {
			return err
		}
		m.stack.Push(values[0] / values[1])
		m.ip += 1
		return nil

	case INST_MOD:
		values, err := popN(&m.stack, 2)
		if err != nil {
			return err
		}
		m.stack.Push(values[0] % values[1])
		m.ip += 1
		return nil

	case INST_INC:
		value, err := m.stack.Pop()
		if err != nil {
			return err
		}
		m.stack.Push(value + 1)
		m.ip += 1
		return nil

	case INST_DEC:
		value, err := m.stack.Pop()
		if err != nil {
			return err
		}
		m.stack.Push(value - 1)
		m.ip += 1
		return nil

	case INST_BRA:
		m.ip = inst.Operands[0]
		return nil

	case INST_BRE:
		values, err := popN(&m.stack, 2)
		if err != nil {
			return err
		}
		if values[0] == values[1] {
			m.ip = inst.Operands[0]
		} else {
			m.ip += 1
		}
		return nil

	case INST_BRT:
		value, err := m.stack.Peek()
		if err != nil {
			return err
		}
		if value != 0 {
			m.ip = inst.Operands[0]
		} else {
			m.ip += 1
		}
		return nil

	case INST_BRZ:
		value, err := m.stack.Peek()
		if err != nil {
			return err
		}
		if value == 0 {
			m.ip = inst.Operands[0]
		} else {
			m.ip += 1
		}
		return nil

	case INST_BRP:
		value, err := m.stack.Peek()
		if err != nil {
			return err
		}
		if value > 0 {
			m.ip = inst.Operands[0]
		} else {
			m.ip += 1
		}
		return nil

	case INST_BRN:
		value, err := m.stack.Peek()
		if err != nil {
			return err
		}
		if value < 0 {
			m.ip = inst.Operands[0]
		} else {
			m.ip += 1
		}
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
