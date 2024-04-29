package teq_test

import (
	"fmt"

	"github.com/seiyab/teq"
)

type mockT struct {
	errors []string
}

var _ teq.TestingT = &mockT{}

func (t *mockT) Helper() {}

func (t *mockT) Error(args ...interface{}) {
	t.errors = append(t.errors, fmt.Sprint(args...))
}

func (t *mockT) Errorf(format string, args ...interface{}) {
	t.errors = append(t.errors, fmt.Sprintf(format, args...))
}

func (t *mockT) Log(args ...interface{}) {}
