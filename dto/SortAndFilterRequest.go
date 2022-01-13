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
	Sorts   SortBy
	Filters []FilterBy
	Limit   int
	Offset  int
}

type FilterSql struct {
	SqlOperator  string
	ValueReplace string
}

var Operators = []string{"eq", "neq", "ct", "sw", "ew", "gt", "lt", "gte", "lte"}

var SqlOperatorReplacement = map[string]FilterSql{
	"eq": {
		SqlOperator:  "=",
		ValueReplace: "@@"},
	"neq": {
		SqlOperator:  "!=",
		ValueReplace: "@@",
	},
	"ct": {
		SqlOperator:  "LIKE",
		ValueReplace: "%@@%",
	},
	"sw": {
		SqlOperator:  "LIKE",
		ValueReplace: "@@%",
	},
	"ew": {
		SqlOperator:  "LIKE",
		ValueReplace: "%@@",
	},
	"gt": {
		SqlOperator:  ">",
		ValueReplace: "@@",
	},
	"lt": {
		SqlOperator:  "<",
		ValueReplace: "@@",
	},
	"gte": {
		SqlOperator:  ">=",
		ValueReplace: "@@",
	},
	"lte": {
		SqlOperator:  "<=",
		ValueReplace: "@@",
	},
}
