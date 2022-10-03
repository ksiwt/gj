package lexer

import (
	"testing"

	"github.com/pohedev/gj.git/token"
)

type lexTest struct {
	name      string
	input     string
	wantItems []Item
}

func mkItem(tok token.Token, text string) Item {
	return Item{
		token: tok,
		val:   text,
	}
}

var (
	tLeftBrace    = mkItem(token.LeftBrace, "{")
	tRightBrace   = mkItem(token.RightBrace, "}")
	tLeftBracket  = mkItem(token.LeftBracket, "[")
	tRightBracket = mkItem(token.RightBracket, "]")
	tColon        = mkItem(token.Colon, ":")
	tComma        = mkItem(token.Comma, ",")
	tTrue         = mkItem(token.True, "true")
	tFalse        = mkItem(token.False, "false")
	tNull         = mkItem(token.Null, "null")
	tEOF          = mkItem(token.EOF, "")
)

func lexToSlice(input string) []Item {
	var items []Item
	lexer := Lex(input)
	for {
		item := lexer.nextItem()
		items = append(items, item)
		if item.token == token.EOF || item.token == token.Error {
			break
		}
	}
	return items
}

func equal(i1, i2 []Item, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].token != i2[k].token {
			return false
		}
		if i1[k].val != i2[k].val {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

func TestLexToken(t *testing.T) {
	var tests = []lexTest{
		{
			"string",
			`{"color" : "blue"}`,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"color"`),
				tColon,
				mkItem(token.String, `"blue"`),
				tRightBrace,
				tEOF,
			},
		},
		{
			"booleans",
			`{
					"boolean1" : true,
					"boolean2" : false
				}`,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"boolean1"`),
				tColon,
				tTrue,
				tComma,
				mkItem(token.String, `"boolean2"`),
				tColon,
				tFalse,
				tRightBrace,
				tEOF,
			},
		},
		{
			"numbers",
			`{
					"number_1" : 210,
					"number_2" : -210,
					"number_3" : 21.05,
					"number_4" : 1.0E+2,
					"number_5" : 2e+308
				}`,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"number_1"`),
				tColon,
				mkItem(token.Number, "210"),
				tComma,
				mkItem(token.String, `"number_2"`),
				tColon,
				mkItem(token.Number, "-210"),
				tComma,
				mkItem(token.String, `"number_3"`),
				tColon,
				mkItem(token.Number, "21.05"),
				tComma,
				mkItem(token.String, `"number_4"`),
				tColon,
				mkItem(token.Number, "1.0E+2"),
				tComma,
				mkItem(token.String, `"number_5"`),
				tColon,
				mkItem(token.Number, "2e+308"),
				tRightBrace,
				tEOF,
			},
		},
		{
			"alphanumeric",
			`{"value" : "abc123"}`,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"value"`),
				tColon,
				mkItem(token.String, `"abc123"`),
				tRightBrace,
				tEOF,
			},
		},
		{
			"null",
			`{"value" : null}`,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"value"`),
				tColon,
				tNull,
				tRightBrace,
				tEOF,
			},
		},
		{
			"array",
			`{"cars":["Ford", "BMW"]} `,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"cars"`),
				tColon,
				tLeftBracket,
				mkItem(token.String, `"Ford"`),
				tComma,
				mkItem(token.String, `"BMW"`),
				tRightBracket,
				tRightBrace,
				tEOF,
			},
		},
		{
			"nested array",
			`{"id": [ [ 12, 23 ], [ 34, 45 ] ]}`,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"id"`),
				tColon,
				tLeftBracket,
				tLeftBracket,
				mkItem(token.Number, "12"),
				tComma,
				mkItem(token.Number, "23"),
				tRightBracket,
				tComma,
				tLeftBracket,
				mkItem(token.Number, "34"),
				tComma,
				mkItem(token.Number, "45"),
				tRightBracket,
				tRightBracket,
				tRightBrace,
				tEOF,
			},
		},
		{
			"nested object",
			`{
					"id": "123",
					"product": {
						"model": {
						"property": {
								"battery": "li-ion"
							}
						}
					},
					"status": "In-Stock"
				}`,
			[]Item{
				tLeftBrace,
				mkItem(token.String, `"id"`),
				tColon,
				mkItem(token.String, `"123"`),
				tComma,
				mkItem(token.String, `"product"`),
				tColon,
				tLeftBrace,
				mkItem(token.String, `"model"`),
				tColon,
				tLeftBrace,
				mkItem(token.String, `"property"`),
				tColon,
				tLeftBrace,
				mkItem(token.String, `"battery"`),
				tColon,
				mkItem(token.String, `"li-ion"`),
				tRightBrace,
				tRightBrace,
				tRightBrace,
				tComma,
				mkItem(token.String, `"status"`),
				tColon,
				mkItem(token.String, `"In-Stock"`),
				tRightBrace,
				tEOF,
			},
		},
		{
			"unknown",
			`*`,
			[]Item{
				mkItem(token.Unknown, `*`),
				tEOF,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := lexToSlice(tt.input)
			if !equal(items, tt.wantItems, false) {
				t.Errorf("%s: got\n\t%v\nexpected\n\t%v", tt.name, items, tt.wantItems)
			}
		})
	}
}
