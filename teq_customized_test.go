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

		secondsEastOfUTC := int((8 * time.Hour).Seconds())
		beijing := time.FixedZone("Beijing Time", secondsEastOfUTC)
		d1 := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
		d2 := time.Date(2000, 2, 1, 20, 30, 0, 0, beijing)

		defaultTeq.NotEqual(t, d1, d2)
		customizedTeq.Equal(t, d1, d2)
		if reflect.DeepEqual(d1, d2) {
			t.Error("expected d1 != d2, got d1 == d2 with reflect.DeepEqual")
		}

		type twoDates struct {
			d1 time.Time
			d2 time.Time
		}
		dt1 := twoDates{d1, d2}
		dt2 := twoDates{d2, d1}

		defaultTeq.NotEqual(t, dt1, dt2)
		customizedTeq.Equal(t, dt1, dt2)

		if reflect.DeepEqual(dt1, dt2) {
			t.Error("expected dt1 != dt2, got dt1 == dt2 with reflect.DeepEqual")
		}

		ds1 := []time.Time{d1, d1, d2}
		ds2 := []time.Time{d2, d1, d1}
		defaultTeq.NotEqual(t, ds1, ds2)
		customizedTeq.Equal(t, ds1, ds2)
		if reflect.DeepEqual(ds1, ds2) {
			t.Error("expected ds1 != ds2, got ds1 == ds2 with reflect.DeepEqual")
		}
	})
}

func utc(d time.Time) time.Time {
	return d.UTC()
}
