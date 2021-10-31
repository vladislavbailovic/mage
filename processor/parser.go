package processor

import (
	"os"
	"bufio"
	"errors"

	"mage/typedefs"
)

type Parser struct {
	file string
	source []string
	Tasks map[string]typedefs.TaskDefinition
	Macros map[string]typedefs.MacroDefinition
}

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

func NewParser(file string) (Parser, error) {
	lines, err := loadFile(file)
	if err != nil {
		return Parser{}, err
	}
	return Parser{
		file,
		lines,
		map[string]typedefs.TaskDefinition{},
		map[string]typedefs.MacroDefinition{},
	}, nil
}

func (p Parser)Parse() {
	allTokens := lex(p.file, p.source)
	dependencies := []string{}
	commands := []string{}
	for i := len(allTokens)-1; i >= 0; i-- {
		switch allTokens[i].kind {
		case typedefs.TOKEN_MACRO_DFN_OPEN:
			pos := allTokens[i].pos
			name := allTokens[i+1].name
			value := ""
			for j := i+2; j < len(allTokens) - 1; j++ {
				if allTokens[j].kind == typedefs.TOKEN_MACRO_DFN_CLOSE {
					break
				}
				if len(value) == 0 {
					value = allTokens[j].name
				} else {
					value = value + " " + allTokens[j].name
				}
			}
			p.Macros[ name ] = typedefs.MacroDefinition{
				pos, name, value,
			}
		}
	}

	for iter := 0; iter < 100; iter++ {
		replaced := false
		for i := len(allTokens)-1; i >= 0; i-- {
			switch allTokens[i].kind {
			case typedefs.TOKEN_MACRO_CALL:
				rawName := allTokens[i].name
				name := string(rawName[2:len(rawName)-1])
				macro, ok := p.Macros[ name ]
				if !ok {
					panic("Unknown macro: " + name)
				}
				allTokens[i].name = macro.Value
				if !isMacroCall(macro.Value) {
					allTokens[i].kind = typedefs.TOKEN_WORD
				} else {
				}
				replaced = true
			}
		}
		if !replaced {
			break
		}
	}

	inCommandBlock := false;
	for i := len(allTokens)-1; i >= 0; i-- {
		switch allTokens[i].kind {
		case typedefs.TOKEN_RULE:
			name := allTokens[i].name
			p.Tasks[ name ] = typedefs.TaskDefinition{
				allTokens[i].pos,
				allTokens[i].name,
				string(allTokens[i].name[:len(allTokens[i].name)-1]),
				dependencies,
				commands,
			}
			dependencies = []string{}
			commands = []string{}
		case typedefs.TOKEN_WORD:
			if inCommandBlock {
				commands = append(commands, allTokens[i].name)
			} else {
				dependencies = append(dependencies, allTokens[i].name)
			}
		case typedefs.TOKEN_COMMAND_OPEN:
			inCommandBlock = false
			commands = append(commands, allTokens[i].name)
		case typedefs.TOKEN_COMMAND_CLOSE:
			inCommandBlock = true
		}
	}
}

func (p Parser)knowsAboutTask(name string) bool {
	for tname, _ := range p.Tasks {
		if name + ":" == tname {
			return true
		}
	}
	return false
}
