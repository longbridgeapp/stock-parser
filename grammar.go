package stock_parser

// Code generated by peg -inline -switch -output grammar.go grammar.peg DO NOT EDIT.

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	ruleItem
	ruleLine
	ruleOTHER
	ruleStock
	ruleXLStock
	ruleFTStock
	ruleStockName
	ruleCode
	ruleUSCode
	ruleHKCode
	ruleACode
	ruleHSCODE
	ruleLetter
	ruleNumber
	ruleSuffix
	ruleMarket
	ruleCNMarket
	ruleHKMarket
	ruleUSMarket
	ruleSGMarket
	ruleSP
)

var rul3s = [...]string{
	"Unknown",
	"Item",
	"Line",
	"OTHER",
	"Stock",
	"XLStock",
	"FTStock",
	"StockName",
	"Code",
	"USCode",
	"HKCode",
	"ACode",
	"HSCODE",
	"Letter",
	"Number",
	"Suffix",
	"Market",
	"CNMarket",
	"HKMarket",
	"USMarket",
	"SGMarket",
	"SP",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(w io.Writer, pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Fprintf(w, " ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Fprintf(w, "%v %v\n", rule, quote)
			} else {
				fmt.Fprintf(w, "\x1B[36m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(w io.Writer, buffer string) {
	node.print(w, false, buffer)
}

func (node *node32) PrettyPrint(w io.Writer, buffer string) {
	node.print(w, true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(os.Stdout, buffer)
}

func (t *tokens32) WriteSyntaxTree(w io.Writer, buffer string) {
	t.AST().Print(w, buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(os.Stdout, buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	tree, i := t.tree, int(index)
	if i >= len(tree) {
		t.tree = append(tree, token32{pegRule: rule, begin: begin, end: end})
		return
	}
	tree[i] = token32{pegRule: rule, begin: begin, end: end}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type StockCodeParser struct {
	pos     int
	peekPos int

	Buffer string
	buffer []rune
	rules  [22]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *StockCodeParser) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *StockCodeParser) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *StockCodeParser
	max token32
}

func (e *parseError) Error() string {
	tokens, err := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		err += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return err
}

func (p *StockCodeParser) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *StockCodeParser) WriteSyntaxTree(w io.Writer) {
	p.tokens32.WriteSyntaxTree(w, p.Buffer)
}

func (p *StockCodeParser) SprintSyntaxTree() string {
	var bldr strings.Builder
	p.WriteSyntaxTree(&bldr)
	return bldr.String()
}

func Pretty(pretty bool) func(*StockCodeParser) error {
	return func(p *StockCodeParser) error {
		p.Pretty = pretty
		return nil
	}
}

func Size(size int) func(*StockCodeParser) error {
	return func(p *StockCodeParser) error {
		p.tokens32 = tokens32{tree: make([]token32, 0, size)}
		return nil
	}
}
func (p *StockCodeParser) Init(options ...func(*StockCodeParser) error) error {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	for _, option := range options {
		err := option(p)
		if err != nil {
			return err
		}
	}
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := p.tokens32
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 Item <- <(Line* !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					{
						position4 := position
						{
							position5, tokenIndex5 := position, tokenIndex
							{
								position7 := position
								if buffer[position] != rune('$') {
									goto l6
								}
								position++
								{
									position8, tokenIndex8 := position, tokenIndex
									if !_rules[ruleStockName]() {
										goto l9
									}
									goto l8
								l9:
									position, tokenIndex = position8, tokenIndex8
									if !_rules[ruleStockName]() {
										goto l6
									}
									if buffer[position] != rune('-') {
										goto l6
									}
									position++
									if !_rules[ruleLetter]() {
										goto l6
									}
								l10:
									{
										position11, tokenIndex11 := position, tokenIndex
										if !_rules[ruleLetter]() {
											goto l11
										}
										goto l10
									l11:
										position, tokenIndex = position11, tokenIndex11
									}
								}
							l8:
								if buffer[position] != rune('(') {
									goto l6
								}
								position++
								if !_rules[ruleCode]() {
									goto l6
								}
								if buffer[position] != rune('.') {
									goto l6
								}
								position++
								if !_rules[ruleMarket]() {
									goto l6
								}
								if buffer[position] != rune(')') {
									goto l6
								}
								position++
								if buffer[position] != rune('$') {
									goto l6
								}
								position++
								add(ruleFTStock, position7)
							}
							goto l5
						l6:
							position, tokenIndex = position5, tokenIndex5
							{
								position13 := position
								{
									position14, tokenIndex14 := position, tokenIndex
									if buffer[position] != rune('$') {
										goto l15
									}
									position++
									if !_rules[ruleCode]() {
										goto l15
									}
									{
										position16, tokenIndex16 := position, tokenIndex
										if !_rules[ruleSuffix]() {
											goto l17
										}
										goto l16
									l17:
										position, tokenIndex = position16, tokenIndex16
										{
											position18, tokenIndex18 := position, tokenIndex
											if !_rules[ruleSuffix]() {
												goto l18
											}
											goto l19
										l18:
											position, tokenIndex = position18, tokenIndex18
										}
									l19:
									}
								l16:
									if buffer[position] != rune('$') {
										goto l15
									}
									position++
									goto l14
								l15:
									position, tokenIndex = position14, tokenIndex14
									if buffer[position] != rune('$') {
										goto l20
									}
									position++
									if !_rules[ruleCode]() {
										goto l20
									}
									if !_rules[ruleSuffix]() {
										goto l20
									}
									goto l14
								l20:
									position, tokenIndex = position14, tokenIndex14
									if buffer[position] != rune('(') {
										goto l21
									}
									position++
									{
										position22, tokenIndex22 := position, tokenIndex
										if buffer[position] != rune('N') {
											goto l23
										}
										position++
										if buffer[position] != rune('Y') {
											goto l23
										}
										position++
										if buffer[position] != rune('S') {
											goto l23
										}
										position++
										if buffer[position] != rune('E') {
											goto l23
										}
										position++
										goto l22
									l23:
										position, tokenIndex = position22, tokenIndex22
										if buffer[position] != rune('N') {
											goto l21
										}
										position++
										if buffer[position] != rune('A') {
											goto l21
										}
										position++
										if buffer[position] != rune('S') {
											goto l21
										}
										position++
										if buffer[position] != rune('D') {
											goto l21
										}
										position++
										if buffer[position] != rune('A') {
											goto l21
										}
										position++
										if buffer[position] != rune('Q') {
											goto l21
										}
										position++
									}
								l22:
									{
										position24, tokenIndex24 := position, tokenIndex
										if buffer[position] != rune('：') {
											goto l25
										}
										position++
										goto l24
									l25:
										position, tokenIndex = position24, tokenIndex24
										if buffer[position] != rune(':') {
											goto l21
										}
										position++
									}
								l24:
								l26:
									{
										position27, tokenIndex27 := position, tokenIndex
										{
											position28 := position
											{
												position29, tokenIndex29 := position, tokenIndex
												if buffer[position] != rune(' ') {
													goto l30
												}
												position++
												goto l29
											l30:
												position, tokenIndex = position29, tokenIndex29
												if buffer[position] != rune('\t') {
													goto l27
												}
												position++
											}
										l29:
											add(ruleSP, position28)
										}
										goto l26
									l27:
										position, tokenIndex = position27, tokenIndex27
									}
									if !_rules[ruleUSCode]() {
										goto l21
									}
									if buffer[position] != rune(')') {
										goto l21
									}
									position++
									goto l14
								l21:
									position, tokenIndex = position14, tokenIndex14
									{
										switch buffer[position] {
										case '(':
											if buffer[position] != rune('(') {
												goto l12
											}
											position++
											if !_rules[ruleMarket]() {
												goto l12
											}
											if buffer[position] != rune(':') {
												goto l12
											}
											position++
											if !_rules[ruleCode]() {
												goto l12
											}
											if buffer[position] != rune(')') {
												goto l12
											}
											position++
										case '$':
											if buffer[position] != rune('$') {
												goto l12
											}
											position++
											if !_rules[ruleUSCode]() {
												goto l12
											}
										default:
											if buffer[position] != rune(' ') {
												goto l12
											}
											position++
										l32:
											{
												position33, tokenIndex33 := position, tokenIndex
												if buffer[position] != rune(' ') {
													goto l33
												}
												position++
												goto l32
											l33:
												position, tokenIndex = position33, tokenIndex33
											}
											if !_rules[ruleCode]() {
												goto l12
											}
											if !_rules[ruleSuffix]() {
												goto l12
											}
										}
									}

								}
							l14:
								add(ruleStock, position13)
							}
							goto l5
						l12:
							position, tokenIndex = position5, tokenIndex5
							{
								position35 := position
								if buffer[position] != rune('$') {
									goto l34
								}
								position++
								if !_rules[ruleStockName]() {
									goto l34
								}
								if buffer[position] != rune('(') {
									goto l34
								}
								position++
								{
									position36, tokenIndex36 := position, tokenIndex
									if !_rules[ruleCNMarket]() {
										goto l37
									}
									if !_rules[ruleACode]() {
										goto l37
									}
									goto l36
								l37:
									position, tokenIndex = position36, tokenIndex36
									if !_rules[ruleHKMarket]() {
										goto l38
									}
									{
										position39 := position
										{
											position40, tokenIndex40 := position, tokenIndex
											if buffer[position] != rune('H') {
												goto l41
											}
											position++
											if buffer[position] != rune('S') {
												goto l41
											}
											position++
											if buffer[position] != rune('T') {
												goto l41
											}
											position++
											if buffer[position] != rune('E') {
												goto l41
											}
											position++
											if buffer[position] != rune('C') {
												goto l41
											}
											position++
											if buffer[position] != rune('H') {
												goto l41
											}
											position++
											goto l40
										l41:
											position, tokenIndex = position40, tokenIndex40
											if buffer[position] != rune('H') {
												goto l38
											}
											position++
											if buffer[position] != rune('S') {
												goto l38
											}
											position++
											if buffer[position] != rune('I') {
												goto l38
											}
											position++
										}
									l40:
										add(ruleHSCODE, position39)
									}
									goto l36
								l38:
									position, tokenIndex = position36, tokenIndex36
									if !_rules[ruleUSCode]() {
										goto l42
									}
									goto l36
								l42:
									position, tokenIndex = position36, tokenIndex36
									if !_rules[ruleHKCode]() {
										goto l34
									}
								}
							l36:
								if buffer[position] != rune(')') {
									goto l34
								}
								position++
								if buffer[position] != rune('$') {
									goto l34
								}
								position++
								add(ruleXLStock, position35)
							}
							goto l5
						l34:
							position, tokenIndex = position5, tokenIndex5
							{
								position43 := position
								if !matchDot() {
									goto l3
								}
								add(ruleOTHER, position43)
							}
						}
					l5:
						add(ruleLine, position4)
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position44, tokenIndex44 := position, tokenIndex
					if !matchDot() {
						goto l44
					}
					goto l0
				l44:
					position, tokenIndex = position44, tokenIndex44
				}
				add(ruleItem, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Line <- <(FTStock / Stock / XLStock / OTHER)> */
		nil,
		/* 2 OTHER <- <.> */
		nil,
		/* 3 Stock <- <(('$' Code (Suffix / Suffix?) '$') / ('$' Code Suffix) / ('(' (('N' 'Y' 'S' 'E') / ('N' 'A' 'S' 'D' 'A' 'Q')) ('：' / ':') SP* USCode ')') / ((&('(') ('(' Market ':' Code ')')) | (&('$') ('$' USCode)) | (&(' ') (' '+ Code Suffix))))> */
		nil,
		/* 4 XLStock <- <('$' StockName '(' ((CNMarket ACode) / (HKMarket HSCODE) / USCode / HKCode) ')' '$')> */
		nil,
		/* 5 FTStock <- <('$' (StockName / (StockName '-' Letter+)) '(' Code '.' Market ')' '$')> */
		nil,
		/* 6 StockName <- <(!((&('\'') [\'-\']) | (&(')') ')') | (&('(') '(')) .)+> */
		func() bool {
			position50, tokenIndex50 := position, tokenIndex
			{
				position51 := position
				{
					position54, tokenIndex54 := position, tokenIndex
					{
						switch buffer[position] {
						case '\'':
							if c := buffer[position]; c < rune('\'') || c > rune('\'') {
								goto l54
							}
							position++
						case ')':
							if buffer[position] != rune(')') {
								goto l54
							}
							position++
						default:
							if buffer[position] != rune('(') {
								goto l54
							}
							position++
						}
					}

					goto l50
				l54:
					position, tokenIndex = position54, tokenIndex54
				}
				if !matchDot() {
					goto l50
				}
			l52:
				{
					position53, tokenIndex53 := position, tokenIndex
					{
						position56, tokenIndex56 := position, tokenIndex
						{
							switch buffer[position] {
							case '\'':
								if c := buffer[position]; c < rune('\'') || c > rune('\'') {
									goto l56
								}
								position++
							case ')':
								if buffer[position] != rune(')') {
									goto l56
								}
								position++
							default:
								if buffer[position] != rune('(') {
									goto l56
								}
								position++
							}
						}

						goto l53
					l56:
						position, tokenIndex = position56, tokenIndex56
					}
					if !matchDot() {
						goto l53
					}
					goto l52
				l53:
					position, tokenIndex = position53, tokenIndex53
				}
				add(ruleStockName, position51)
			}
			return true
		l50:
			position, tokenIndex = position50, tokenIndex50
			return false
		},
		/* 7 Code <- <(USCode / HKCode / ACode)> */
		func() bool {
			position58, tokenIndex58 := position, tokenIndex
			{
				position59 := position
				{
					position60, tokenIndex60 := position, tokenIndex
					if !_rules[ruleUSCode]() {
						goto l61
					}
					goto l60
				l61:
					position, tokenIndex = position60, tokenIndex60
					if !_rules[ruleHKCode]() {
						goto l62
					}
					goto l60
				l62:
					position, tokenIndex = position60, tokenIndex60
					if !_rules[ruleACode]() {
						goto l58
					}
				}
			l60:
				add(ruleCode, position59)
			}
			return true
		l58:
			position, tokenIndex = position58, tokenIndex58
			return false
		},
		/* 8 USCode <- <Letter+> */
		func() bool {
			position63, tokenIndex63 := position, tokenIndex
			{
				position64 := position
				if !_rules[ruleLetter]() {
					goto l63
				}
			l65:
				{
					position66, tokenIndex66 := position, tokenIndex
					if !_rules[ruleLetter]() {
						goto l66
					}
					goto l65
				l66:
					position, tokenIndex = position66, tokenIndex66
				}
				add(ruleUSCode, position64)
			}
			return true
		l63:
			position, tokenIndex = position63, tokenIndex63
			return false
		},
		/* 9 HKCode <- <Number+> */
		func() bool {
			position67, tokenIndex67 := position, tokenIndex
			{
				position68 := position
				if !_rules[ruleNumber]() {
					goto l67
				}
			l69:
				{
					position70, tokenIndex70 := position, tokenIndex
					if !_rules[ruleNumber]() {
						goto l70
					}
					goto l69
				l70:
					position, tokenIndex = position70, tokenIndex70
				}
				add(ruleHKCode, position68)
			}
			return true
		l67:
			position, tokenIndex = position67, tokenIndex67
			return false
		},
		/* 10 ACode <- <Number+> */
		func() bool {
			position71, tokenIndex71 := position, tokenIndex
			{
				position72 := position
				if !_rules[ruleNumber]() {
					goto l71
				}
			l73:
				{
					position74, tokenIndex74 := position, tokenIndex
					if !_rules[ruleNumber]() {
						goto l74
					}
					goto l73
				l74:
					position, tokenIndex = position74, tokenIndex74
				}
				add(ruleACode, position72)
			}
			return true
		l71:
			position, tokenIndex = position71, tokenIndex71
			return false
		},
		/* 11 HSCODE <- <(('H' 'S' 'T' 'E' 'C' 'H') / ('H' 'S' 'I'))> */
		nil,
		/* 12 Letter <- <[A-Z]> */
		func() bool {
			position76, tokenIndex76 := position, tokenIndex
			{
				position77 := position
				if c := buffer[position]; c < rune('A') || c > rune('Z') {
					goto l76
				}
				position++
				add(ruleLetter, position77)
			}
			return true
		l76:
			position, tokenIndex = position76, tokenIndex76
			return false
		},
		/* 13 Number <- <[0-9]> */
		func() bool {
			position78, tokenIndex78 := position, tokenIndex
			{
				position79 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l78
				}
				position++
				add(ruleNumber, position79)
			}
			return true
		l78:
			position, tokenIndex = position78, tokenIndex78
			return false
		},
		/* 14 Suffix <- <('.' (Market / ('o' / 'O')))> */
		func() bool {
			position80, tokenIndex80 := position, tokenIndex
			{
				position81 := position
				if buffer[position] != rune('.') {
					goto l80
				}
				position++
				{
					position82, tokenIndex82 := position, tokenIndex
					if !_rules[ruleMarket]() {
						goto l83
					}
					goto l82
				l83:
					position, tokenIndex = position82, tokenIndex82
					{
						position84, tokenIndex84 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l85
						}
						position++
						goto l84
					l85:
						position, tokenIndex = position84, tokenIndex84
						if buffer[position] != rune('O') {
							goto l80
						}
						position++
					}
				l84:
				}
			l82:
				add(ruleSuffix, position81)
			}
			return true
		l80:
			position, tokenIndex = position80, tokenIndex80
			return false
		},
		/* 15 Market <- <(CNMarket / ((&('S') SGMarket) | (&('U') USMarket) | (&('H') HKMarket)))> */
		func() bool {
			position86, tokenIndex86 := position, tokenIndex
			{
				position87 := position
				{
					position88, tokenIndex88 := position, tokenIndex
					if !_rules[ruleCNMarket]() {
						goto l89
					}
					goto l88
				l89:
					position, tokenIndex = position88, tokenIndex88
					{
						switch buffer[position] {
						case 'S':
							{
								position91 := position
								if buffer[position] != rune('S') {
									goto l86
								}
								position++
								if buffer[position] != rune('G') {
									goto l86
								}
								position++
								add(ruleSGMarket, position91)
							}
						case 'U':
							{
								position92 := position
								if buffer[position] != rune('U') {
									goto l86
								}
								position++
								if buffer[position] != rune('S') {
									goto l86
								}
								position++
								add(ruleUSMarket, position92)
							}
						default:
							if !_rules[ruleHKMarket]() {
								goto l86
							}
						}
					}

				}
			l88:
				add(ruleMarket, position87)
			}
			return true
		l86:
			position, tokenIndex = position86, tokenIndex86
			return false
		},
		/* 16 CNMarket <- <(('S' 'H') / ('S' 'Z'))> */
		func() bool {
			position93, tokenIndex93 := position, tokenIndex
			{
				position94 := position
				{
					position95, tokenIndex95 := position, tokenIndex
					if buffer[position] != rune('S') {
						goto l96
					}
					position++
					if buffer[position] != rune('H') {
						goto l96
					}
					position++
					goto l95
				l96:
					position, tokenIndex = position95, tokenIndex95
					if buffer[position] != rune('S') {
						goto l93
					}
					position++
					if buffer[position] != rune('Z') {
						goto l93
					}
					position++
				}
			l95:
				add(ruleCNMarket, position94)
			}
			return true
		l93:
			position, tokenIndex = position93, tokenIndex93
			return false
		},
		/* 17 HKMarket <- <('H' 'K')> */
		func() bool {
			position97, tokenIndex97 := position, tokenIndex
			{
				position98 := position
				if buffer[position] != rune('H') {
					goto l97
				}
				position++
				if buffer[position] != rune('K') {
					goto l97
				}
				position++
				add(ruleHKMarket, position98)
			}
			return true
		l97:
			position, tokenIndex = position97, tokenIndex97
			return false
		},
		/* 18 USMarket <- <('U' 'S')> */
		nil,
		/* 19 SGMarket <- <('S' 'G')> */
		nil,
		/* 20 SP <- <(' ' / '\t')> */
		nil,
	}
	p.rules = _rules
	return nil
}
