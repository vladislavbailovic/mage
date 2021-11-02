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

func (t *tokenizer)tokenizeSubstring(substr string) []token {
	subt := newTokenizer(t.position.source, substr)
	subt.position.advanceLine(t.position.currentLine)
	subt.position.advanceChar(t.position.currentChar)
	return subt.tokenize()
}

func (t *tokenizer)processNewline(word string) string {
	if len(word) > 0 {
		t.addNewToken(typedefs.TOKEN_WORD, word)
		word = ""
	}
	t.position.advance(1)
	t.position.advanceLine(1)
	return word
}

func (t *tokenizer)processCommand(word string) string {
	// Command
	content := t.content + " \n"
	t.position.advanceCursor(1)
	command := consumeUntil("\n", content, t.position.cursor)
	t.addNewToken(typedefs.TOKEN_COMMAND_OPEN, "")
	for _, tk := range t.tokenizeSubstring(command) {
		t.addToken(tk)
	}
	t.addNewToken(typedefs.TOKEN_COMMAND_CLOSE, "")
	t.position.advanceCursor(len(command))
	return ""
}

func (t *tokenizer)processMacroDfn(word string) string {
	content := t.content + " \n"
	kwLen := len("macro")
	// Macro definition
	t.position.advanceCursor(kwLen + 1)
	macro := consumeUntil("\n", content, t.position.cursor)
	t.addNewToken(typedefs.TOKEN_MACRO_DFN_OPEN, "")
	for _, tk := range t.tokenizeSubstring(macro) {
		t.addToken(tk)
	}
	t.addNewToken(typedefs.TOKEN_MACRO_DFN_CLOSE, "")
	t.position.advanceCursor(len(macro))
	return ""
}

func (t *tokenizer)processMacroCall(word string) string {
	content := t.content + " \n"
	if len(word) > 0 {
		t.addNewToken(typedefs.TOKEN_WORD, word)
	}
	// Macro call
	t.position.advanceCursor(2)
	macro := consumeUntil(")", content, t.position.cursor)
	t.addNewToken(typedefs.TOKEN_MACRO_CALL_OPEN, "")
	for _, tk := range t.tokenizeSubstring(macro) {
		t.addToken(tk)
	}
	t.addNewToken(typedefs.TOKEN_MACRO_CALL_CLOSE, "")
	t.position.advance(len(macro)+1)
	return ""
}

func (t *tokenizer)tokenize() []token {
	content := t.content + " \n"
	word := ""
	for t.position.cursor < len(content) - 1 {
		chr := string(content[t.position.cursor])

		if "\n" == chr {
			word = t.processNewline(word)
			continue
		}

		if "\t" == chr {
			word = t.processCommand(word)
			continue
		}

		kw := "macro"
		kwLen := len(kw)
		matchEnd := t.position.cursor + kwLen
		if len(content) >= matchEnd {
			match := string(content[t.position.cursor:matchEnd])
			if kw == match {
				word = t.processMacroDfn(word)
				continue
			}
		}

		kw = "$("
		kwLen = len(kw)
		matchEnd = t.position.cursor + kwLen
		if len(content) >= matchEnd {
			match := string(content[t.position.cursor:matchEnd])
			if kw == match {
				word = t.processMacroCall(word)
				continue
			}
		}

		if ":" == chr {
			// Rule definition
			name := consumeBackUntil("\n", content, t.position.cursor-1)
			dependencies := consumeUntil("\n", content, t.position.cursor+1)
			t.addNewToken(typedefs.TOKEN_RULE_OPEN, "")
			for _, tk := range t.tokenizeSubstring(name) {
				t.addToken(tk)
			}
			for _, tk := range t.tokenizeSubstring(dependencies) {
				t.addToken(tk)
			}
			t.addNewToken(typedefs.TOKEN_RULE_CLOSE, "")
			t.position.advance(len(dependencies)+1)
			word = ""
			continue
		}

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

func (t tokenizer)filter(expected typedefs.TokenType) []token {
	result := []token{}
	for _, tk := range t.tokens {
		if expected != tk.kind {
			continue
		}
		result = append(result, tk)
	}
	return result
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
