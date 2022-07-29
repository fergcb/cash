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

		mnemonic := fields[0]
		fields = fields[1:]

		args := make([]word.Word, len(fields))
		for i, field := range fields {
			if val, ok := labels[field]; ok {
				args[i] = val
			} else {
				val, err := strconv.Atoi(field)
				if err != nil {
					return nil, fmt.Errorf("malformed operand in '%s' at field '%s'", line, field)
				}
				args[i] = word.Word(val)
			}
		}

		instruction := inst.FromMnemonic(mnemonic, args...)
		program = append(program, *instruction)
	}
	return program, nil
}
