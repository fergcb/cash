package main

import (
	"bufio"
	"cash/error"
	"cash/machine"
	"cash/machinecode"
	"os"
	"strings"

	"github.com/alexflint/go-arg"
)

type CompileCommand struct {
	Input string `arg:"positional"`
}

type AssembleCommand struct {
	Input string `arg:"positional"`
}

type RunCommand struct {
	Input string `arg:"positional"`
}

func main() {

	var args struct {
		Compile  *CompileCommand  `arg:"subcommand:compile"`
		Assemble *AssembleCommand `arg:"subcommand:assemble"`
		Run      *RunCommand      `arg:"subcommand:run"`
	}

	arg.MustParse(&args)
	if args.Compile != nil {
		compile(args.Compile)
	} else if args.Assemble != nil {
		assemble(args.Assemble)
	} else if args.Run != nil {
		run(args.Run)
	}
}

func compile(args *CompileCommand) {
	panic("function not implemented")
}

func assemble(args *AssembleCommand) {
	panic("function not implemented")
}

func run(args *RunCommand) {
	inputPath := args.Input

	file, err := os.Open(inputPath)
	error.Check(err)
	defer file.Close()

	if strings.HasSuffix(inputPath, ".cash") {
		runCashCodeFromFile(file)
	} else if strings.HasSuffix(inputPath, ".cm") {
		runMachineCodeFromFile(file)
	}
}

func runMachineCodeFromFile(file *os.File) {
	scanner := bufio.NewScanner(file)

	program, err := machinecode.Parse(scanner)
	error.Check(err)

	vm := machine.NewMachine()
	vm.LoadProgram(program)
	vm.Run()

	err = scanner.Err()
	error.Check(err)
}

func runCashCodeFromFile(file *os.File) {
	panic("function not implemented")
}
