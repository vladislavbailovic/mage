package typedefs

type Epoch int64

type Milestone interface {
	GetMilestone() Epoch
}
