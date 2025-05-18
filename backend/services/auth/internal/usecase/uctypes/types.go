package uctypes

type QueryGetListParams struct {
	WithDeleted bool
	ForShare    bool
	ForUpdate   bool
	Limit       uint64
	Offset      uint64
}

type QueryGetOneParams struct {
	WithDeleted bool
	ForShare    bool
	ForUpdate   bool
}

type CompareType string

const (
	CompareTypeEqual       CompareType = "equal"
	CompareTypeLess        CompareType = "less"
	CompareTypeMore        CompareType = "more"
	CompareTypeLessOrEqual CompareType = "less_or_equal"
	CompareTypeMoreOrEqual CompareType = "more_or_equal"
)
