package postgres

type Condition struct {
	Field    string
	Operator Operator
	Value    interface{}
}

type (
	Conditions []Condition
	Operator   string
)

const (
	OperatorEqual            Operator = "="
	OperatorIn               Operator = "IN"
	OperatorNotIn            Operator = "NOT IN"
	OperatorIs               Operator = "IS"
	OperatorIsNot            Operator = "IS NOT"
	OperatorGreaterThan      Operator = ">"
	OperatorGreaterThanEqual Operator = ">="
	OperatorLessThanEqual    Operator = "<="
	OperatorLike             Operator = "LIKE"
)

type OrderBy int8

const (
	OrderByDateDescending OrderBy = iota
	OrderByDateAscending
	OrderByDateUnsorted
)

func (o OrderBy) String() string {
	switch o {
	case OrderByDateDescending:
		return "DESC"
	case OrderByDateAscending:
		return "ASC"
	default:
		return ""
	}
}
