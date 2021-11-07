package typedefs

const (
	TOKEN_RULE_OPEN TokenType = iota + TOKEN_WORD + 1
	TOKEN_RULE_CLOSE

	TOKEN_COMMAND_OPEN
	TOKEN_COMMAND_CLOSE
)

type Task interface {
	GetMilestone() Epoch
	RecordTime()
	GetName() string
	GetCommands() []string
}
