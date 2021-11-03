package processor

import (
	"fmt"
	"os"
	"bufio"
	"errors"
	"mage/typedefs"
)

func loadFile(fpath string) ([]string, error) {
	fp, err := os.Open(fpath)
	if err != nil {
		return nil, errors.New("Error reading file: " + fpath)
	}
	defer fp.Close()

	lines := []string{}
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func dbgdefs(mds map[string]typedefs.MacroDefinition) {
	for n,m := range mds {
		fmt.Printf("- [%s]:\n", n)
		for _, t := range m.Tokens {
			fmt.Printf("\t> [%s] (%s)\n", toktype(t.Kind), t.Value)
		}
	}
}

func toktype(kind typedefs.TokenType) string {
	switch(kind) {
	case typedefs.TOKEN_WORD:
		return "word"
	case typedefs.TOKEN_MACRO_DFN_OPEN:
		return "macro dfn open"
	case typedefs.TOKEN_MACRO_DFN_CLOSE:
		return "macro dfn CLOSE"
	case typedefs.TOKEN_MACRO_CALL_OPEN:
		return "macro call open"
	case typedefs.TOKEN_MACRO_CALL_CLOSE:
		return "macro call CLOSE"
	case typedefs.TOKEN_RULE_OPEN:
		return "rule open"
	case typedefs.TOKEN_RULE_CLOSE:
		return "rule close"
	case typedefs.TOKEN_COMMAND_OPEN:
		return "command open"
	case typedefs.TOKEN_COMMAND_CLOSE:
		return "command close"
	}
	return "wut"
}
