package typedefs

const (
	TOKEN_MACRO_DFN_OPEN TokenType = iota + TOKEN_COMMAND_CLOSE + 1
	TOKEN_MACRO_DFN_CLOSE
	TOKEN_MACRO_DFN
	TOKEN_MACRO_CALL_OPEN
	TOKEN_MACRO_CALL_CLOSE
	TOKEN_MACRO_CALL
)
