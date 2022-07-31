package machine

import (
	e "cash/error"
	inst "cash/instruction"
	"cash/word"
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
	program      []inst.Inst
	program_size word.Word
	ip           word.Word
	callStack    CallStack
	heap         Heap
	halted       bool
}

func NewMachine() *Machine {
	return &Machine{
		program:      []inst.Inst{},
		program_size: 0,
		ip:           0,
		callStack:    *NewCallStack(),
		heap:         *NewHeap(1024),
		halted:       false,
	}
}

func (m *Machine) LoadProgram(program []inst.Inst) {
	m.program = program
	m.program_size = word.Word(len(program))
}

func peekN(stack *Stack, count int) []word.Word {
	values := []word.Word{}
	for i := 0; i < count; i++ {
		value, err := stack.AccessRandom(i)
		e.Check(err)
		values = append([]word.Word{value}, values...)
	}
	return values
}

func popN(stack *Stack, count int) []word.Word {
	values := []word.Word{}
	for i := 0; i < count; i++ {
		value, err := stack.Pop()
		e.Check(err)
		values = append([]word.Word{value}, values...)
	}
	return values
}

func (m *Machine) Execute(instruction inst.Inst) error {

	frame, err := m.callStack.CurrentFrame()
	if err != nil {
		return err
	}

	stack := &frame.Stack

	switch instruction.Type {
	case inst.NOP:
		m.ip += 1
		return nil

	case inst.HALT:
		m.halted = true
		return nil

	case inst.DUMP:
		fmt.Printf("Stack: %v\nHeap: %v\nip: %d\n", stack.Data, m.heap.Data, m.ip)
		m.ip += 1
		return nil

	case inst.PRINT:
		value, err := stack.Pop()
		e.Check(err)
		fmt.Print(value)
		m.ip += 1
		return nil

	case inst.PUSH:
		value := instruction.Operand
		stack.Push(value)
		m.ip += 1
		return nil

	case inst.VOID:
		stack.Pop()
		m.ip += 1
		return nil

	case inst.DUP:
		value, err := stack.Peek()
		if err != nil {
			return err
		}
		stack.Push(value)
		m.ip += 1
		return nil

	case inst.DUP2:
		values := peekN(stack, 2)
		stack.Push(values[0])
		stack.Push(values[1])
		m.ip += 1
		return nil

	case inst.SWAP:
		values := popN(stack, 2)
		stack.Push(values[1])
		stack.Push(values[0])
		m.ip += 1
		return nil

	case inst.ADD:
		values := popN(stack, 2)
		res := values[0] + values[1]
		stack.Push(res)
		m.ip += 1
		return nil

	case inst.SUB:
		values := popN(stack, 2)
		stack.Push(values[0] - values[1])
		m.ip += 1
		return nil

	case inst.MUL:
		values := popN(stack, 2)
		stack.Push(values[0] * values[1])
		m.ip += 1
		return nil

	case inst.DIV:
		values := popN(stack, 2)
		stack.Push(values[0] / values[1])
		m.ip += 1
		return nil

	case inst.MOD:
		values := popN(stack, 2)
		stack.Push(values[0] % values[1])
		m.ip += 1
		return nil

	case inst.INC:
		value, err := stack.Pop()
		if err != nil {
			return err
		}
		stack.Push(value + 1)
		m.ip += 1
		return nil

	case inst.DEC:
		value, err := stack.Pop()
		e.Check(err)
		stack.Push(value - 1)
		m.ip += 1
		return nil

	case inst.BRA:
		m.ip = instruction.Operand
		return nil

	case inst.BRE:
		values := popN(stack, 2)
		if values[0] == values[1] {
			m.ip = instruction.Operand
		} else {
			m.ip += 1
		}
		return nil

	case inst.BRT:
		value, err := stack.Pop()
		e.Check(err)
		if value != 0 {
			m.ip = instruction.Operand
		} else {
			m.ip += 1
		}
		return nil

	case inst.BRZ:
		value, err := stack.Pop()
		e.Check(err)
		if value == 0 {
			m.ip = instruction.Operand
		} else {
			m.ip += 1
		}
		return nil

	case inst.BRP:
		value, err := stack.Pop()
		e.Check(err)
		if value > 0 {
			m.ip = instruction.Operand
		} else {
			m.ip += 1
		}
		return nil

	case inst.BRN:
		value, err := stack.Pop()
		e.Check(err)
		if value < 0 {
			m.ip = instruction.Operand
		} else {
			m.ip += 1
		}
		return nil

	case inst.CALL:
		m.callStack.PushFrame(m.ip)
		m.ip = instruction.Operand
		return nil

	case inst.ARG:
		parent, err := m.callStack.ParentFrame()
		e.Check(err)
		index := int(instruction.Operand)
		value, err := parent.Stack.AccessRandom(index)
		e.Check(err)
		stack.Push(value)
		m.ip += 1
		return nil

	case inst.RETURN:
		parent, err := m.callStack.ParentFrame()
		e.Check(err)

		argCount := int(instruction.Operand)
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

	case inst.ALLOC:
		size, err := stack.Pop()
		e.Check(err)

		offset, err := m.heap.Allocate(int(size))
		e.Check(err)

		stack.Push(word.Word(offset))

		m.ip += 1
		return nil

	case inst.DEALLOC:
		offset, err := stack.Pop()
		e.Check(err)

		size, err := stack.Pop()
		e.Check(err)

		m.heap.Deallocate(int(offset), int(size))

		m.ip += 1
		return nil

	case inst.WRITE:
		offset, err := stack.Pop()
		e.Check(err)

		value, err := stack.Pop()
		e.Check(err)

		next := m.heap.Write(int(offset), 1, []word.Word{value})

		stack.Push(word.Word(next))

		m.ip += 1
		return nil

	case inst.READ:
		offset, err := stack.Pop()
		e.Check(err)

		value := m.heap.Read(int(offset), 1)[0]

		stack.Push(value)

		m.ip += 1
		return nil

	}

	return errors.New("invalid opcode")
}

func (m *Machine) Run() int {
	for !m.halted {
		if m.ip >= word.Word(m.program_size) {
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
