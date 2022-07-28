package main

import (
	inst "cash/instruction"
	"cash/machine"
)

func main() {
	program := []inst.Inst{
		*inst.INST_PUSH.New(35),
		*inst.INST_PUSH.New(34),
		*inst.INST_DUP2.New(),
		*inst.INST_ADD.New(),
		*inst.INST_DUMP.New(),
		*inst.INST_HALT.New(),
	}

	vm := machine.NewMachine()
	vm.LoadProgram(program)
	vm.Run()
}
