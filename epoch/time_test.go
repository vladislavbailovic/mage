package epoch

import (
	"testing"
	"time"
)

func Test_TimeNow(t *testing.T) {
	now := time.Now().Unix()
	epoch := Now()
	if now != int64(epoch) {
		t.Fatalf("time now mismatch: %d is now, but got %v", now, epoch)
	}
}
