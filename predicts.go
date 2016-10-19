package gormgen

type Predict interface {
	String() string
}

type PredictImpl struct {
	Operator string
}

func (p PredictImpl) String() string {
	return p.Operator
}

var (
	EqualPredict              = PredictImpl{Operator: "="}
	NotEqualPredict           = PredictImpl{Operator: "<>"}
	GreaterThanPredict        = PredictImpl{Operator: ">"}
	GreaterThanOrEqualPredict = PredictImpl{Operator: ">="}
	SmallerThanPredict        = PredictImpl{Operator: "<"}
	SmallerThanOrEqualPredict = PredictImpl{Operator: "<="}
	LikePredict               = PredictImpl{Operator: "LIKE"}
)
