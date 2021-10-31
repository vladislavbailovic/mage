package ruleset

import (
	"os"
	"bufio"
	"errors"
)

type parser struct {
	file string
	source []string
	tasks map[string]taskDefinition
}

type taskDefinition struct {
	pos sourcePosition
	name string
	normalizedName string
	dependencies []string
	commands []string
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

func NewParser(file string) (parser, error) {
	lines, err := loadFile(file)
	if err != nil {
		return parser{}, err
	}
	return parser{
		file,
		lines,
		map[string]taskDefinition{},
	}, nil
}

func (p parser)parse() {
	allTokens := lex(p.file, p.source)
	dependencies := []string{}
	commands := []string{}
	for i := len(allTokens)-1; i >= 0; i-- {
		switch allTokens[i].kind {
		case TYPE_RULE:
			name := allTokens[i].name
			p.tasks[ name ] = taskDefinition{
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

func (p parser)knowsAboutTask(name string) bool {
	for tname, _ := range p.tasks {
		if name + ":" == tname {
			return true
		}
	}
	return false
}
