package tokenizing

// Builds a list of tokens to be used further in the
// processing pipeline.

import (
	"mage/typedefs"
)

type tokenizer struct {
	tokens   []typedefs.Token
	content  string
	position *tokenizerPosition
}

func NewTokenizer(source string, content string) *tokenizer {
	pos := tokenizerPosition{source, 0, 1, 1}
	tkn := tokenizer{
		content:  content,
		position: &pos,
	}
	return &tkn
}

func (t *tokenizer) addNewToken(kind typedefs.TokenType, value string) {
	t.addToken(typedefs.Token{
		t.position.getPosition(),
		kind,
		value,
	})
}

func (t *tokenizer) addToken(tk typedefs.Token) {
	t.tokens = append(t.tokens, tk)
}

func (t *tokenizer) tokenizeSubstring(substr string) []typedefs.Token {
	subt := NewTokenizer(t.position.source, substr)
	subt.position.currentLine = t.position.currentLine
	subt.position.currentChar = t.position.currentChar
	return subt.Tokenize()
}

func (t *tokenizer) processNewline(word string) string {
	if len(word) > 0 {
		t.addNewToken(typedefs.TOKEN_WORD, word)
		word = ""
	}
	t.position.advanceCursor(1)
	t.position.advanceLine(1)
	return word
}

func (t *tokenizer) processCommand(word string) string {
	// Command
	content := t.content + " \n"
	t.position.advance(1)
	command := consumeUntil("\n", content, t.position.cursor)
	t.addNewToken(typedefs.TOKEN_COMMAND_OPEN, "")
	for _, tk := range t.tokenizeSubstring(command) {
		t.addToken(tk)
	}
	t.addNewToken(typedefs.TOKEN_COMMAND_CLOSE, "")
	t.position.advance(len(command))
	return ""
}

func (t *tokenizer) processMacroDfn(word string) string {
	content := t.content + " \n"
	kwLen := len(":macro")
	// Macro definition
	t.position.advance(kwLen + 1)
	macro := consumeUntil("\n", content, t.position.cursor)
	t.addNewToken(typedefs.TOKEN_MACRO_DFN_OPEN, "")
	for _, tk := range t.tokenizeSubstring(macro) {
		t.addToken(tk)
	}
	t.addNewToken(typedefs.TOKEN_MACRO_DFN_CLOSE, "")
	t.position.advance(len(macro))
	return ""
}

func (t *tokenizer) processMacroCall(word string) string {
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
	t.position.advance(len(macro) + 1)
	return ""
}

func (t *tokenizer) processRule(word string) string {
	content := t.content + " \n"
	// Rule definition
	name := consumeBackUntil("\n", content, t.position.cursor-1)
	dependencies := consumeUntil("\n", content, t.position.cursor+1)
	t.addNewToken(typedefs.TOKEN_RULE_OPEN, "")
	for _, tk := range t.tokenizeSubstring(name) {
		t.addToken(tk)
	}
	t.position.advanceChar(len(name) + 1)
	for _, tk := range t.tokenizeSubstring(dependencies) {
		t.addToken(tk)
	}
	t.addNewToken(typedefs.TOKEN_RULE_CLOSE, "")
	t.position.advance(len(dependencies) + 1)
	return ""
}

func (t *tokenizer) processIncludeCall(word string) string {
	content := t.content + " \n"
	kwLen := len(":include")
	// Include definition
	t.position.advance(kwLen + 1)
	include := consumeUntil("\n", content, t.position.cursor)
	t.addNewToken(typedefs.TOKEN_INCLUDE_CALL_OPEN, "")
	for _, tk := range t.tokenizeSubstring(include) {
		t.addToken(tk)
	}
	t.addNewToken(typedefs.TOKEN_INCLUDE_CALL_CLOSE, "")
	t.position.advance(len(include))
	return ""
}

func (t *tokenizer) Tokenize() []typedefs.Token {
	content := t.content + " \n"
	word := ""
	for t.position.cursor < len(content)-1 {
		chr := string(content[t.position.cursor])

		if "\t" == chr {
			word = t.processCommand(word)
			continue
		}

		kw := ":macro"
		kwLen := len(kw)
		matchEnd := t.position.cursor + kwLen
		if len(content) >= matchEnd {
			match := string(content[t.position.cursor:matchEnd])
			if kw == match {
				word = t.processMacroDfn(word)
				continue
			}
		}

		kw = ":include"
		kwLen = len(kw)
		matchEnd = t.position.cursor + kwLen
		if len(content) >= matchEnd {
			match := string(content[t.position.cursor:matchEnd])
			if kw == match {
				word = t.processIncludeCall(word)
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
			word = t.processRule(word)
			continue
		}

		if "\n" == chr {
			word = t.processNewline(word)
			continue
		}

		if " " == chr && len(word) > 0 {
			t.addNewToken(typedefs.TOKEN_WORD, word)
			t.position.advanceCursor(1)
			t.position.advanceChar(len(word) + 1)
			word = ""
			continue
		}

		word += chr

		t.position.advanceCursor(1)
	}
	return t.tokens
}

func (t tokenizer) filter(expected typedefs.TokenType) []typedefs.Token {
	result := []typedefs.Token{}
	for _, tk := range t.tokens {
		if expected != tk.Kind {
			continue
		}
		result = append(result, tk)
	}
	return result
}

func consumeUntil(what string, source string, pos int) string {
	item := ""
	for i := pos; i < len(source)-1; i++ {
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
