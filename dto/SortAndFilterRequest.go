package dto

type SortBy struct {
	Field string
	Dir   string
}

type FilterBy struct {
	Field    string
	Operator string
	Value    interface{}
}

type SortAndFilterRequest struct {
	Sorts   []SortBy
	Filters []FilterBy
	Limit   int
	Anchor  string
}
