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
