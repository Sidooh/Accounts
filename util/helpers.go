package util

import (
	"encoding/json"
	"sort"
)

func ConvertStruct(from interface{}, to interface{}) {
	record, _ := json.Marshal(from)
	_ = json.Unmarshal(record, &to)
}

func ReverseSlice[T comparable](s []T) {
	sort.SliceStable(s, func(i, j int) bool {
		return i > j
	})
}

func ReverseInterfaceSlice(s []interface{}) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
