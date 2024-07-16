package diffshow

import (
	"reflect"

	"github.com/pmezard/go-difflib/difflib"
)

type Diff struct {
	path []segment
	a    reflect.Value
	b    reflect.Value
}

func (d Diff) Format(f func(v reflect.Value) string) (string, error) {
	a := difflib.SplitLines(f(d.a))
	b := difflib.SplitLines(f(d.b))
	diff := difflib.UnifiedDiff{
		A:        a,
		B:        b,
		FromFile: "expected",
		ToFile:   "actual",
		Context:  1,
	}
	diffTxt, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return "", err
	}
	return diffTxt, nil
}
