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
		Token: tok,
		Val:   text,
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
		item := lexer.NextItem()
		items = append(items, item)
		if item.Token == token.EOF || item.Token == token.Error {
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
		if i1[k].Token != i2[k].Token {
			return false
		}
		if i1[k].Val != i2[k].Val {
			return false
		}
		if checkPos && i1[k].Pos != i2[k].Pos {
			return false
		}
	}
	return true
}

func TestLex(t *testing.T) {
	input := `{
					"glossary": {
						"title": "example glossary",
						"GlossDiv": {
							"title": "S",
							"GlossList": {
								"GlossEntry": {
									"GlossTerm": "Standard Generalized Markup Language",
									"Abbrev": "ISO 8879:1986",
									"GlossDef": {
										"para": "A meta-markup language, used to create markup languages such as DocBook.",
										"GlossSeeAlso": ["GML", "XML"]
									},
									"GlossSee": "markup"
								}
							},
							"Nums": 5245243
						}
					}
				}`

	wantItems := []Item{
		tLeftBrace,
		mkItem(token.String, `"glossary"`),
		tColon,
		tLeftBrace,
		mkItem(token.String, `"title"`),
		tColon,
		mkItem(token.String, `"example glossary"`),
		tComma,
		mkItem(token.String, `"GlossDiv"`),
		tColon,
		tLeftBrace,
		mkItem(token.String, `"title"`),
		tColon,
		mkItem(token.String, `"S"`),
		tComma,
		mkItem(token.String, `"GlossList"`),
		tColon,
		tLeftBrace,
		mkItem(token.String, `"GlossEntry"`),
		tColon,
		tLeftBrace,
		mkItem(token.String, `"GlossTerm"`),
		tColon,
		mkItem(token.String, `"Standard Generalized Markup Language"`),
		tComma,
		mkItem(token.String, `"Abbrev"`),
		tColon,
		mkItem(token.String, `"ISO 8879:1986"`),
		tComma,
		mkItem(token.String, `"GlossDef"`),
		tColon,
		tLeftBrace,
		mkItem(token.String, `"para"`),
		tColon,
		mkItem(token.String, `"A meta-markup language, used to create markup languages such as DocBook."`),
		tComma,
		mkItem(token.String, `"GlossSeeAlso"`),
		tColon,
		tLeftBracket,
		mkItem(token.String, `"GML"`),
		tComma,
		mkItem(token.String, `"XML"`),
		tRightBracket,
		tRightBrace,
		tComma,
		mkItem(token.String, `"GlossSee"`),
		tColon,
		mkItem(token.String, `"markup"`),
		tRightBrace,
		tRightBrace,
		tComma,
		mkItem(token.String, `"Nums"`),
		tColon,
		mkItem(token.Number, "5245243"),
		tRightBrace,
		tRightBrace,
		tRightBrace,
		tEOF,
	}

	items := lexToSlice(input)
	if !equal(items, wantItems, false) {
		t.Errorf("got\n\t%v\nexpected\n\t%v", items, wantItems)
	}
}

func TestLexToken(t *testing.T) {
	var tests = []lexTest{
		{
			"string",
			`{"color": "blue"}`,
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
					"boolean1": true,
					"boolean2": false
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
					"number_1": 210,
					"number_2": -210,
					"number_3": 21.05,
					"number_4": 1.0E+2,
					"number_5": 2e+308
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
			`{"value": "abc123"}`,
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
			`{"value": null}`,
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
			"root array",
			`[{"id": 1, "name":"water"}, {"id": 2,"name":"knife"}]`,
			[]Item{
				tLeftBracket,
				tLeftBrace,
				mkItem(token.String, `"id"`),
				tColon,
				mkItem(token.Number, "1"),
				tComma,
				mkItem(token.String, `"name"`),
				tColon,
				mkItem(token.String, `"water"`),
				tRightBrace,
				tComma,
				tLeftBrace,
				mkItem(token.String, `"id"`),
				tColon,
				mkItem(token.Number, "2"),
				tComma,
				mkItem(token.String, `"name"`),
				tColon,
				mkItem(token.String, `"knife"`),
				tRightBrace,
				tRightBracket,
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
