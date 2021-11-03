package debug

import (
	"fmt"
	"mage/typedefs"
)

func MacroDefinitions(mds map[string]typedefs.MacroDefinition) {
	for n, m := range mds {
		fmt.Printf("- [%s]:\n", n)
		for _, t := range m.Tokens {
			fmt.Printf("\t> [%s] (%s)\n", GetTokenType(t.Kind), t.Value)
		}
	}
}

func TaskDefinitions(tds map[string]typedefs.TaskDefinition) {
	for n, t := range tds {
		fmt.Printf(
			"[%v] (%v), at %s %d:%d:\n",
			n, t.Name,
			t.Pos.File,
			t.Pos.Line,
			t.Pos.Char,
		)
		fmt.Printf("  deps: ")
		for di, dep := range t.Dependencies {
			fmt.Printf("[%d: {%v}], ", di, dep)
		}
		fmt.Printf("\n  cmds: ")
		for dc, cmd := range t.Commands {
			fmt.Printf("[%d: {%v}], ", dc, cmd)
		}
		fmt.Println("")
	}
}

func Tokens(tokens []typedefs.Token) {
	for idx, token := range tokens {
		fmt.Printf("%d: (%s) ", idx, GetTokenType(token.Kind))
		fmt.Printf(
			"[%s, %d:%d] [%v]\n",
			token.Pos.File,
			token.Pos.Line,
			token.Pos.Char,
			token.Value,
		)
	}
}

func GetTokenType(kind typedefs.TokenType) string {
	switch kind {
	case typedefs.TOKEN_WORD:
		return "word"
	case typedefs.TOKEN_MACRO_DFN_OPEN:
		return "macro dfn open"
	case typedefs.TOKEN_MACRO_DFN_CLOSE:
		return "macro dfn CLOSE"
	case typedefs.TOKEN_MACRO_CALL_OPEN:
		return "macro call open"
	case typedefs.TOKEN_MACRO_CALL_CLOSE:
		return "macro call close"
	case typedefs.TOKEN_RULE_OPEN:
		return "rule open"
	case typedefs.TOKEN_RULE_CLOSE:
		return "rule close"
	case typedefs.TOKEN_COMMAND_OPEN:
		return "command open"
	case typedefs.TOKEN_COMMAND_CLOSE:
		return "command close"
	case typedefs.TOKEN_INCLUDE_CALL_OPEN:
		return "include open"
	case typedefs.TOKEN_INCLUDE_CALL_CLOSE:
		return "include close"
	}
	return "UNKNOWN"
}

func TokenError(tk typedefs.Token, msg string) error {
	return fmt.Errorf(
		"ERROR: %s %d %d (%v): %v",
		tk.Pos.File,
		tk.Pos.Line,
		tk.Pos.Char,
		GetTokenType(tk.Kind),
		msg,
	)

}
