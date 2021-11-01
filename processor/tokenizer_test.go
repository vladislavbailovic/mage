package processor

import (
	"testing"
	"strings"

	"mage/typedefs"
)

func Test_Tokenize(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tokens := tokenize("test", strings.Join(lines, "\n"))
	expected := 63
	if expected != len(tokens) {
		t.Fatalf("expected %d tokens, but got %d", expected, len(tokens))
	}
	tokMacros := filterTokens(tokens, typedefs.TOKEN_MACRO_DFN_OPEN)
	if 5 != len(tokMacros) {
		t.Fatalf("there should be 5 macros, not %d", len(tokMacros))
	}
	tokRules := filterTokens(tokens, typedefs.TOKEN_RULE_OPEN)
	if 2 != len(tokRules) {
		t.Fatalf("there should be 2 rule dfns, not %d", len(tokRules))
	}
	for _, tk := range tokens {
		t.Log(tk)
	}
}

// func Test_Transform(t *testing.T) {
// 	lines, _ := loadFile("../fixtures/macro.mg")
// 	tokens := tokenize("test", strings.Join(lines, "\n"))
// 	transformedTokens := transform(tokens)
// 	for _, tk := range transformedTokens {
// 		t.Log(tk)
// 	}
// }
