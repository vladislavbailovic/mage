package processor

import (
	"mage/typedefs"
)

type token struct {
	pos typedefs.SourcePosition
	kind typedefs.TokenType
	value string
}

type tokenizerPosition struct {
	source string
	cursor int
	currentLine int
	currentChar int
}

func (tp *tokenizerPosition)advanceCursor(chr int) {
	tp.cursor += chr
}
func (tp *tokenizerPosition)advanceChar(chr int) {
	tp.currentChar += chr
}
func (tp *tokenizerPosition)advanceLine(chr int) {
	tp.currentLine += chr
	tp.currentChar = 1
}
func (tp *tokenizerPosition)advance(chr int) {
	tp.advanceCursor(chr)
	tp.advanceChar(chr)
}
func (tp tokenizerPosition)getPosition() typedefs.SourcePosition {
	return typedefs.SourcePosition{
		tp.source,
		tp.currentLine,
		tp.currentChar,
	}
}

func tokenize(file string, content string, initPos ...int) []token {
	content += " \n"
	currentLine := 1
	if len(initPos) > 0 {
		currentLine = initPos[0]
	}
	currentChar := 1
	if len(initPos) > 1 {
		currentChar = initPos[1]
	}
	allTokens := []token{}
	word := ""
	for pos := 0; pos < len(content) - 1; {

		if "\n" == string(content[pos]) {
			if len(word) > 0 {
				allTokens = append(allTokens, token{
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
			allTokens = append(allTokens, token{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_COMMAND_OPEN,
				"",
			})
			for _, tk := range tokenize(file, command, currentLine, currentChar) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, token{
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
			allTokens = append(allTokens, token{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_MACRO_DFN_OPEN,
				"",
			})
			for _, tk := range tokenize(file, macro, currentLine, currentChar) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, token{
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
				allTokens = append(allTokens, token{
					typedefs.SourcePosition{file, currentLine, currentChar + 1},
					typedefs.TOKEN_WORD,
					word,
				})
				word = ""
			}
			// Macro call
			pos += 2
			macro := consumeUntil(")", content, pos)
			allTokens = append(allTokens, token{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_MACRO_CALL_OPEN,
				"",
			})
			for _, tk := range tokenize(file, macro, currentLine, currentChar) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, token{
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
			allTokens = append(allTokens, token{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_RULE_OPEN,
				"",
			})
			for _, tk := range tokenize(file, name, currentLine, currentChar) {
				allTokens = append(allTokens, tk)
			}
			for _, tk := range tokenize(file, dependencies, currentLine, currentChar) {
				allTokens = append(allTokens, tk)
			}
			allTokens = append(allTokens, token{
				typedefs.SourcePosition{file, currentLine, currentChar + 1},
				typedefs.TOKEN_RULE_CLOSE,
				"",
			})
			pos += len(dependencies) + 1
			word = ""
			continue
		}

		if " " == string(content[pos]) && len(word) > 0 {
			allTokens = append(allTokens, token{
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

func filterTokens(tokens []token, expected typedefs.TokenType) []token {
	result := []token{}
	for _, tk := range tokens {
		if expected != tk.kind {
			continue
		}
		result = append(result, tk)
	}
	return result
}

// func transform(tokens []token) []token {
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

type tokenizer struct {
	tokens []token
	content string
	position *tokenizerPosition
}

func newTokenizer(source string, content string) *tokenizer {
	pos := tokenizerPosition{ source, 0, 1, 1 }
	tkn := tokenizer{
		content: content,
		position: &pos,
	}
	return &tkn
}

func newTokenizerAt(src string, content string, line int, char int) *tokenizer {
	tkn := newTokenizer(src, content)
	tkn.position.currentLine = line
	tkn.position.currentChar = char
	return tkn
}

func (t *tokenizer)addNewToken(kind typedefs.TokenType, value string) {
	t.addToken(token{
		t.position.getPosition(),
		kind,
		value,
	})
}

func (t *tokenizer)addToken(tk token) {
	t.tokens = append(t.tokens, tk)
}

func (t *tokenizer)getSubtokens(substr string) []token {
	subt := newTokenizer(t.position.source, substr)
	subt.position.advanceLine(t.position.currentLine)
	subt.position.advanceChar(t.position.currentChar)
	return subt.tokenize()
}

func (t *tokenizer)tokenize() []token {
	content := t.content + " \n"
	word := ""
	for t.position.cursor < len(content) - 1 {
		chr := string(content[t.position.cursor])

		if "\n" == chr {
			if len(word) > 0 {
				t.addNewToken(typedefs.TOKEN_WORD, word)
				word = ""
			}
			t.position.advance(1)
			t.position.advanceLine(1)
			continue
		}

		if "\t" == chr {
			// Command
			t.position.advanceCursor(1)
			command := consumeUntil("\n", content, t.position.cursor)
			t.addNewToken(typedefs.TOKEN_COMMAND_OPEN, "")
			for _, tk := range t.getSubtokens(command) {
				t.addToken(tk)
			}
			t.addNewToken(typedefs.TOKEN_COMMAND_CLOSE, "")
			t.position.advanceCursor(len(command))
			word = ""
			continue
		}

		kw := "macro"
		kwLen := len(kw)
		matchEnd := t.position.cursor + kwLen
		if len(content) >= matchEnd {
			match := string(content[t.position.cursor:matchEnd])
			if kw == match {
				// Macro definition
				t.position.advanceCursor(kwLen + 1)
				macro := consumeUntil("\n", content, t.position.cursor)
				t.addNewToken(typedefs.TOKEN_MACRO_DFN_OPEN, "")
				for _, tk := range t.getSubtokens(macro) {
					t.addToken(tk)
				}
				t.addNewToken(typedefs.TOKEN_MACRO_DFN_CLOSE, "")
				t.position.advanceCursor(len(macro))
				word = ""
				continue
			}
		}

		kw = "$("
		kwLen = len(kw)
		matchEnd = t.position.cursor + kwLen
		if len(content) >= matchEnd {
			match := string(content[t.position.cursor:matchEnd])
			if kw == match {
				if len(word) > 0 {
					t.addNewToken(typedefs.TOKEN_WORD, word)
					word = ""
				}
				// Macro call
				t.position.advanceCursor(2)
				macro := consumeUntil(")", content, t.position.cursor)
				t.addNewToken(typedefs.TOKEN_MACRO_CALL_OPEN, "")
				for _, tk := range t.getSubtokens(macro) {
					t.addToken(tk)
				}
				t.addNewToken(typedefs.TOKEN_MACRO_CALL_CLOSE, "")
				t.position.advance(len(macro)+1)
				word = ""
				continue
			}
		}

		// if ":" == string(content[pos]) {
		// 	// Rule definition
		// 	name := consumeBackUntil("\n", content, pos-1)
		// 	dependencies := consumeUntil("\n", content, pos+1)
		// 	allTokens = append(allTokens, token{
		// 		typedefs.SourcePosition{file, currentLine, currentChar + 1},
		// 		typedefs.TOKEN_RULE_OPEN,
		// 		"",
		// 	})
		// 	for _, tk := range tokenize(file, name, currentLine, currentChar) {
		// 		allTokens = append(allTokens, tk)
		// 	}
		// 	for _, tk := range tokenize(file, dependencies, currentLine, currentChar) {
		// 		allTokens = append(allTokens, tk)
		// 	}
		// 	allTokens = append(allTokens, token{
		// 		typedefs.SourcePosition{file, currentLine, currentChar + 1},
		// 		typedefs.TOKEN_RULE_CLOSE,
		// 		"",
		// 	})
		// 	pos += len(dependencies) + 1
		// 	word = ""
		// 	continue
		// }

		if " " == chr && len(word) > 0 {
			t.addNewToken(typedefs.TOKEN_WORD, word)
			word = ""
			t.position.advance(1)
			continue
		}

		word += chr

		t.position.advance(1)
	}
	return t.tokens
}
