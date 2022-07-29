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
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}

	labels := make(map[string]word.Word)
	i := 0
	for _, line := range lines {
		match := labelPattern.MatchString(line)
		if match {
			labelText := strings.TrimSpace(line[:1])
			labels[labelText] = word.Word(i)
		} else {
			lines[i] = line
			i += 1
		}
	}

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
			arg := fields[1]
			if val, ok := labels[arg]; ok {
				operand = val
			} else {
				val, err := strconv.Atoi(arg)
				if err != nil {
					return nil, fmt.Errorf("malformed operand in '%s'", line)
				}
				operand = word.Word(val)
			}
		}

		instruction := inst.FromMnemonic(mnemonic, operand)
		program = append(program, *instruction)
	}
	return program, nil
}
