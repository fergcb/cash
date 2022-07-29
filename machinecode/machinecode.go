package machinecode

import (
	"bufio"
	inst "cash/instruction"
	"cash/word"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var labelPattern, _ = regexp.Compile(`^[a-z_]+\s*:$`)

func Parse(scanner *bufio.Scanner) ([]inst.Inst, error) {
	lines := readLines(scanner)
	labels := parseLabels(lines)
	program, err := translateLines(lines, labels)

	return program, err
}

func readLines(scanner *bufio.Scanner) []string {
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(strings.Split(line, ";")[0])
		lines = append(lines, line)
	}
	return lines
}

func parseLabels(lines []string) map[string]word.Word {
	labels := make(map[string]word.Word)
	i := 0
	for _, line := range lines {
		match := labelPattern.MatchString(line)
		if match {
			labelText := strings.TrimSpace(line[:len(line)-1])
			labels[labelText] = word.Word(i)
		} else {
			lines[i] = line
			i += 1
		}
	}
	return labels
}

func translateLines(lines []string, labels map[string]word.Word) ([]inst.Inst, error) {
	program := make([]inst.Inst, 0)
	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) == 0 {
			continue
		}
		if len(fields) > 2 {
			return nil, fmt.Errorf("malformed instruction '%s', should have no more than 1 operand", line)
		}

		mnemonic := fields[0]

		operand := word.Word(0)
		if len(fields) == 2 {
			op, err := resolveOperand(fields[1], labels)
			if err != nil {
				return nil, err
			}
			operand = op
		}

		instruction := inst.FromMnemonic(mnemonic, operand)
		program = append(program, *instruction)
	}

	return program, nil
}

func resolveOperand(arg string, labels map[string]word.Word) (word.Word, error) {
	if val, ok := labels[arg]; ok {
		return val, nil
	} else {
		val, err := strconv.Atoi(arg)
		if err != nil {
			return 0, fmt.Errorf("malformed operand '%s'", arg)
		}
		return word.Word(val), nil
	}
}
