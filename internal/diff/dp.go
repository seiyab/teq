package diff

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

type dpCell struct {
	loss  float64
	entry entry
	fromA int
	fromB int
}

func mixedEntries[Value any, List any](
	v1, v2 List,
	length func(List) int,
	getValue func(List, int) Value,
	getReflect func(List, int) reflect.Value,
	nx func(Value, Value) (diffTree, error),
) ([]entry, error) {
	leading := make([]entry, 0)
	for i := 0; i < length(v1) && i < length(v2); i++ {
		t, err := nx(getValue(v1, i), getValue(v2, i))
		if err != nil {
			return nil, err
		}
		if t.loss() > 0 {
			break
		}
		leading = append(leading, entry{value: same(getReflect(v1, i))})
	}
	k := len(leading)

	dp := make([][]dpCell, length(v1)-k+1)
	for i := range dp {
		dp[i] = make([]dpCell, length(v2)-k+1)
		for j := range dp[i] {
			dp[i][j] = dpCell{loss: math.MaxFloat64}
		}
	}
	dp[0][0] = dpCell{loss: 0}
	for b := 0; k+b < length(v2)+1; b++ {
		for a := 0; k+a < length(v1)+1; a++ {
			l := dp[a][b].loss
			if k+a < length(v1) {
				if l+1 < dp[a+1][b].loss {
					dp[a+1][b] = dpCell{
						loss: l + 1,
						entry: entry{
							leftOnly: true,
							value:    imbalanced(getReflect(v1, k+a)),
						},
						fromA: a,
						fromB: b,
					}
				}
			}
			if k+b < length(v2) {
				if l+1 < dp[a][b+1].loss {
					dp[a][b+1] = dpCell{
						loss: l + 1,
						entry: entry{
							rightOnly: true,
							value:     imbalanced(getReflect(v2, k+b)),
						},
						fromA: a,
						fromB: b,
					}
				}
			}
			if k+a < length(v1) && k+b < length(v2) {
				t, err := nx(getValue(v1, k+a), getValue(v2, k+b))
				if err != nil {
					return nil, err
				}
				tl := t.loss()
				m, ok := t.(mixed)
				if ok && l+tl < dp[a+1][b+1].loss {
					dp[a+1][b+1] = dpCell{
						loss:  l + tl,
						entry: entry{value: m},
						fromA: a,
						fromB: b,
					}
				}
			}
		}
	}
	a := len(dp) - 1
	b := len(dp[a]) - 1
	if dp[a][b].loss > 1_000_000 {
		return nil, fmt.Errorf("failed to compute diff")
	}

	trailing := make([]entry, 0, length(v1)+length(v2))
	for a, b := length(v1)-k, length(v2)-k; a > 0 || b > 0; {
		cell := dp[a][b]
		trailing = append(trailing, cell.entry)
		if !(cell.fromA < a || cell.fromB < b) {
			return nil, fmt.Errorf("infinite loop")
		}
		a = cell.fromA
		b = cell.fromB
	}
	reverse(trailing)

	entries := append(leading, trailing...)

	return entries, nil
}

func sliceMixedEntries(v1, v2 reflect.Value, nx next) ([]entry, error) {
	if v1.Kind() != reflect.Slice && v1.Kind() != reflect.Array ||
		v2.Kind() != reflect.Slice && v2.Kind() != reflect.Array {
		return nil, errors.New("unexpected kind")
	}
	es, err := mixedEntries(
		v1, v2,
		func(v reflect.Value) int { return v.Len() },
		func(v reflect.Value, i int) reflect.Value { return v.Index(i) },
		func(v reflect.Value, i int) reflect.Value { return v.Index(i) },
		nx,
	)
	if err != nil {
		return nil, err
	}
	return es, nil
}

func multiLineStringEntries(v1, v2 []string) ([]entry, error) {
	return mixedEntries(
		v1, v2,
		func(v []string) int { return len(v) },
		func(v []string, i int) string { return v[i] },
		func(v []string, i int) reflect.Value { return reflect.ValueOf(v[i]) },
		func(v1, v2 string) (diffTree, error) {
			return stringDiff(reflect.ValueOf(v1), reflect.ValueOf(v2), nil)
		},
	)
}

func reverse(entries []entry) {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
}
