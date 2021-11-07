package typedefs

type Task interface {
	GetMilestone() Epoch
	GetName() string
	GetCommands() []string
}
