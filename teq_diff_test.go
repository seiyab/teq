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
@@ -1,3 +1,2 @@
 abc
-def
 ghi
`
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
@@ -1,3 +1,3 @@
 [1]int{
-  int(0),
+  int(1),
 }
`
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
@@ -1 +1 @@
-*int(100)
+<nil>
`
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
@@ -4,3 +4,3 @@
   Jar: http.CookieJar(<nil>),
-  Timeout: time.Duration(1000000000),
+  Timeout: time.Duration(0),
 }
`
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
@@ -3,3 +3,3 @@
   *"b",
-  *"c",
+  *"d",
   <nil>,
`
			if mt.errors[0] != expected {
				t.Errorf("expected %q, got %q", expected, mt.errors[0])
			}
			assert.Equal(t, expected, mt.errors[0])
		})
	})
}
