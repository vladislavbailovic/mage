package processor

import (
	"fmt"
	"mage/typedefs"
)

type token struct {
	pos typedefs.SourcePosition
	kind typedefs.TokenType
	name string
}


func newToken(ttype typedefs.TokenType, file string, line int, pos int, item string) token {
	return token{
		typedefs.SourcePosition{ file, line, pos - len(item) },
		ttype,
		item,
	}
}

func lexLine(file string, linePos int, line string) []token {
	items := []token{}
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

func lex(file string, fileLines []string) []token {
	result := []token{}
	for idx, line := range fileLines {
		for _, tk := range lexLine(file, idx+1, line) {
			result = append(result, tk)
		}
	}
	return result
}







type xtoken struct {
	pos typedefs.SourcePosition
	kind typedefs.TokenType
	value string
}

func xlex(file string, content string) []xtoken {
	content += " \n"
	currentLine := 1
	currentChar := 1
	allTokens := []xtoken{}
	word := ""
	for pos := 0; pos < len(content) - 1; {

		if "\n" == string(content[pos]) {
			if len(word) > 0 {
				allTokens = append(allTokens, xtoken{
					typedefs.SourcePosition{file, currentLine, currentChar + 1},
					typedefs.TOKEN_WORD,
					word,
				})
				word = ""
			}
			currentLine += 1
			currentChar = 1
			pos++
			continue
		}

		if "\t" == string(content[pos]) {
			// Command
			pos++
			command := consumeUntil("\n", content, pos)
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_COMMAND_OPEN,
				"",
			})
			for _, tk := range xlex(file, command) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_COMMAND_CLOSE,
				"",
			})
			pos += len(command)
			word = ""
			continue
		}

		if len(content) >= pos+5 && "macro" == string(content[pos:pos+5]) {
			// Macro definition
			pos += 5 + 1
			macro := consumeUntil("\n", content, pos)
			fmt.Printf("\tmacro [%s]\n", macro)
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_MACRO_DFN_OPEN,
				"",
			})
			for _, tk := range xlex(file, macro) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_MACRO_DFN_CLOSE,
				"",
			})
			pos += len(macro)
			word = ""
			continue
		}

		if len(content) >= pos+2 && "$(" == string(content[pos:pos+2]) {
			if len(word) > 0 {
				allTokens = append(allTokens, xtoken{
					typedefs.SourcePosition{file, currentLine, currentChar + 1},
					typedefs.TOKEN_WORD,
					word,
				})
				word = ""
			}
			// Macro call
			pos += 2
			macro := consumeUntil(")", content, pos)
			fmt.Printf("macro call: [%s]\n", macro)
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_MACRO_CALL_OPEN,
				"",
			})
			for _, tk := range xlex(file, macro) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_MACRO_CALL_CLOSE,
				"",
			})
			pos += len(macro) + 1
			word = ""
			continue
		}

		if ":" == string(content[pos]) {
			// Rule definition
			name := consumeBackUntil("\n", content, pos-1)
			dependencies := consumeUntil("\n", content, pos+1)
			fmt.Printf("adding rule dfn: [%v] [%v]\n", name, dependencies)
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_RULE_OPEN,
				"",
			})
			for _, tk := range xlex(file, name) {
				allTokens = append(allTokens, tk)
			}
			for _, tk := range xlex(file, dependencies) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_RULE_CLOSE,
				"",
			})
			pos += len(dependencies) + len(name)
			word = ""
			continue
		}

		if " " == string(content[pos]) && len(word) > 0 {
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_WORD,
				word,
			})
			word = ""
			pos += 1
			currentChar += 1
			continue
		}

		word += string(content[pos])

		pos += 1
		currentChar += 1
	}
	return allTokens
}

func positionAfter(what string, source string, pos int) int {
	for i := pos; i<len(source)-1; i++ {
		if string(source[i:i+1]) == what {
			return i+1
		}
	}
	return len(source)-1
}

func consumeUntil(what string, source string, pos int) string {
	item := ""
	for i := pos; i < len(source) - 1; i++ {
		chr := string(source[i])
		if chr == what {
			break
		}
		item += chr
	}
	return item
}

func consumeBackUntil(what string, source string, pos int) string {
	item := ""
	for i := pos; i > 0; i-- {
		chr := string(source[i])
		if chr == what {
			break
		}
		item = chr + item
	}
	return item
}

func filterTokens(tokens []xtoken, expected typedefs.TokenType) []xtoken {
	result := []xtoken{}
	for _, tk := range tokens {
		if expected != tk.kind {
			continue
		}
		result = append(result, tk)
	}
	return result
}

// func transform(tokens []xtoken) []xtoken {
// 	macros := filterTokens(tokens, typedefs.TOKEN_MACRO_DFN)
// 	for _, macroToken := range macros {
// 		for i, tk := range tokens {
// 			rpl := "$(" + macroToken.name + ")"
// 			fmt.Println(rpl)
// 			tokens[i].name = strings.Replace(tk.name, rpl, macroToken.value, -1)
// 			tokens[i].value = strings.Replace(tk.value, rpl, macroToken.value, -1)
// 		}
// 	}
// 	return tokens
// }
