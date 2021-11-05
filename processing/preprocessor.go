package processing

// Preprocesses the list of tokens by expanding macros.

import (
	"fmt"
	"mage/debug"
	"mage/shell"
	"mage/typedefs"
	"strings"
)

const MACRO_EXPANSION_RECURSE_LIMIT = 10

func preprocess(tokens []typedefs.Token) ([]typedefs.Token, error) {
	combined, err := preprocessIncludes(tokens)
	if err != nil {
		return nil, err
	}
	return preprocessMacros(combined)
}

func preprocessIncludes(tokens []typedefs.Token) ([]typedefs.Token, error) {
	for safety := 0; safety < MACRO_EXPANSION_RECURSE_LIMIT; safety++ {
		changed := false
		result := []typedefs.Token{}
		for i := 0; i < len(tokens); i++ {
			if tokens[i].Kind != typedefs.TOKEN_INCLUDE_CALL_OPEN {
				continue
			}
			start := i
			end := i + 3
			i += 1
			if tokens[i].Kind != typedefs.TOKEN_WORD {
				return nil, debug.TokenError(tokens[i], "include can only have words")
			}
			filepath := tokens[i].Value
			loadedTokens, err := includeFile(filepath, tokens[i].Pos.File)
			if err != nil {
				return nil, err
			}
			result = append(result, tokens[:start]...)
			result = append(result, loadedTokens...)
			result = append(result, tokens[end:]...)
			changed = true
			break
		}
		if !changed {
			return tokens, nil
		}
		tokens = result[:]
	}
	return nil, fmt.Errorf("exceeded includes recursion")
}

func preprocessMacros(tokens []typedefs.Token) ([]typedefs.Token, error) {
	macros, err := getMacroDefinitions(tokens)
	if err != nil {
		return nil, err
	}
	result := []typedefs.Token{}

	for i := 0; i < len(tokens); i++ {
		if tokens[i].Kind == typedefs.TOKEN_MACRO_DFN_OPEN {
			// SKip over macro definitions, already have those
			for tokens[i].Kind != typedefs.TOKEN_MACRO_DFN_CLOSE {
				i += 1
			}
			continue
		}

		if tokens[i].Kind == typedefs.TOKEN_MACRO_CALL_OPEN {
			// Expand macro calls
			i += 1
			if tokens[i].Kind != typedefs.TOKEN_WORD {
				return nil, debug.TokenError(tokens[i], fmt.Sprintf("expected word as macro name, got [%v]", debug.GetTokenType(tokens[i].Kind)))
			}
			if tokens[i+1].Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
				return nil, debug.TokenError(tokens[i], "macro call not closed")
			}
			macroName := tokens[i].Value
			macro, ok := macros[macroName]
			if !ok {
				return nil, debug.TokenError(tokens[i], fmt.Sprintf("unknown macro: [%v]", macroName))
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

func includeFile(filepath string, relativeToSource string) ([]typedefs.Token, error) {
	relpath := shell.PathRelativeTo(filepath, relativeToSource)
	lines, err := shell.LoadFile(relpath)
	if err != nil {
		return nil, err
	}
	tkn := newTokenizer(relpath, lines)
	return tkn.tokenize(), nil
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
					return nil, debug.TokenError(nameTok, fmt.Sprintf("macro call has to be a word, not [%v]", debug.GetTokenType(nameTok.Kind)))
				}

				// @TODO this is horrible, refactor
				var newTokens []typedefs.Token
				var endIndex int
				if "!" != string(nameTok.Value[0]) {
					// Normal macro
					macro, ok := dfns[nameTok.Value]
					if !ok {
						return nil, debug.TokenError(nameTok, fmt.Sprintf("unknown token [%v]", nameTok.Value))
					}

					closeTok := md.Tokens[idx+2]
					if closeTok.Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
						return nil, debug.TokenError(closeTok, "macro call not closed off")
					}
					endIndex = idx + 3
					newTokens = macro.Tokens
				} else {
					// Starts with "!" - shellcall macro
					cmd := []string{nameTok.Value[1:]}
					for i := idx + 2; i < len(md.Tokens); i++ {
						// Can't mix shellcalls with nested macros.
						if md.Tokens[i].Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
							cmd = append(cmd, md.Tokens[i].Value)
							continue
						}
						endIndex = i + 1
						break
					}
					if endIndex == 0 {
						return nil, debug.TokenError(nameTok, "shellcall macro call not closed off")
					}
					command := shell.NewCommand([]string{strings.Join(cmd, " ")})
					out, err := command.GetStdout()
					if err != nil {
						return nil, debug.TokenError(nameTok, fmt.Sprintf("[%v] shellcall error: [%v]", cmd, err))
					}
					newTokens = []typedefs.Token{
						typedefs.Token{nameTok.Pos, typedefs.TOKEN_WORD, strings.TrimSpace(out)},
					}
				}

				tks := append(md.Tokens[0:idx], newTokens...)
				tks = append(tks, md.Tokens[endIndex:]...)
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
			return nil, debug.TokenError(nameToken, "macro name missing")
		}
		i += 1

		valueTokens := []typedefs.Token{}
		for j := i; j < len(tokens); j++ {
			if tokens[j].Kind == typedefs.TOKEN_MACRO_DFN_CLOSE {
				break
			}
			if tokens[j].Kind != typedefs.TOKEN_WORD && tokens[j].Kind != typedefs.TOKEN_MACRO_CALL_OPEN && tokens[j].Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
				return nil, debug.TokenError(tokens[j], fmt.Sprintf("unexpected macro content; only words and macro calls are allowed but we got %v", debug.GetTokenType(tokens[j].Kind)))
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
