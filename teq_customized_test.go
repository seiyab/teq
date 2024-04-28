package teq_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/seiyab/teq"
)

func TestEqual_Customized(t *testing.T) {
	t.Run("time.Time", func(t *testing.T) {
		defaultTeq := teq.New()
		customizedTeq := teq.New()
		customizedTeq.AddTransform(utc)

		mt := &mockT{}

		secondsEastOfUTC := int((8 * time.Hour).Seconds())
		beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
		d1 := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
		d2 := time.Date(2000, 2, 1, 20, 30, 0, 0, beijing)

		if defaultTeq.Equal(mt, d1, d2) {
			t.Error("expected d1 != d2, got d1 == d2 with defaultTeq")
		}
		if !customizedTeq.Equal(mt, d1, d2) {
			t.Error("expected d1 == d2, got d1 != d2 with customizedTeq")
		}
		if reflect.DeepEqual(d1, d2) {
			t.Error("expected d1 != d2, got d1 == d2 with reflect.DeepEqual")
		}
	})
}

func utc(d time.Time) time.Time {
	return d.UTC()
}
