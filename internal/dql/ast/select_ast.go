package ast

type AggregateOrder struct {
	Func  string // COUNT, SUM, etc.
	Field string // duration, id
	Desc  bool   // true if DESC
}

type SelectQueryAST struct {
	Fields       []string
	Table        string
	Conditions   map[string]string
	Aggregates   []AggregateFunc
	GroupBy      []string
	OrderBy      []OrderByField
	OrderByAgg   []AggregateOrder // NEW: Aggregate ORDER BY
	Limit        int
	IsCount      bool
	HasCountStar bool
}

type AggregateFunc struct {
	Func  string // SUM, AVG, MAX, MIN, COUNT
	Field string // duration, retries, etc.
}

type OrderByField struct {
	Field   string
	AggFunc string // e.g., SUM, COUNT
	Desc    bool
}
