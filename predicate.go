package gormgen

// Predicate is a string that acts as a condition in the where clause
type Predicate string

var (
	EqualPredicate              = "="
	NotEqualPredicate           = "<>"
	GreaterThanPredicate        = ">"
	GreaterThanOrEqualPredicate = ">="
	SmallerThanPredicate        = "<"
	SmallerThanOrEqualPredicate = "<="
	LikePredicate               = "LIKE"
)
