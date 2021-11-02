package processor

import (
	"fmt"
	"mage/typedefs"
)

func preprocess(tokens []token) ([]typedefs.MacroDefinition, error) {
	result := []typedefs.MacroDefinition{}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].kind != typedefs.TOKEN_MACRO_DFN_OPEN {
			continue
		}
		i += 1
		nameToken := tokens[i]
		if nameToken.kind != typedefs.TOKEN_WORD {
			return nil, fmt.Errorf(
				"ERROR %s %d %d: macro name missing",
				nameToken.pos.File,
				nameToken.pos.Line,
				nameToken.pos.Char,
			)
		}

		valueTokens := []string{}
		for j := i; j < len(tokens); j++ {
			if tokens[j].kind == typedefs.TOKEN_MACRO_DFN_CLOSE {
				break
			}
			if tokens[j].kind != typedefs.TOKEN_WORD && tokens[j].kind != typedefs.TOKEN_MACRO_CALL_OPEN && tokens[j].kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
				return nil, fmt.Errorf(
					"ERROR %s %d %d: unexpected macro content; only words and macro calls are allowed but we got %v",
					tokens[j].pos.File,
					tokens[j].pos.Line,
					tokens[j].pos.Char,
					toktype(tokens[j].kind),
				)
			}
			valueTokens = append(valueTokens, tokens[j].value)
		}

		i += len(valueTokens)
	}

	return result, nil
}

func toktype(kind typedefs.TokenType) string {
	switch(kind) {
	case typedefs.TOKEN_WORD:
		return "word"
	case typedefs.TOKEN_MACRO_DFN_OPEN:
		return "macro dfn open"
	case typedefs.TOKEN_MACRO_DFN_CLOSE:
		return "macro dfn CLOSE"
	case typedefs.TOKEN_MACRO_CALL_OPEN:
		return "macro call open"
	case typedefs.TOKEN_MACRO_CALL_CLOSE:
		return "macro call CLOSE"
	case typedefs.TOKEN_RULE_OPEN:
		return "rule open"
	case typedefs.TOKEN_RULE_CLOSE:
		return "rule close"
	case typedefs.TOKEN_COMMAND_OPEN:
		return "command open"
	case typedefs.TOKEN_COMMAND_CLOSE:
		return "command close"
	}
	return "wut"
}
