package processor

import (
	"os"
	"bufio"
	"errors"
)

type Parser struct {
	file string
	source []string
	Tasks map[string]TaskDefinition
}

type TaskDefinition struct {
	Pos SourcePosition
	Name string
	NormalizedName string
	Dependencies []string
	Commands []string
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
		map[string]TaskDefinition{},
	}, nil
}

func (p Parser)Parse() {
	allTokens := lex(p.file, p.source)
	dependencies := []string{}
	commands := []string{}
	for i := len(allTokens)-1; i >= 0; i-- {
		switch allTokens[i].kind {
		case TYPE_RULE:
			name := allTokens[i].name
			p.Tasks[ name ] = TaskDefinition{
				allTokens[i].pos,
				allTokens[i].name,
				string(allTokens[i].name[:len(allTokens[i].name)-1]),
				dependencies,
				commands,
			}
			dependencies = []string{}
			commands = []string{}
		case TYPE_WORD:
			dependencies = append(dependencies, allTokens[i].name)
		case TYPE_COMMAND:
			commands = append(commands, allTokens[i].name)
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
