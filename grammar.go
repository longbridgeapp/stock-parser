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
	rules  [21]func() bool
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
								{
									position8, tokenIndex8 := position, tokenIndex
									if buffer[position] != rune('$') {
										goto l9
									}
									position++
									if !_rules[ruleCode]() {
										goto l9
									}
									{
										position10, tokenIndex10 := position, tokenIndex
										if !_rules[ruleSuffix]() {
											goto l11
										}
										goto l10
									l11:
										position, tokenIndex = position10, tokenIndex10
										{
											position12, tokenIndex12 := position, tokenIndex
											if !_rules[ruleSuffix]() {
												goto l12
											}
											goto l13
										l12:
											position, tokenIndex = position12, tokenIndex12
										}
									l13:
									}
								l10:
									if buffer[position] != rune('$') {
										goto l9
									}
									position++
									goto l8
								l9:
									position, tokenIndex = position8, tokenIndex8
									{
										position15, tokenIndex15 := position, tokenIndex
										if buffer[position] != rune('$') {
											goto l15
										}
										position++
										goto l16
									l15:
										position, tokenIndex = position15, tokenIndex15
									}
								l16:
									if !_rules[ruleCode]() {
										goto l14
									}
									if !_rules[ruleSuffix]() {
										goto l14
									}
									goto l8
								l14:
									position, tokenIndex = position8, tokenIndex8
									if buffer[position] != rune('$') {
										goto l17
									}
									position++
									if !_rules[ruleUSCode]() {
										goto l17
									}
									goto l8
								l17:
									position, tokenIndex = position8, tokenIndex8
									if buffer[position] != rune('(') {
										goto l6
									}
									position++
									{
										position18, tokenIndex18 := position, tokenIndex
										if buffer[position] != rune('N') {
											goto l19
										}
										position++
										if buffer[position] != rune('Y') {
											goto l19
										}
										position++
										if buffer[position] != rune('S') {
											goto l19
										}
										position++
										if buffer[position] != rune('E') {
											goto l19
										}
										position++
										goto l18
									l19:
										position, tokenIndex = position18, tokenIndex18
										if buffer[position] != rune('N') {
											goto l6
										}
										position++
										if buffer[position] != rune('A') {
											goto l6
										}
										position++
										if buffer[position] != rune('S') {
											goto l6
										}
										position++
										if buffer[position] != rune('D') {
											goto l6
										}
										position++
										if buffer[position] != rune('A') {
											goto l6
										}
										position++
										if buffer[position] != rune('Q') {
											goto l6
										}
										position++
									}
								l18:
									{
										position20, tokenIndex20 := position, tokenIndex
										if buffer[position] != rune('：') {
											goto l21
										}
										position++
										goto l20
									l21:
										position, tokenIndex = position20, tokenIndex20
										if buffer[position] != rune(':') {
											goto l6
										}
										position++
									}
								l20:
								l22:
									{
										position23, tokenIndex23 := position, tokenIndex
										{
											position24 := position
											{
												position25, tokenIndex25 := position, tokenIndex
												if buffer[position] != rune(' ') {
													goto l26
												}
												position++
												goto l25
											l26:
												position, tokenIndex = position25, tokenIndex25
												if buffer[position] != rune('\t') {
													goto l23
												}
												position++
											}
										l25:
											add(ruleSP, position24)
										}
										goto l22
									l23:
										position, tokenIndex = position23, tokenIndex23
									}
									if !_rules[ruleUSCode]() {
										goto l6
									}
									if buffer[position] != rune(')') {
										goto l6
									}
									position++
								}
							l8:
								add(ruleStock, position7)
							}
							goto l5
						l6:
							position, tokenIndex = position5, tokenIndex5
							{
								position28 := position
								if buffer[position] != rune('$') {
									goto l27
								}
								position++
								if !_rules[ruleStockName]() {
									goto l27
								}
								if buffer[position] != rune('(') {
									goto l27
								}
								position++
								if !_rules[ruleCode]() {
									goto l27
								}
								if buffer[position] != rune('.') {
									goto l27
								}
								position++
								if !_rules[ruleMarket]() {
									goto l27
								}
								if buffer[position] != rune(')') {
									goto l27
								}
								position++
								if buffer[position] != rune('$') {
									goto l27
								}
								position++
								add(ruleFTStock, position28)
							}
							goto l5
						l27:
							position, tokenIndex = position5, tokenIndex5
							{
								position30 := position
								if buffer[position] != rune('$') {
									goto l29
								}
								position++
								if !_rules[ruleStockName]() {
									goto l29
								}
								if buffer[position] != rune('(') {
									goto l29
								}
								position++
								{
									position31, tokenIndex31 := position, tokenIndex
									if !_rules[ruleCNMarket]() {
										goto l32
									}
									if !_rules[ruleACode]() {
										goto l32
									}
									goto l31
								l32:
									position, tokenIndex = position31, tokenIndex31
									if !_rules[ruleUSCode]() {
										goto l33
									}
									goto l31
								l33:
									position, tokenIndex = position31, tokenIndex31
									if !_rules[ruleHKCode]() {
										goto l29
									}
								}
							l31:
								if buffer[position] != rune(')') {
									goto l29
								}
								position++
								if buffer[position] != rune('$') {
									goto l29
								}
								position++
								add(ruleXLStock, position30)
							}
							goto l5
						l29:
							position, tokenIndex = position5, tokenIndex5
							{
								position34 := position
								if !matchDot() {
									goto l3
								}
								add(ruleOTHER, position34)
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
					position35, tokenIndex35 := position, tokenIndex
					if !matchDot() {
						goto l35
					}
					goto l0
				l35:
					position, tokenIndex = position35, tokenIndex35
				}
				add(ruleItem, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 Line <- <(Stock / FTStock / XLStock / OTHER)> */
		nil,
		/* 2 OTHER <- <.> */
		nil,
		/* 3 Stock <- <(('$' Code (Suffix / Suffix?) '$') / ('$'? Code Suffix) / ('$' USCode) / ('(' (('N' 'Y' 'S' 'E') / ('N' 'A' 'S' 'D' 'A' 'Q')) ('：' / ':') SP* USCode ')'))> */
		nil,
		/* 4 XLStock <- <('$' StockName '(' ((CNMarket ACode) / USCode / HKCode) ')' '$')> */
		nil,
		/* 5 FTStock <- <('$' StockName '(' Code '.' Market ')' '$')> */
		nil,
		/* 6 StockName <- <(!((&('P') 'P') | (&('S') 'S') | (&(')') ')') | (&('(') '(')) .)+> */
		func() bool {
			position41, tokenIndex41 := position, tokenIndex
			{
				position42 := position
				{
					position45, tokenIndex45 := position, tokenIndex
					{
						switch buffer[position] {
						case 'P':
							if buffer[position] != rune('P') {
								goto l45
							}
							position++
						case 'S':
							if buffer[position] != rune('S') {
								goto l45
							}
							position++
						case ')':
							if buffer[position] != rune(')') {
								goto l45
							}
							position++
						default:
							if buffer[position] != rune('(') {
								goto l45
							}
							position++
						}
					}

					goto l41
				l45:
					position, tokenIndex = position45, tokenIndex45
				}
				if !matchDot() {
					goto l41
				}
			l43:
				{
					position44, tokenIndex44 := position, tokenIndex
					{
						position47, tokenIndex47 := position, tokenIndex
						{
							switch buffer[position] {
							case 'P':
								if buffer[position] != rune('P') {
									goto l47
								}
								position++
							case 'S':
								if buffer[position] != rune('S') {
									goto l47
								}
								position++
							case ')':
								if buffer[position] != rune(')') {
									goto l47
								}
								position++
							default:
								if buffer[position] != rune('(') {
									goto l47
								}
								position++
							}
						}

						goto l44
					l47:
						position, tokenIndex = position47, tokenIndex47
					}
					if !matchDot() {
						goto l44
					}
					goto l43
				l44:
					position, tokenIndex = position44, tokenIndex44
				}
				add(ruleStockName, position42)
			}
			return true
		l41:
			position, tokenIndex = position41, tokenIndex41
			return false
		},
		/* 7 Code <- <(USCode / HKCode / ACode)> */
		func() bool {
			position49, tokenIndex49 := position, tokenIndex
			{
				position50 := position
				{
					position51, tokenIndex51 := position, tokenIndex
					if !_rules[ruleUSCode]() {
						goto l52
					}
					goto l51
				l52:
					position, tokenIndex = position51, tokenIndex51
					if !_rules[ruleHKCode]() {
						goto l53
					}
					goto l51
				l53:
					position, tokenIndex = position51, tokenIndex51
					if !_rules[ruleACode]() {
						goto l49
					}
				}
			l51:
				add(ruleCode, position50)
			}
			return true
		l49:
			position, tokenIndex = position49, tokenIndex49
			return false
		},
		/* 8 USCode <- <Letter+> */
		func() bool {
			position54, tokenIndex54 := position, tokenIndex
			{
				position55 := position
				{
					position58 := position
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l54
					}
					position++
					add(ruleLetter, position58)
				}
			l56:
				{
					position57, tokenIndex57 := position, tokenIndex
					{
						position59 := position
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l57
						}
						position++
						add(ruleLetter, position59)
					}
					goto l56
				l57:
					position, tokenIndex = position57, tokenIndex57
				}
				add(ruleUSCode, position55)
			}
			return true
		l54:
			position, tokenIndex = position54, tokenIndex54
			return false
		},
		/* 9 HKCode <- <Number+> */
		func() bool {
			position60, tokenIndex60 := position, tokenIndex
			{
				position61 := position
				if !_rules[ruleNumber]() {
					goto l60
				}
			l62:
				{
					position63, tokenIndex63 := position, tokenIndex
					if !_rules[ruleNumber]() {
						goto l63
					}
					goto l62
				l63:
					position, tokenIndex = position63, tokenIndex63
				}
				add(ruleHKCode, position61)
			}
			return true
		l60:
			position, tokenIndex = position60, tokenIndex60
			return false
		},
		/* 10 ACode <- <Number+> */
		func() bool {
			position64, tokenIndex64 := position, tokenIndex
			{
				position65 := position
				if !_rules[ruleNumber]() {
					goto l64
				}
			l66:
				{
					position67, tokenIndex67 := position, tokenIndex
					if !_rules[ruleNumber]() {
						goto l67
					}
					goto l66
				l67:
					position, tokenIndex = position67, tokenIndex67
				}
				add(ruleACode, position65)
			}
			return true
		l64:
			position, tokenIndex = position64, tokenIndex64
			return false
		},
		/* 11 Letter <- <[A-Z]> */
		nil,
		/* 12 Number <- <[0-9]> */
		func() bool {
			position69, tokenIndex69 := position, tokenIndex
			{
				position70 := position
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l69
				}
				position++
				add(ruleNumber, position70)
			}
			return true
		l69:
			position, tokenIndex = position69, tokenIndex69
			return false
		},
		/* 13 Suffix <- <('.' (Market / ('o' / 'O')))> */
		func() bool {
			position71, tokenIndex71 := position, tokenIndex
			{
				position72 := position
				if buffer[position] != rune('.') {
					goto l71
				}
				position++
				{
					position73, tokenIndex73 := position, tokenIndex
					if !_rules[ruleMarket]() {
						goto l74
					}
					goto l73
				l74:
					position, tokenIndex = position73, tokenIndex73
					{
						position75, tokenIndex75 := position, tokenIndex
						if buffer[position] != rune('o') {
							goto l76
						}
						position++
						goto l75
					l76:
						position, tokenIndex = position75, tokenIndex75
						if buffer[position] != rune('O') {
							goto l71
						}
						position++
					}
				l75:
				}
			l73:
				add(ruleSuffix, position72)
			}
			return true
		l71:
			position, tokenIndex = position71, tokenIndex71
			return false
		},
		/* 14 Market <- <(CNMarket / ((&('S') SGMarket) | (&('U') USMarket) | (&('H') HKMarket)))> */
		func() bool {
			position77, tokenIndex77 := position, tokenIndex
			{
				position78 := position
				{
					position79, tokenIndex79 := position, tokenIndex
					if !_rules[ruleCNMarket]() {
						goto l80
					}
					goto l79
				l80:
					position, tokenIndex = position79, tokenIndex79
					{
						switch buffer[position] {
						case 'S':
							{
								position82 := position
								if buffer[position] != rune('S') {
									goto l77
								}
								position++
								if buffer[position] != rune('G') {
									goto l77
								}
								position++
								add(ruleSGMarket, position82)
							}
						case 'U':
							{
								position83 := position
								if buffer[position] != rune('U') {
									goto l77
								}
								position++
								if buffer[position] != rune('S') {
									goto l77
								}
								position++
								add(ruleUSMarket, position83)
							}
						default:
							{
								position84 := position
								if buffer[position] != rune('H') {
									goto l77
								}
								position++
								if buffer[position] != rune('K') {
									goto l77
								}
								position++
								add(ruleHKMarket, position84)
							}
						}
					}

				}
			l79:
				add(ruleMarket, position78)
			}
			return true
		l77:
			position, tokenIndex = position77, tokenIndex77
			return false
		},
		/* 15 CNMarket <- <(('S' 'H') / ('S' 'Z'))> */
		func() bool {
			position85, tokenIndex85 := position, tokenIndex
			{
				position86 := position
				{
					position87, tokenIndex87 := position, tokenIndex
					if buffer[position] != rune('S') {
						goto l88
					}
					position++
					if buffer[position] != rune('H') {
						goto l88
					}
					position++
					goto l87
				l88:
					position, tokenIndex = position87, tokenIndex87
					if buffer[position] != rune('S') {
						goto l85
					}
					position++
					if buffer[position] != rune('Z') {
						goto l85
					}
					position++
				}
			l87:
				add(ruleCNMarket, position86)
			}
			return true
		l85:
			position, tokenIndex = position85, tokenIndex85
			return false
		},
		/* 16 HKMarket <- <('H' 'K')> */
		nil,
		/* 17 USMarket <- <('U' 'S')> */
		nil,
		/* 18 SGMarket <- <('S' 'G')> */
		nil,
		/* 19 SP <- <(' ' / '\t')> */
		nil,
	}
	p.rules = _rules
	return nil
}
