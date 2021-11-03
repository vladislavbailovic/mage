package typedefs

const (
	TOKEN_MACRO_DFN_OPEN TokenType = iota + TOKEN_COMMAND_CLOSE + 1
	TOKEN_MACRO_DFN_CLOSE

	TOKEN_MACRO_CALL_OPEN
	TOKEN_MACRO_CALL_CLOSE
)

type SourcePosition struct {
	File string
	Line int
	Char int
}

type TaskDefinition struct {
	Pos            SourcePosition
	Name           string
	NormalizedName string
	Dependencies   []string
	Commands       []string
}

type MacroDefinition struct {
	Pos    SourcePosition
	Name   string
	Tokens []Token
}
