package epoch

import (
	"mage/typedefs"
	"time"
)

func Now() typedefs.Epoch {
	return typedefs.Epoch(
		time.Now().Unix())
}
