package util

import "sort"

func GetSortedKeys(m interface{}) []int {
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
