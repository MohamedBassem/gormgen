package gormgen

type Predicate interface {
	String() string
}

type PredicateImpl struct {
	Operator string
}

func (p PredicateImpl) String() string {
	return p.Operator
}

var (
	EqualPredicate              = PredicateImpl{Operator: "="}
	NotEqualPredicate           = PredicateImpl{Operator: "<>"}
	GreaterThanPredicate        = PredicateImpl{Operator: ">"}
	GreaterThanOrEqualPredicate = PredicateImpl{Operator: ">="}
	SmallerThanPredicate        = PredicateImpl{Operator: "<"}
	SmallerThanOrEqualPredicate = PredicateImpl{Operator: "<="}
	LikePredicate               = PredicateImpl{Operator: "LIKE"}
)
