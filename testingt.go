package teq

import "testing"

type TestingT interface {
	Errorf(format string, args ...interface{})
}

var _ TestingT = &testing.T{}
