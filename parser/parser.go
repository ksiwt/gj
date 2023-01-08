package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pohedev/gj.git/ast"
	"github.com/pohedev/gj.git/lexer"
	"github.com/pohedev/gj.git/token"
)

// Parser represents iterating Lexer and building AST,
// holds three tokens.
type Parser struct {
	lex      *lexer.Lexer // Lexer.
	previous lexer.Item   // Previous Item.
	current  lexer.Item   // Current Item.
	peek     lexer.Item   // Peek Item.
}

// New takes a Lexer and initialize Parser,
// set current and peek Item,.
func New(lex *lexer.Lexer) *Parser {
	p := Parser{
		lex: lex,
	}

	p.next()
	p.next()

	return &p
}

// Parse parses Items and creates an AST.
func (p *Parser) Parse() (*ast.RootNode, error) {
	var node ast.RootNode
	switch p.current.Token {
	case token.LeftBrace:
		node.RootNodeType = ast.RootNodeTypeObject
	case token.LeftBracket:
		node.RootNodeType = ast.RootNodeTypeArray
	}

	if err := p.validateStartingSyntax(node); err != nil {
		return nil, err
	}

	val, parseErr := p.parseValue()
	if parseErr != nil {
		return nil, parseErr
	}
	node.Value = val

	if err := p.validateClosingSyntax(node); err != nil {
		return nil, err
	}

	return &node, nil
}

// validateStartingSyntax validate JSON starting syntax.
func (p *Parser) validateStartingSyntax(n ast.RootNode) error {
	switch n.RootNodeType {
	case ast.RootNodeTypeObject:
		if p.isCurrentToken(token.LeftBrace) {
			return nil
		}
	case ast.RootNodeTypeArray:
		if p.isCurrentToken(token.LeftBracket) {
			return nil
		}
	}
	return errors.New("failed to parse: missing JSON starting brace or bracket")
}

// validateClosingSyntax validate JSON closing syntax.
func (p *Parser) validateClosingSyntax(n ast.RootNode) error {
	switch n.RootNodeType {
	case ast.RootNodeTypeObject:
		if p.isCurrentToken(token.EOF) && p.isPreviousToken(token.RightBrace) {
			return nil
		}
	case ast.RootNodeTypeArray:
		if p.isCurrentToken(token.EOF) && p.isPreviousToken(token.RightBracket) {
			return nil
		}
	}
	return errors.New("failed to parse: missing JSON closing brace or bracket")
}

// next sets and advance Item which include token.
// in the process,
// - set current to previous.
// - set peek to current.
// - set returned from p.lex.NextItem() to peek.
func (p *Parser) next() {
	p.previous = p.current
	p.current = p.peek
	p.peek = p.lex.NextItem()
}

// parseValue is the entry point for parsing JSON values.
func (p *Parser) parseValue() (*ast.Value, error) {
	value := ast.Value{}

	switch p.current.Token {
	case token.LeftBrace:
		objValue, parseErr := p.parseObject()
		if parseErr != nil {
			return nil, parseErr
		}
		value.Value = objValue

	case token.LeftBracket:
		arrayVal, parseErr := p.parseArray()
		if parseErr != nil {
			return nil, parseErr
		}
		value.Value = arrayVal

	default:
		litValue, parseErr := p.parseLiteral()
		if parseErr != nil {
			return nil, parseErr
		}
		value.Value = litValue
	}

	return &value, nil
}

// parseObject parses JSON object.
func (p *Parser) parseObject() (*ast.Object, error) {
	obj := ast.Object{}
	objState := ast.StateObjectStart

	for {
		if p.isCurrentToken(token.EOF) {
			break
		}

		switch objState {
		case ast.StateObjectStart:
			if p.isCurrentToken(token.LeftBrace) {
				obj.Start = p.current.Pos
				objState = ast.StateObjectOpen
				p.next()
			} else {
				return nil, fmt.Errorf(
					"failed to parse object: expected LeftBrace token but got: %v",
					p.current.Val,
				)
			}

		case ast.StateObjectOpen:
			if p.isCurrentToken(token.RightBrace) {
				obj.End = p.current.Pos
				return &obj, nil
			}
			prop, parseErr := p.parseProperty()
			if parseErr != nil {
				return nil, parseErr
			}
			obj.Children = append(obj.Children, *prop)
			objState = ast.StateObjectProperty

		case ast.StateObjectProperty:
			if p.isCurrentToken(token.RightBrace) {
				p.next()
				obj.End = p.current.Pos
				return &obj, nil
			} else if p.isCurrentToken(token.Comma) {
				objState = ast.StateObjectComma
				p.next()
			} else {
				return nil, fmt.Errorf(
					"failed to parse property: expected RightBrace or Comma token but got: %v",
					p.current.Val,
				)
			}

		case ast.StateObjectComma:
			if p.isCurrentToken(token.RightBrace) {
				obj.End = p.current.Pos
				return &obj, nil
			}
			prop, parseErr := p.parseProperty()
			if parseErr != nil {
				return nil, parseErr
			}
			obj.Children = append(obj.Children, *prop)
			objState = ast.StateObjectProperty
		}
	}

	obj.End = p.current.Pos
	return &obj, nil
}

// parseProperty parses JSON key value pair property.
func (p *Parser) parseProperty() (*ast.Property, error) {
	prop := ast.Property{}
	propertyState := ast.StatePropertyStart

	for {
		if p.isCurrentToken(token.EOF) {
			break
		}

		switch propertyState {
		case ast.StatePropertyStart:
			if p.isCurrentToken(token.String) {
				prop.Identifier = ast.Identifier{Value: p.parseString()}
				propertyState = ast.StatePropertyKey
				p.next()
			} else {
				return nil, fmt.Errorf(
					"failed to parse property start: expected String token but got: %v",
					p.current.Val,
				)
			}

		case ast.StatePropertyKey:
			if p.isCurrentToken(token.Colon) {
				propertyState = ast.StatePropertyColon
				p.next()
			} else {
				return nil, fmt.Errorf(
					"failed to parse property key: expected Colon token but got: %v",
					p.current.Val,
				)
			}

		case ast.StatePropertyColon:
			value, parseErr := p.parseValue()
			if parseErr != nil {
				return nil, parseErr
			}
			prop.Value = value
			propertyState = ast.StatePropertyValue

		case ast.StatePropertyValue:
			return &prop, nil
		}
	}

	return &prop, nil
}

// parseArray parses JSON array.
func (p *Parser) parseArray() (*ast.Array, error) {
	array := ast.Array{}
	arrayState := ast.StateArrayStart

	for {
		if p.isCurrentToken(token.EOF) {
			break
		}

		switch arrayState {
		case ast.StateArrayStart:
			if p.isCurrentToken(token.LeftBracket) {
				array.Start = p.current.Pos
				arrayState = ast.StateArrayOpen
				p.next()
			}

		case ast.StateArrayOpen:
			if p.isCurrentToken(token.RightBracket) {
				array.End = p.current.Pos
				return &array, nil
			}
			arrayItem, parseErr := p.parseArrayItem()
			if parseErr != nil {
				return nil, parseErr
			}
			array.Children = append(array.Children, *arrayItem)
			arrayState = ast.StateArrayValue
			if p.isPeekToken(token.RightBracket) {
				p.next()
			}

		case ast.StateArrayValue:
			if p.isCurrentToken(token.RightBracket) {
				array.End = p.current.Pos
				p.next()
				return &array, nil
			} else if p.isCurrentToken(token.Comma) {
				arrayState = ast.StateArrayComma
				p.next()
			} else {
				return nil, fmt.Errorf(
					"failed to parse array: expected RightBrace or Comma token but got: %v",
					p.current.Val,
				)
			}

		case ast.StateArrayComma:
			if p.isCurrentToken(token.RightBracket) {
				array.End = p.current.Pos
				return &array, nil
			}
			arrayItem, parseErr := p.parseArrayItem()
			if parseErr != nil {
				return nil, parseErr
			}
			array.Children = append(array.Children, *arrayItem)
			arrayState = ast.StateArrayValue
		}
	}

	array.End = p.current.Pos
	return &array, nil
}

// parseArrayItem parses item inside JSON array.
func (p *Parser) parseArrayItem() (*ast.ArrayItem, error) {
	item := ast.ArrayItem{}

	switch p.current.Token {
	case token.LeftBrace:
		objValue, parseErr := p.parseObject()
		if parseErr != nil {
			return nil, parseErr
		}
		item.Value = objValue

	case token.LeftBracket:
		arrayValue, parseErr := p.parseArray()
		if parseErr != nil {
			return nil, parseErr
		}
		item.Value = arrayValue

	default:
		litValue, parseErr := p.parseLiteral()
		if parseErr != nil {
			return nil, parseErr
		}
		item.Value = litValue
	}

	return &item, nil
}

// parseLiteral parse JSON literal.
func (p *Parser) parseLiteral() (*ast.Literal, error) {
	lit := ast.Literal{}

	defer p.next()

	switch p.current.Token {
	case token.String:
		lit.LiteralType = ast.LiteralTypeString
		lit.Val = p.parseString()

	case token.Number:
		lit.LiteralType = ast.LiteralTypeNumber
		ct := p.current.Val
		i, parseIntErr := strconv.ParseInt(ct, 10, 64)
		if parseIntErr == nil {
			lit.Val = i
		} else {
			f, parseFloatErr := strconv.ParseFloat(ct, 64)
			if parseFloatErr != nil {
				return nil, fmt.Errorf(
					"failed to parse number: incorrect syntax %v",
					p.current.Val,
				)
			}
			lit.Val = f
		}

	case token.True:
		lit.LiteralType = ast.LiteralTypeTrue
		lit.Val = true

	case token.False:
		lit.LiteralType = ast.LiteralTypeFalse
		lit.Val = false

	case token.Null:
		lit.LiteralType = ast.LiteralTypeNull
		lit.Val = "null"

	default:
		return nil, fmt.Errorf(
			"failed to parse literal: incorrect syntax %v",
			p.current.Val,
		)
	}

	return &lit, nil
}

// parseString parses JSON string literal.
func (p *Parser) parseString() string {
	s, _ := strconv.Unquote(p.current.Val)
	return s
}

// isPreviousToken reports whether t is previous Token.
func (p *Parser) isPreviousToken(t token.Token) bool {
	return p.previous.Token == t
}

// isCurrentToken reports whether t is current Token.
func (p *Parser) isCurrentToken(t token.Token) bool {
	return p.current.Token == t
}

// isPeekToken reports whether t is peek Token.
func (p *Parser) isPeekToken(t token.Token) bool {
	return p.peek.Token == t
}
