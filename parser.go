package main

type parser struct {
	file string
	source []string
	tasks map[string]taskDefinition
}

type taskDefinition struct {
	pos sourcePosition
	name string
	dependencies []string
	commands []string
}

func newParser(file string) parser {
	return parser{
		file,
		[]string {
			"root: dep2 dep1",
			"\t# root commands to execute",
			"",
			"subdep1:",
			"\t#subdep1 commands",
			"",
			"dep2: subdep1",
			"\t#dep2 commands to execute",
			"",
			"dep1:",
			"\t# dep 1 commands to execute",
		},
		map[string]taskDefinition{},
	}
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
