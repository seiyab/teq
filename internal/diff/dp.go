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

func mixedEntries[List any](
	p diffProcess,
	v1, v2 List,
	length func(List) int,
	getReflect func(List, int) reflect.Value,
) ([]entry, error) {
	leading := make([]entry, 0)
	for i := 0; i < length(v1) && i < length(v2); i++ {
		t := p.diff(getReflect(v1, i), getReflect(v2, i))
		if t.loss() > 0 {
			break
		}
		leading = append(leading, entry{value: p.leftPure(getReflect(v1, i))})
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
							value:    p.leftPure(getReflect(v1, k+a)),
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
							value:     p.rightPure(getReflect(v2, k+b)),
						},
						fromA: a,
						fromB: b,
					}
				}
			}
			if k+a < length(v1) && k+b < length(v2) {
				t := p.diff(getReflect(v1, k+a), getReflect(v2, k+b))
				tl := t.loss()
				switch t.(type) {
				case mixed, cycle, nilNode, format1:
					if l+tl < dp[a+1][b+1].loss {
						dp[a+1][b+1] = dpCell{
							loss:  l + tl,
							entry: entry{value: t},
							fromA: a,
							fromB: b,
						}
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

func sliceMixedEntries(v1, v2 reflect.Value, p diffProcess) ([]entry, error) {
	if v1.Kind() != reflect.Slice && v1.Kind() != reflect.Array ||
		v2.Kind() != reflect.Slice && v2.Kind() != reflect.Array {
		return nil, errors.New("unexpected kind")
	}
	es, err := mixedEntries(
		p,
		v1, v2,
		func(v reflect.Value) int { return v.Len() },
		func(v reflect.Value, i int) reflect.Value { return v.Index(i) },
	)
	if err != nil {
		return nil, err
	}
	return es, nil
}

func multiLineStringEntries(v1, v2 []string, p diffProcess) ([]entry, error) {
	return mixedEntries(
		p,
		v1, v2,
		func(v []string) int { return len(v) },
		func(v []string, i int) reflect.Value { return reflect.ValueOf(v[i]) },
	)
}

func reverse(entries []entry) {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
}
