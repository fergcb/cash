package main

import (
	// inst "cash/instruction"
	// "cash/machine"
	"bufio"
	"cash/error"
	"cash/machine"
	"cash/machinecode"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

func main() {
	var args struct {
		Input string `arg:"positional"`
	}

	arg.MustParse(&args)

	if strings.HasSuffix(args.Input, ".cm") {
		runMachineCode(args.Input)
	}

	// program := []inst.Inst{
	// 	*inst.INST_PUSH.New(15),
	// 	*inst.INST_CALL.New(4),
	// 	*inst.INST_DUMP.New(),
	// 	*inst.INST_HALT.New(),

	// 	*inst.INST_ARG.New(0),
	// 	*inst.INST_BRZ.New(21),
	// 	*inst.INST_ARG.New(0),
	// 	*inst.INST_PUSH.New(1),
	// 	*inst.INST_SUB.New(),
	// 	*inst.INST_BRZ.New(21),
	// 	*inst.INST_ARG.New(0),
	// 	*inst.INST_PUSH.New(1),
	// 	*inst.INST_SUB.New(),
	// 	*inst.INST_ARG.New(0),
	// 	*inst.INST_PUSH.New(2),
	// 	*inst.INST_SUB.New(),
	// 	*inst.INST_CALL.New(4),
	// 	*inst.INST_SWAP.New(),
	// 	*inst.INST_CALL.New(4),
	// 	*inst.INST_ADD.New(),
	// 	*inst.INST_BRA.New(22),
	// 	*inst.INST_ARG.New(0),
	// 	*inst.INST_RETURN.New(1),
	// }

	// vm := machine.NewMachine()
	// vm.LoadProgram(program)
	// vm.Run()
}

func runMachineCode(inputPath string) {
	file, err := os.Open(inputPath)
	error.Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	program, err := machinecode.Parse(scanner)
	error.Check(err)

	vm := machine.NewMachine()
	vm.LoadProgram(program)
	vm.Run()

	err = scanner.Err()
	error.Check(err)
}
