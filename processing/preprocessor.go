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

type preprocessor struct {
	tokens         []typedefs.Token
	macros         map[string]typedefs.MacroDefinition
	includes       []string
	shellcalls     []string
	expansionDepth int
	currentPos     int
}

func NewPreprocessor(tokens []typedefs.Token) *preprocessor {
	return &preprocessor{tokens, map[string]typedefs.MacroDefinition{}, []string{}, []string{}, MACRO_EXPANSION_RECURSE_LIMIT, 0}
}

func (p *preprocessor) run() error {
	err := p.doIncludes()
	if err != nil {
		return err
	}
	return p.doMacros()
}

func preprocess(tokens []typedefs.Token) ([]typedefs.Token, error) {
	combined, err := preprocessIncludes(tokens)
	if err != nil {
		return nil, err
	}
	return preprocessMacros(combined)
}

func (p *preprocessor) doIncludes() error {
	for depth := 0; depth < p.expansionDepth; depth++ {
		changed := false
		var err error
		p.reset()
		for p.nextType(typedefs.TOKEN_INCLUDE_CALL_OPEN) == nil {
			at := p.currentPos
			if p.next() != nil {
				return p.tokenError("unfinished include")
			}

			include := p.current()
			if include.Kind != typedefs.TOKEN_WORD {
				return p.tokenError("include can only have words")
			}

			changed, err = p.includeFile(include, at)
			if err != nil {
				return err
			}
			break
		}
		if !changed {
			return nil
		}
	}
	return fmt.Errorf("exceeded includes recursion")
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

func (p *preprocessor) includeFile(from typedefs.Token, at int) (bool, error) {
	if p.nextType(typedefs.TOKEN_INCLUDE_CALL_CLOSE) != nil {
		return false, p.tokenErrorAt(at, "unfinished include")
	}
	end := p.currentPos + 1

	loaded, err := includeFile(from.Value, from.Pos.File)
	if err != nil {
		return false, err
	}
	if len(loaded) == 0 {
		return false, nil
	}

	result := []typedefs.Token{}
	result = append(result, p.tokens[:at]...)
	result = append(result, loaded...)
	result = append(result, p.tokens[end:]...)

	p.tokens = result[:]
	return true, nil
}

func (p *preprocessor) nextType(kind typedefs.TokenType) error {
	for p.currentPos < len(p.tokens)-1 {
		if p.tokens[p.currentPos].Kind != kind {
			p.currentPos++
			continue
		}
		return nil
	}
	return fmt.Errorf("unable to find [%v]", debug.GetTokenType(kind))
}

func (p *preprocessor) next() error {
	if p.currentPos < len(p.tokens) {
		p.currentPos++
		return nil
	}
	return fmt.Errorf("no more tokens")
}

func (p preprocessor) current() typedefs.Token {
	return p.tokens[p.currentPos]
}

func (p *preprocessor) reset() {
	p.currentPos = 0
}

func (p preprocessor) tokenErrorAt(pos int, msg string) error {
	return debug.TokenError(p.tokens[pos], msg)
}

func (p preprocessor) tokenError(msg string) error {
	return p.tokenErrorAt(p.currentPos, msg)
}

func (p *preprocessor) stripMacroDefinitions() error {
	result := []typedefs.Token{}
	for p.currentPos = 0; p.currentPos < len(p.tokens); p.currentPos++ {
		if p.current().Kind == typedefs.TOKEN_MACRO_DFN_OPEN {
			for p.current().Kind != typedefs.TOKEN_MACRO_DFN_CLOSE {
				p.currentPos++
			}
		}
		result = append(result, p.current())
	}

	p.tokens = result[:]
	return nil
}

func (p *preprocessor) doMacros() error {
	p.reset()
	err := p.doMacroDefinitions()
	if err != nil {
		return err
	}

	p.reset()
	err = p.stripMacroDefinitions()
	if err != nil {
		return err
	}

	p.reset()
	err = p.expandMacros()
	if err != nil {
		return err
	}

	return nil
}

func (p *preprocessor) expandMacros() error {
	result := []typedefs.Token{}
	for p.currentPos = 0; p.currentPos < len(p.tokens); p.currentPos++ {
		if p.current().Kind == typedefs.TOKEN_MACRO_CALL_OPEN {
			p.next()
			if p.current().Kind != typedefs.TOKEN_WORD {
				return p.tokenError(fmt.Sprintf(
					"expected word as macro name, got [%v]", debug.GetTokenType(p.current().Kind)))
			}

			macroName := p.current().Value

			p.next()
			if p.current().Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
				return p.tokenError("macro call not closed")
			}

			macro, ok := p.macros[macroName]
			if !ok {
				return p.tokenError(fmt.Sprintf("unknown macro: [%v]", macroName))
			}

			for _, tk := range macro.Tokens {
				result = append(result, tk)
			}

			p.next()
			continue
		}
		result = append(result, p.current())
	}

	p.tokens = result[:]
	return nil
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

func (p *preprocessor) doMacroDefinitions() error {
	dfns, err := getRawMacroDefinitions(p.tokens)
	if err != nil {
		return err
	}
	collection := newDfnCollection(dfns)

	for depth := 0; depth < p.expansionDepth; depth++ {
		changed, err := collection.process()
		if err != nil {
			return err
		}
		if !changed {
			break
		}
	}
	p.macros = collection.dfns
	return nil
}

type dfnCollection struct {
	dfns map[string]typedefs.MacroDefinition
}

func newDfnCollection(dfns map[string]typedefs.MacroDefinition) *dfnCollection {
	return &dfnCollection{dfns}
}

func (dc *dfnCollection) process() (bool, error) {
	changed := false
	for name, _ := range dc.dfns {
		dfnChanged, err := dc.processDfn(name)
		if err != nil {
			return changed, err
		}
		if dfnChanged {
			changed = true
		}
	}
	return changed, nil
}

func (dc *dfnCollection) processDfn(name string) (bool, error) {
	changed := false
	for idx, token := range dc.dfns[name].Tokens {
		if token.Kind != typedefs.TOKEN_MACRO_CALL_OPEN {
			continue
		}
		next := dc.dfns[name].Tokens[idx+1]
		if next.Kind != typedefs.TOKEN_WORD {
			return changed, debug.TokenError(
				next, fmt.Sprintf("macro call has to be a word, not [%v]", debug.GetTokenType(next.Kind)))
		}

		if "!" == string(next.Value[0]) {
			changed, err := dc.expandShellcodeIn(name, idx)
			if err != nil {
				return changed, err
			}
		} else {
			changed, err := dc.expandMacroIn(name, idx)
			if err != nil {
				return changed, err
			}
		}
	}
	return changed, nil
}

func (dc *dfnCollection) expandMacroIn(name string, start int) (bool, error) {
	dfn := dc.dfns[name]
	nameTok := dfn.Tokens[start+1]
	macro, ok := dc.dfns[nameTok.Value]
	if !ok {
		return false, debug.TokenError(nameTok, fmt.Sprintf("unknown token [%v]", nameTok.Value))
	}

	closeTok := dfn.Tokens[start+2]
	if closeTok.Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
		return false, debug.TokenError(closeTok, "macro call not closed off")
	}
	end := start + 3

	tks := append(dfn.Tokens[0:start], macro.Tokens...)
	tks = append(tks, dfn.Tokens[end:]...)
	dfn.Tokens = tks
	dc.dfns[name] = dfn

	return true, nil
}

func (dc *dfnCollection) expandShellcodeIn(name string, start int) (bool, error) {
	dfn := dc.dfns[name]
	nameTok := dfn.Tokens[start+1]
	end := 0
	cmd := []string{nameTok.Value[1:]}
	for i := start + 2; i < len(dfn.Tokens); i++ {
		// Can't mix shellcalls with nested macros.
		if dfn.Tokens[i].Kind != typedefs.TOKEN_MACRO_CALL_CLOSE {
			cmd = append(cmd, dfn.Tokens[i].Value)
			continue
		}
		end = i + 1
		break
	}
	if end == 0 {
		return false, debug.TokenError(nameTok, "shellcall macro call not closed off")
	}
	command := shell.NewCommand([]string{strings.Join(cmd, " ")})
	out, err := command.GetStdout()
	if err != nil {
		return false, debug.TokenError(nameTok, fmt.Sprintf("[%v] shellcall error: [%v]", cmd, err))
	}

	tks := append(dfn.Tokens[0:start], typedefs.Token{nameTok.Pos, typedefs.TOKEN_WORD, strings.TrimSpace(out)})
	tks = append(tks, dfn.Tokens[end:]...)
	dfn.Tokens = tks
	dc.dfns[name] = dfn

	return true, nil
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
