package processor

import (
	"strings"
	"testing"
)

func Test_Process(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	rawTokens := tkn.tokenize()
	tokens, _ := preprocess(rawTokens)
	dfns, err := process(tokens)
	if err != nil {
		t.Fatalf("expected processing to succeed")
	}
	debugTaskDefinitions(dfns)
}
