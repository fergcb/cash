package main

import (
	inst "cash/instruction"
	"cash/machine"
)

func main() {
	program := []inst.Inst{
		*inst.INST_PUSH.New(15),
		*inst.INST_CALL.New(4),
		*inst.INST_DUMP.New(),
		*inst.INST_HALT.New(),

		*inst.INST_ARG.New(0),
		*inst.INST_BRZ.New(21),
		*inst.INST_ARG.New(0),
		*inst.INST_PUSH.New(1),
		*inst.INST_SUB.New(),
		*inst.INST_BRZ.New(21),
		*inst.INST_ARG.New(0),
		*inst.INST_PUSH.New(1),
		*inst.INST_SUB.New(),
		*inst.INST_ARG.New(0),
		*inst.INST_PUSH.New(2),
		*inst.INST_SUB.New(),
		*inst.INST_CALL.New(4),
		*inst.INST_SWAP.New(),
		*inst.INST_CALL.New(4),
		*inst.INST_ADD.New(),
		*inst.INST_BRA.New(22),
		*inst.INST_ARG.New(0),
		*inst.INST_RETURN.New(1),
	}

	vm := machine.NewMachine()
	vm.LoadProgram(program)
	vm.Run()
}
