package util

import (
	"sort"

	"github.com/shopspring/decimal"
)

func GetSortedKeysInt(m interface{}) []int {
	keys := make([]int, 0)
	switch mm := m.(type) {
	case map[int]DateRange:
		for k := range mm {
			keys = append(keys, k)
		}
	}
	sort.Ints(keys)
	return keys
}

func GetSortedKeysString(m interface{}) []string {
	keys := make([]string, 0)
	switch mm := m.(type) {
	case map[string]decimal.Decimal:
		for k := range mm {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}
	return keys
}
