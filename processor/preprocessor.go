package processor

import (
	"fmt"
	"mage/typedefs"
)

const MACRO_EXPANSION_RECURSE_LIMIT = 10


type macroDefinition struct {
	Pos typedefs.SourcePosition
	Name string
	tokens []token
}

func preprocess(tokens []token) ([]token, error) {
	mds, err := getMacroDefinitions(tokens)
	if err != nil {
		return nil, err
	}

	fmt.Println("----- before -----")
	dbgdefs(mds)

	// Prepare macro definitions by expanding calls
	recursionCounter := 0
	for recursionCounter < MACRO_EXPANSION_RECURSE_LIMIT {
		didReplacement := false
		for name, md := range mds {
			for idx, token := range md.tokens {
				if token.kind != typedefs.TOKEN_MACRO_CALL_OPEN {
					continue
				}
				nameTok := md.tokens[idx+1]
				if nameTok.kind != typedefs.TOKEN_WORD {
					panic("macro call has to be a word")
				}

				macro, ok := mds[nameTok.value]
				if !ok {
					panic("can't find token: " + nameTok.value)
				}

				closeTok := md.tokens[idx+2]
				if closeTok.kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
					panic("macro call not closed off")
				}

				tks := append(md.tokens[0:idx], macro.tokens...)
				tks = append(tks, md.tokens[idx+3:]...)
				md.tokens = tks
				mds[name] = md
				didReplacement = true
			}
		}
		recursionCounter++
		if !didReplacement {
			break
		}
	}

	fmt.Println("----- after -----")
	dbgdefs(mds)

	return tokens, nil
}

func dbgdefs(mds map[string]macroDefinition) {
	for n,m := range mds {
		fmt.Printf("- [%s]:\n", n)
		for _, t := range m.tokens {
			fmt.Printf("\t> [%s] (%s)\n", toktype(t.kind), t.value)
		}
	}
}

func getMacroDefinitions(tokens []token) (map[string]macroDefinition, error) {
	result := map[string]macroDefinition{}

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
		i += 1

		valueTokens := []token{}
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
			valueTokens = append(valueTokens, tokens[j])
		}

		result[nameToken.value] = macroDefinition{
			nameToken.pos,
			nameToken.value,
			valueTokens,
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
