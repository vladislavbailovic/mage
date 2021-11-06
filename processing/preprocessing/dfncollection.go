package preprocessing

import (
	"fmt"
	"mage/debug"
	"mage/shell"
	"mage/typedefs"
	"strings"
)

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
	somethingChanged := false
	for idx, token := range dc.dfns[name].Tokens {
		if token.Kind != typedefs.TOKEN_MACRO_CALL_OPEN {
			continue
		}
		next := dc.dfns[name].Tokens[idx+1]
		if next.Kind != typedefs.TOKEN_WORD {
			return somethingChanged, debug.TokenError(
				next, fmt.Sprintf("macro call has to be a word, not [%v]", debug.GetTokenType(next.Kind)))
		}

		if "!" == string(next.Value[0]) {
			changed, err := dc.expandShellcodeIn(name, idx)
			if err != nil {
				return changed, err
			}
			if changed {
				somethingChanged = true
			}
		} else {
			changed, err := dc.expandMacroIn(name, idx)
			if err != nil {
				return changed, err
			}
			if changed {
				somethingChanged = true
			}
		}
	}
	return somethingChanged, nil
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
