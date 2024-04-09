package time_wheel

import (
	"context"
	"testing"
	"time"
)

func TestTimeWheel(t *testing.T) {
	tw := NewTimeWheel(context.TODO(), 60, 1*time.Second)
	pos, cyc := tw.getPositionAndCycle(time.Now().Add(90 * time.Second))
	t.Log(pos, cyc)
}
