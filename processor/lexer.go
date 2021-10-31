package processor

type SourcePosition struct {
	File string
	Line int
	Char int
}

type token struct {
	pos SourcePosition
	kind tokenType
	name string
}

type tokenType int

const (
	TYPE_WORD tokenType = iota
	TYPE_RULE
	TYPE_COMMAND
)

func newToken(ttype tokenType, file string, line int, pos int, item string) token {
	return token{
		SourcePosition{ file, line, pos - len(item) },
		ttype,
		item,
	}
}

func lexLine(file string, linePos int, line string) []token {
	items := []token{}
	item := ""
	if len(line) == 0 {
		return items
	}
	if string(line[0]) == "\t" {
		return []token{newToken(TYPE_COMMAND, file, linePos, 0, line)}
	}
	for i := 0; i < len(line); i++ {
		if string(line[i]) == "#" {
			// comment
			break
		}
		if string(line[i]) != " " {
			item += string(line[i])
			continue
		}

		currentType := TYPE_WORD
		// can we be more specific about the type?
		if string(item[len(item)-1]) == ":" {
			currentType = TYPE_RULE
		}

		tok := newToken(currentType, file, linePos, i, item)
		items = append(items, tok)
		item = ""
	}
	if len(item) > 0 {
		currentType := TYPE_WORD
		// can we be more specific about the type?
		if string(item[len(item)-1]) == ":" {
			currentType = TYPE_RULE
		}
		tok := newToken(currentType, file, linePos, len(line), item)
		items = append(items, tok)
	}
	return items
}

func lex(file string, fileLines []string) []token {
	result := []token{}
	for idx, line := range fileLines {
		for _, tk := range lexLine(file, idx+1, line) {
			result = append(result, tk)
		}
	}
	return result
}
