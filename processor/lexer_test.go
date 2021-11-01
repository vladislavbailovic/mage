package processor

import (
	"testing"
	"strings"
)

func Test_LexFile(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tokens := xlex("test", strings.Join(lines, "\n"))
	// if 10 != len(tokens) {
	// 	t.Fatalf("expected 10 tokens")
	// }
	for _, tk := range tokens {
		t.Log(tk)
	}
}

// func Test_Transform(t *testing.T) {
// 	lines, _ := loadFile("../fixtures/macro.mg")
// 	tokens := xlex("test", strings.Join(lines, "\n"))
// 	transformedTokens := transform(tokens)
// 	for _, tk := range transformedTokens {
// 		t.Log(tk)
// 	}
// }
