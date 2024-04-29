package teq

import "testing"

type TestingT interface {
	Helper()
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Log(args ...interface{})
}

var _ TestingT = &testing.T{}
