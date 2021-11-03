package processing

import (
	"strings"
	"testing"
)

func Test_Includes(t *testing.T) {
	lines, _ := loadFile("../fixtures/includes.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	rawTokens := tkn.tokenize()
	debugTokens(rawTokens)
}
