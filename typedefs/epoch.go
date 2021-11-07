package typedefs

type Epoch int64

type Milestone interface {
	GetMilestone() Epoch
}

type Record struct {
	Name      string
	Timestamp Epoch
}
