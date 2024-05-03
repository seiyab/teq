package teq_test

import (
	"testing"

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
