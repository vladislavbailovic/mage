package typedefs

const (
	TOKEN_RULE TokenType = iota + TOKEN_WORD + 1
	TOKEN_RULE_OPEN
	TOKEN_RULE_CLOSE
	TOKEN_COMMAND
	TOKEN_COMMAND_OPEN
	TOKEN_COMMAND_CLOSE
)
