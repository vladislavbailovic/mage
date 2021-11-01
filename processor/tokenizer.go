package processor

import (
	"fmt"
	"mage/typedefs"
)

type xtoken struct {
	pos typedefs.SourcePosition
	kind typedefs.TokenType
	value string
}

func tokenize(file string, content string) []xtoken {
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
			pos += 1
			command := consumeUntil("\n", content, pos)
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_COMMAND_OPEN,
				"",
			})
			for _, tk := range tokenize(file, command) {
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
			currentChar += 5
			macro := consumeUntil("\n", content, pos)
			fmt.Printf("\tmacro [%s]\n", macro)
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_MACRO_DFN_OPEN,
				"",
			})
			for _, tk := range tokenize(file, macro) {
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
			for _, tk := range tokenize(file, macro) {
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
			for _, tk := range tokenize(file, name) {
				allTokens = append(allTokens, tk)
			}
			for _, tk := range tokenize(file, dependencies) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, xtoken{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_RULE_CLOSE,
				"",
			})
			pos += len(dependencies) + 1
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
