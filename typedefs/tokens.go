package typedefs

type TokenType int

const (
	TOKEN_WORD TokenType = iota
)

type Token struct {
	Pos   SourcePosition
	Kind  TokenType
	Value string
}
