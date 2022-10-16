package ast

type RootNodeType int

const (
	RootNodeTypeObject RootNodeType = iota + 1
	RootNodeTypeArray
)

type LiteralType int

const (
	LiteralTypeString LiteralType = iota + 1
	LiteralTypeNumber
	LiteralTypeNull
	LiteralTypeTrue
	LiteralTypeFalse
)

type State int

const (
	StateObjectStart State = iota + 1
	StateObjectOpen
	StateObjectProperty
	StateObjectComma
	StatePropertyStart
	StatePropertyKey
	StatePropertyColon
	StatePropertyValue
	StateArrayStart
	StateArrayOpen
	StateArrayValue
	StateArrayComma
)

type RootNode struct {
	RootNodeType
	*Value
}

type Object struct {
	Children []Property
	Start    int
	End      int
}

type Property struct {
	Identifier
	Value any
}

type Identifier struct {
	Value string
}

type Array struct {
	Children []ArrayItem
	Start    int
	End      int
}

type ArrayItem struct {
	Value any
}

type Value struct {
	Value any
}

type Literal struct {
	LiteralType
	Val any
}

type ValueSet interface{}
