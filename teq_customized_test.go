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

func TestEqual_Format(t *testing.T) {
	assert := teq.New()
	assert.AddFormat(func(d time.Time) string {
		return d.Format(time.RFC3339)
	})
	assert.AddFormat(func(d time.Duration) string {
		return d.String()
	})

	t.Run("time.Time", func(t *testing.T) {
		t.Run("naive", func(t *testing.T) {
			mt := &mockT{}
			assert.Equal(
				mt,
				time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
				time.Date(2000, 2, 1, 20, 30, 0, 0, time.UTC),
			)
			if len(mt.errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(mt.errors))
			}
			expected := `not equal
differences:
--- expected
+++ actual
@@ -1 +1 @@
-2000-02-01T12:30:00Z
+2000-02-01T20:30:00Z
`
			if mt.errors[0] != expected {
				t.Errorf("expected %q, got %q", expected, mt.errors[0])
			}
		})

		t.Run("nested", func(t *testing.T) {
			mt := &mockT{}
			assert.Equal(
				mt,
				map[int]time.Time{
					1: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
					2: time.Date(2000, 2, 1, 20, 30, 0, 0, time.UTC),
					3: time.Date(2000, 2, 2, 10, 0, 0, 0, time.UTC),
				},
				map[int]time.Time{
					1: time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC),
					2: time.Date(2000, 2, 1, 20, 0, 0, 0, time.UTC),
					3: time.Date(2000, 2, 2, 10, 0, 0, 0, time.UTC),
					4: time.Date(2000, 2, 2, 20, 30, 0, 0, time.UTC),
				},
			)
			if len(mt.errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(mt.errors))
			}
			expected := `not equal
differences:
--- expected
+++ actual
@@ -2,4 +2,5 @@
   int(1): 2000-02-01T12:30:00Z,
-  int(2): 2000-02-01T20:30:00Z,
+  int(2): 2000-02-01T20:00:00Z,
   int(3): 2000-02-02T10:00:00Z,
+  int(4): 2000-02-02T20:30:00Z,
 }
`
			if mt.errors[0] != expected {
				t.Errorf("expected %q, got %q", expected, mt.errors[0])
			}
		})
	})

	t.Run("Duration", func(t *testing.T) {
		mt := &mockT{}
		assert.Equal(mt, []time.Duration{1 * time.Hour}, []time.Duration{2 * time.Second})
		if len(mt.errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(mt.errors))
		}
		expected := `not equal
differences:
--- expected
+++ actual
@@ -1,3 +1,3 @@
 []time.Duration{
-  1h0m0s,
+  2s,
 }
`
		if mt.errors[0] != expected {
			t.Errorf("expected %q, got %q", expected, mt.errors[0])
		}
		assert.Equal(t, []string{expected}, mt.errors)
	})
}

func utc(d time.Time) time.Time {
	return d.UTC()
}
