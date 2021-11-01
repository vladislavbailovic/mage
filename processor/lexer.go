package processor

import (
	"mage/typedefs"
)

type ttoken struct {
	pos typedefs.SourcePosition
	kind typedefs.TokenType
	name string
}


func newToken(ttype typedefs.TokenType, file string, line int, pos int, item string) ttoken {
	return ttoken{
		typedefs.SourcePosition{ file, line, pos - len(item) },
		ttype,
		item,
	}
}

func lexLine(file string, linePos int, line string) []ttoken {
	items := []ttoken{}
	item := ""
	macroOpen := false
	commandLineOpen := false
	if len(line) == 0 {
		return items
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

		currentType := getTokenType(item, commandLineOpen)
		tok := newToken(currentType, file, linePos, i, item)
		items = append(items, tok)
		if currentType == typedefs.TOKEN_COMMAND_OPEN {
			commandLineOpen = true
		}
		if currentType == typedefs.TOKEN_MACRO_DFN_OPEN {
			macroOpen = true
		}
		item = ""
	}
	if len(item) > 0 {
		currentType := getTokenType(item, commandLineOpen)
		tok := newToken(currentType, file, linePos, len(line), item)
		items = append(items, tok)
	}
	if macroOpen {
		tok := newToken(
			typedefs.TOKEN_MACRO_DFN_CLOSE, file, linePos, len(line), "")
		items = append(items, tok)
	}
	if commandLineOpen {
		tok := newToken(
			typedefs.TOKEN_COMMAND_CLOSE, file, linePos, len(line), "")
		items = append(items, tok)
	}
	return items
}

func isMacroCall(item string) bool {
	if "$(" == string(item[0:2]) {
		return true
	}
	return false
}

func getTokenType(item string, commandLine bool) typedefs.TokenType {
	//fmt.Printf("\t\titem: [%s] [%s]\n", item, string(item[0:2]))
	if isMacroCall(item) {
		return typedefs.TOKEN_MACRO_CALL
	}

	if commandLine {
		return typedefs.TOKEN_WORD
	}

	if "\t" == string(item[:1]) {
		return typedefs.TOKEN_COMMAND_OPEN
	}

	if "macro" == item {
		return typedefs.TOKEN_MACRO_DFN_OPEN
	}

	currentType := typedefs.TOKEN_WORD
	if string(item[len(item)-1]) == ":" {
		currentType = typedefs.TOKEN_RULE
	}

	return currentType
}

func lex(file string, fileLines []string) []ttoken {
	result := []ttoken{}
	for idx, line := range fileLines {
		for _, tk := range lexLine(file, idx+1, line) {
			result = append(result, tk)
		}
	}
	return result
}
