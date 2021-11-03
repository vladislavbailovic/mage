package processor

import (
	"fmt"
	"mage/typedefs"
)

const MACRO_EXPANSION_RECURSE_LIMIT = 10

func preprocess(tokens []typedefs.Token) ([]typedefs.Token, error) {
	macros, err := getMacroDefinitions(tokens)
	if err != nil {
		return nil, err
	}
	result := []typedefs.Token{}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Kind == typedefs.TOKEN_MACRO_DFN_OPEN {
			// SKip over macro definitions, already have those
			for tokens[i].Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
				i += 1
			}
			i += 1
			continue
		}

		if tokens[i].Kind == typedefs.TOKEN_MACRO_CALL_OPEN {
			// Expand macro calls
			i += 1
			if tokens[i].Kind != typedefs.TOKEN_WORD {
				panic("expected word")
			}
			if tokens[i+1].Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
				panic("macro not closed")
			}
			macroName := tokens[i].Value
			macro, ok := macros[macroName]
			if !ok {
				panic("unknown macro")
			}
			for _, tk := range macro.Tokens {
				result = append(result, tk)
			}
			i += 1
			continue
		}

		result = append(result, tokens[i])
	}

	return result, nil
}

func getMacroDefinitions(tokens []typedefs.Token) (map[string]typedefs.MacroDefinition, error) {
	dfns, err := getRawMacroDefinitions(tokens)
	if err != nil {
		return nil, err
	}

	// Prepare macro definitions by expanding calls
	recursionCounter := 0
	for recursionCounter < MACRO_EXPANSION_RECURSE_LIMIT {
		didReplacement := false
		for name, md := range dfns {
			for idx, token := range md.Tokens {
				if token.Kind != typedefs.TOKEN_MACRO_CALL_OPEN {
					continue
				}
				nameTok := md.Tokens[idx+1]
				if nameTok.Kind != typedefs.TOKEN_WORD {
					panic("macro call has to be a word")
				}

				macro, ok := dfns[nameTok.Value]
				if !ok {
					panic("can't find token: " + nameTok.Value)
				}

				closeTok := md.Tokens[idx+2]
				if closeTok.Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
					panic("macro call not closed off")
				}

				tks := append(md.Tokens[0:idx], macro.Tokens...)
				tks = append(tks, md.Tokens[idx+3:]...)
				md.Tokens = tks
				dfns[name] = md
				didReplacement = true
			}
		}
		recursionCounter++
		if !didReplacement {
			break
		}
	}

	return dfns, nil
}

func getRawMacroDefinitions(tokens []typedefs.Token) (map[string]typedefs.MacroDefinition, error) {
	result := map[string]typedefs.MacroDefinition{}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Kind != typedefs.TOKEN_MACRO_DFN_OPEN {
			continue
		}
		i += 1
		nameToken := tokens[i]
		if nameToken.Kind != typedefs.TOKEN_WORD {
			return nil, fmt.Errorf(
				"ERROR %s %d %d: macro name missing",
				nameToken.Pos.File,
				nameToken.Pos.Line,
				nameToken.Pos.Char,
			)
		}
		i += 1

		valueTokens := []typedefs.Token{}
		for j := i; j < len(tokens); j++ {
			if tokens[j].Kind == typedefs.TOKEN_MACRO_DFN_CLOSE {
				break
			}
			if tokens[j].Kind != typedefs.TOKEN_WORD && tokens[j].Kind != typedefs.TOKEN_MACRO_CALL_OPEN && tokens[j].Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
				return nil, fmt.Errorf(
					"ERROR %s %d %d: unexpected macro content; only words and macro calls are allowed but we got %v",
					tokens[j].Pos.File,
					tokens[j].Pos.Line,
					tokens[j].Pos.Char,
					toktype(tokens[j].Kind),
				)
			}
			valueTokens = append(valueTokens, tokens[j])
		}

		result[nameToken.Value] = typedefs.MacroDefinition{
			nameToken.Pos,
			nameToken.Value,
			valueTokens,
		}

		i += len(valueTokens)
	}

	return result, nil
}

