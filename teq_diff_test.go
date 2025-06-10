package teq_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/seiyab/teq"
)

func TestEqual_StringFormat(t *testing.T) {
	assert := teq.New()
	a := `abc
def
ghi
jkl
mno`
	b := `abc
ghi
jkl
mno`
	mt := &mockT{}
	assert.Equal(mt, a, b)
	if len(mt.errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(mt.errors))
	}
	expected := `not equal
differences:
--- expected
+++ actual
  string(
    abc
-   def
    ghi
    jkl
:
  )`
	if mt.errors[0] != expected {
		t.Errorf("expected %q, got %q", expected, mt.errors[0])
	}
	assert.Equal(t, expected, mt.errors[0])
}

func TestEqual_Format(t *testing.T) {
	assert := teq.New()
	t.Run("array", func(t *testing.T) {
		mt := &mockT{}
		assert.Equal(mt, [1]int{0}, [1]int{1})
		if len(mt.errors) != 1 {
			t.Fatalf("expected 1 error, got %d", len(mt.errors))
		}
		expected := `not equal
differences:
--- expected
+++ actual
  [1]int{
-   0,
+   1,
  }`
		if mt.errors[0] != expected {
			t.Errorf("expected %q, got %q", expected, mt.errors[0])
		}
		assert.Equal(t, expected, mt.errors[0])
	})

	t.Run("pointer", func(t *testing.T) {
		t.Run("nil", func(t *testing.T) {
			mt := &mockT{}
			x := 100
			a := &x
			b := (*int)(nil)
			assert.Equal(mt, a, b)

			if len(mt.errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(mt.errors))
			}
			expected := `not equal
differences:
--- expected
+++ actual
- &100
+ *int(nil)`
			if mt.errors[0] != expected {
				t.Errorf("expected %q, got %q", expected, mt.errors[0])
			}
			assert.Equal(t, expected, mt.errors[0])
		})

		t.Run("pointer of struct", func(t *testing.T) {
			mt := &mockT{}
			assert.Equal(mt, &http.Client{Timeout: time.Second}, http.DefaultClient)

			if len(mt.errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(mt.errors))
			}
			expected := `not equal
differences:
--- expected
+++ actual
  &http.Client{
:
    CheckRedirect: func(*http.Request, []*http.Request) error(nil),
    Jar: http.CookieJar(nil),
-   Timeout: time.Duration("1s"),
+   Timeout: time.Duration("0s"),
  }`
			if mt.errors[0] != expected {
				t.Errorf("expected %q, got %q", expected, mt.errors[0])
			}
			assert.Equal(t, expected, mt.errors[0])
		})

		t.Run("slice of pointer", func(t *testing.T) {
			ref := func(s string) *string { return &s }
			mt := &mockT{}
			assert.Equal(mt,
				[]*string{ref("a"), ref("b"), ref("c"), nil},
				[]*string{ref("a"), ref("b"), ref("d"), nil},
			)

			if len(mt.errors) != 1 {
				t.Fatalf("expected 1 error, got %d", len(mt.errors))
			}
			expected := `not equal
differences:
--- expected
+++ actual
  []*string{
    &"a",
    &"b",
-   &"c",
+   &"d",
    *string(nil),
  }`
			if mt.errors[0] != expected {
				t.Errorf("expected %q, got %q", expected, mt.errors[0])
			}
			assert.Equal(t, expected, mt.errors[0])
		})
	})
}
