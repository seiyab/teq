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

func sliceEntries(v1, v2 reflect.Value, nx next) ([]entry, error) {
	if v1.Kind() != reflect.Slice || v2.Kind() != reflect.Slice {
		return nil, errors.New("unexpected kind")
	}
	leading := make([]entry, 0)
	for i := 0; i < v1.Len() && i < v2.Len(); i++ {
		t, err := nx(v1.Index(i), v2.Index(i))
		if err != nil {
			return nil, err
		}
		if t.loss > 0 {
			break
		}
		leading = append(leading, entry{value: same(v1.Index(i))})
	}
	k := len(leading)

	dp := make([][]dpCell, v1.Len()-k+1)
	for i := range dp {
		dp[i] = make([]dpCell, v2.Len()-k+1)
		for j := range dp[i] {
			dp[i][j] = dpCell{loss: math.MaxFloat64}
		}
	}
	dp[0][0] = dpCell{loss: 0}
	for a := 0; k+a < v1.Len()+1; a++ {
		for b := 0; k+b < v2.Len()+1; b++ {
			l := dp[a][b].loss
			if k+a < v1.Len() {
				if l+1 < dp[a+1][b].loss {
					dp[a+1][b] = dpCell{
						loss: l + 1,
						entry: entry{
							leftOnly: true,
							value:    imbalanced(v1.Index(k + a)),
						},
						fromA: a,
						fromB: b,
					}
				}
			}
			if k+b < v2.Len() {
				if l+1 < dp[a][b+1].loss {
					dp[a][b+1] = dpCell{
						loss: l + 1,
						entry: entry{
							rightOnly: true,
							value:     imbalanced(v2.Index(k + b)),
						},
						fromA: a,
						fromB: b,
					}
				}
			}
			if k+a < v1.Len() && k+b < v2.Len() {
				t, err := nx(v1.Index(k+a), v2.Index(k+b))
				if err != nil {
					return nil, err
				}
				if l+t.loss < dp[a+1][b+1].loss {
					dp[a+1][b+1] = dpCell{
						loss:  l + t.loss,
						entry: entry{value: t},
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
		for i := range dp {
			for j := range dp[i] {
				fmt.Printf("%v ", dp[i][j].loss)
			}
			fmt.Println()
		}
		return nil, fmt.Errorf("faild to compute diff")
	}

	trailing := make([]entry, 0, v1.Len()+v2.Len())
	for a, b := v1.Len()-k, v2.Len()-k; a > 0 || b > 0; {
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

func reverse(entries []entry) {
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
}
