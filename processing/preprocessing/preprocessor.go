package preprocessing

// Preprocesses the list of tokens by expanding macros.

import (
	"fmt"
	"mage/debug"
	"mage/processing/tokenizing"
	"mage/shell"
	"mage/typedefs"
)

const MACRO_EXPANSION_RECURSE_LIMIT = 10

func Preprocess(tokens []typedefs.Token) ([]typedefs.Token, error) {
	proc := newPreprocessor(tokens)

	err := proc.doIncludes()
	if err != nil {
		return nil, err
	}

	err = proc.doMacros()
	if err != nil {
		return nil, err
	}

	return proc.tokens, nil
}

type preprocessor struct {
	tokens         []typedefs.Token
	macros         map[string]typedefs.MacroDefinition
	includes       []string
	shellcalls     []string
	expansionDepth int
	currentPos     int
}

func newPreprocessor(tokens []typedefs.Token) *preprocessor {
	return &preprocessor{
		tokens,
		map[string]typedefs.MacroDefinition{},
		[]string{},
		[]string{},
		MACRO_EXPANSION_RECURSE_LIMIT,
		0}
}

func (p *preprocessor) run() error {
	err := p.doIncludes()
	if err != nil {
		return err
	}
	return p.doMacros()
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
	for p.currentPos < len(p.tokens) {
		if p.tokens[p.currentPos].Kind != kind {
			p.currentPos += 1
			continue
		}
		return nil
	}
	return fmt.Errorf("unable to find [%v]", debug.GetTokenType(kind))
}

func (p *preprocessor) next() error {
	if p.currentPos < len(p.tokens) {
		p.currentPos += 1
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
			continue
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

	for p.currentPos < len(p.tokens) {
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
			p.next()

			macro, ok := p.macros[macroName]
			if !ok {
				return p.tokenError(fmt.Sprintf("unknown macro: [%v]", macroName))
			}

			for _, tk := range macro.Tokens {
				result = append(result, tk)
			}
			continue
		}
		result = append(result, p.current())
		p.next()
	}

	p.tokens = result[:]
	return nil
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

func includeFile(filepath string, relativeToSource string) ([]typedefs.Token, error) {
	relpath := shell.PathRelativeTo(filepath, relativeToSource)
	lines, err := shell.LoadFile(relpath)
	if err != nil {
		return nil, err
	}
	tkn := tokenizing.NewTokenizer(relpath, lines)
	return tkn.Tokenize(), nil
}
