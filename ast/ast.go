package ast

// RootNodeType identifies the type of JSON RootNode.
type RootNodeType int

const (
	RootNodeTypeObject RootNodeType = iota + 1 // object
	RootNodeTypeArray                          // array
)

// RootNode represents a what JSON starts every parsed AST.
type RootNode struct {
	RootNodeType
	*Value
}

// LiteralType identifies the type of JSON Literal.
type LiteralType int

const (
	LiteralTypeString LiteralType = iota + 1
	LiteralTypeNumber
	LiteralTypeNull
	LiteralTypeTrue
	LiteralTypeFalse
)

// Literal represents a a JSON literal.
type Literal struct {
	LiteralType
	Val any
}

// State identifies the type of parsing JSON state.
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

// Object represents a JSON object.
type Object struct {
	Children []Property
	Start    int
	End      int
}

// Property represents a JSON object property.
type Property struct {
	Identifier
	Value any
}

// Identifier represents a key identifier of JSON object property.
type Identifier struct {
	Value string
}

// Array represents a JSON array.
type Array struct {
	Children []ArrayItem
	Start    int
	End      int
}

// ArrayItem represents a value of JSON array.
type ArrayItem struct {
	Value any
}

// Value represents a value of JSON value
// (object | array | boolean | string | number | null).
type Value struct {
	Value any
}
