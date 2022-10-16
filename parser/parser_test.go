package parser

import (
	"testing"

	"github.com/pohedev/gj.git/ast"
	"github.com/pohedev/gj.git/lexer"
	"github.com/stretchr/testify/assert"
)

type parserTest struct {
	name  string
	input string
	want  *ast.RootNode
}

type parserErrorTest struct {
	name  string
	input string
}

func TestParser_Parse(t *testing.T) {
	var tests = []parserTest{
		{
			"string",
			`{"color": "blue"}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "color"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "blue"}},
							},
						},
						Start: 0,
						End:   17,
					},
				},
			},
		},
		{
			"booleans",
			`{
					"boolean_1": true,
					"boolean_2": false
				}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "boolean_1"},
								Value: &ast.Value{
									Value: &ast.Literal{LiteralType: ast.LiteralTypeTrue, Val: true},
								},
							},
							{
								Identifier: ast.Identifier{Value: "boolean_2"},
								Value: &ast.Value{
									Value: &ast.Literal{LiteralType: ast.LiteralTypeFalse, Val: false},
								},
							},
						},
						Start: 0,
						End:   55,
					},
				},
			},
		},
		{
			"numbers",
			`{
					"number_1": 210,
					"number_2": -210,
					"number_3": 21.05,
					"number_4": 1.0E+2,
				}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "number_1"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(210)}},
							},
							{
								Identifier: ast.Identifier{Value: "number_2"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(-210)}},
							},
							{
								Identifier: ast.Identifier{Value: "number_3"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: float64(21.05)}},
							},
							{
								Identifier: ast.Identifier{Value: "number_4"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: float64(100)}},
							},
						},
						Start: 0,
						End:   100,
					},
				},
			},
		},
		{
			"alphanumeric",
			`{"value": "abc123"}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "value"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "abc123"}},
							},
						},
						Start: 0,
						End:   19,
					},
				},
			},
		},
		{
			"null",
			`{"value": null}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "value"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNull, Val: "null"}},
							},
						},
						Start: 0,
						End:   15,
					},
				},
			},
		},
		{
			"array",
			`{"cars":["Ford", "BMW"]}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "cars"},
								Value: &ast.Value{
									Value: &ast.Array{
										Children: []ast.ArrayItem{
											{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "Ford"}},
											{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "BMW"}},
										},
										Start: 8,
										End:   22,
									},
								},
							},
						},
						Start: 0,
						End:   24,
					},
				},
			},
		},
		{
			"nested array",
			`{"id": [ [ 12, 23 ], [ 34, 45 ] ]}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "id"},
								Value: &ast.Value{
									Value: &ast.Array{
										Children: []ast.ArrayItem{
											{
												Value: &ast.Array{
													Children: []ast.ArrayItem{
														{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(12)}},
														{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(23)}},
													},
													Start: 9,
													End:   18,
												},
											},
											{
												Value: &ast.Array{
													Children: []ast.ArrayItem{
														{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(34)}},
														{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(45)}},
													},
													Start: 21,
													End:   30,
												},
											},
										},
										Start: 7,
										End:   32,
									},
								},
							},
						},
						Start: 0,
						End:   34,
					},
				},
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
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "id"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "123"}},
							},
							{
								Identifier: ast.Identifier{Value: "product"},
								Value: &ast.Value{
									Value: &ast.Object{
										Children: []ast.Property{
											{
												Identifier: ast.Identifier{Value: "model"},
												Value: &ast.Value{
													Value: &ast.Object{
														Children: []ast.Property{
															{
																Identifier: ast.Identifier{Value: "property"},
																Value: &ast.Value{
																	Value: &ast.Object{
																		Children: []ast.Property{
																			{
																				Identifier: ast.Identifier{Value: "battery"},
																				Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "li-ion"}},
																			},
																		},
																		Start: 73,
																		End:   117,
																	},
																},
															},
														},
														Start: 53,
														End:   124,
													},
												},
											},
										},
										Start: 36,
										End:   125,
									},
								},
							},
							{
								Identifier: ast.Identifier{Value: "status"},
								Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "In-Stock"}},
							},
						},
						Start: 0,
						End:   158,
					},
				},
			},
		},
		{
			"root array",
			`[{"id": 1, "name":"water"}, {"id": 2,"name":"knife"}]`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeArray,
				Value: &ast.Value{
					Value: &ast.Array{
						Children: []ast.ArrayItem{
							{
								Value: &ast.Object{
									Children: []ast.Property{
										{
											Identifier: ast.Identifier{Value: "id"},
											Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(1)}},
										},
										{
											Identifier: ast.Identifier{Value: "name"},
											Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "water"}},
										},
									},
									Start: 1,
									End:   26,
								},
							},
							{
								Value: &ast.Object{
									Children: []ast.Property{
										{
											Identifier: ast.Identifier{Value: "id"},
											Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(2)}},
										},
										{
											Identifier: ast.Identifier{Value: "name"},
											Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "knife"}},
										},
									},
									Start: 28,
									End:   52,
								},
							},
						},
						Start: 0,
						End:   52,
					},
				},
			},
		},
		{
			"full set",
			`{
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
				}`,
			&ast.RootNode{
				RootNodeType: ast.RootNodeTypeObject,
				Value: &ast.Value{
					Value: &ast.Object{
						Children: []ast.Property{
							{
								Identifier: ast.Identifier{Value: "glossary"},
								Value: &ast.Value{
									Value: &ast.Object{
										Children: []ast.Property{
											{
												Identifier: ast.Identifier{Value: "title"},
												Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "example glossary"}},
											},
											{
												Identifier: ast.Identifier{Value: "GlossDiv"},
												Value: &ast.Value{Value: &ast.Object{
													Children: []ast.Property{
														{
															Identifier: ast.Identifier{Value: "title"},
															Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "S"}},
														},
														{
															Identifier: ast.Identifier{Value: "GlossList"},
															Value: &ast.Value{
																Value: &ast.Object{
																	Children: []ast.Property{
																		{
																			Identifier: ast.Identifier{Value: "GlossEntry"},
																			Value: &ast.Value{
																				Value: &ast.Object{
																					Children: []ast.Property{
																						{
																							Identifier: ast.Identifier{Value: "GlossTerm"},
																							Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "Standard Generalized Markup Language"}},
																						},
																						{
																							Identifier: ast.Identifier{Value: "Abbrev"},
																							Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "ISO 8879:1986"}},
																						},
																						{
																							Identifier: ast.Identifier{Value: "GlossDef"},
																							Value: &ast.Value{Value: &ast.Object{
																								Children: []ast.Property{
																									{
																										Identifier: ast.Identifier{Value: "para"},
																										Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "A meta-markup language, used to create markup languages such as DocBook."}},
																									},
																									{
																										Identifier: ast.Identifier{Value: "GlossSeeAlso"},
																										Value: &ast.Value{Value: &ast.Array{
																											Children: []ast.ArrayItem{
																												{
																													Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "GML"},
																												},
																												{
																													Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "XML"},
																												},
																											},
																											Start: 384,
																											End:   397,
																										}},
																									},
																								},
																								Start: 262,
																								End:   409,
																							}},
																						},
																						{
																							Identifier: ast.Identifier{Value: "GlossSee"},
																							Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeString, Val: "markup"}},
																						},
																					},
																					Start: 141,
																					End:   458,
																				},
																			},
																		},
																	},
																	Start: 117,
																	End:   459,
																},
															},
														},
														{
															Identifier: ast.Identifier{Value: "Nums"},
															Value:      &ast.Value{Value: &ast.Literal{LiteralType: ast.LiteralTypeNumber, Val: int64(5245243)}},
														},
													},
													Start: 74,
													End:   497,
												}},
											},
										},
										Start: 19,
										End:   503,
									},
								},
							},
						},
						Start: 0,
						End:   504,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.Lex(tt.input))
			result, err := p.Parse()
			assert.Nil(t, err)
			assert.Equal(t, result, tt.want)
		})
	}
}
