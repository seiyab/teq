package teq

import "testing"

type TestingT interface {
	Helper()
	Errorf(format string, args ...interface{})
}

var _ TestingT = &testing.T{}
