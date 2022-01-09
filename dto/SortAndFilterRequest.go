package dto

type SortBy struct {
	Field string
	Dir   string
}

type SortAndFilterRequest struct {
	Sorts []SortBy
}
