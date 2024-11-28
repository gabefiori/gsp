package finder

import "sort"

type SortType int8

const (
	NoSort SortType = iota
	AscSort
	DescSort
)

func SortTypeFromStr(s string) SortType {
	switch s {
	case "asc":
		return AscSort
	case "desc":
		return DescSort
	default:
		return NoSort
	}
}

func sortResults(r []string, t SortType) {
	switch t {
	case AscSort:
		sort.Strings(r)
	case DescSort:
		sort.Sort(sort.Reverse(sort.StringSlice(r)))
	}
}
