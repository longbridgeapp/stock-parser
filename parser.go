package stock_parser

const (
	MarketHK = "HK"
	MarketUS = "US"
)

func Parse(input string, cb func(code, market, match string) string) (out string) {
	parser := newParser(input)
	parser.consumeNode(parser.root, cb)

	if parser.out == "" {
		return input
	}

	return parser.out
}

type rn []rune

func (r rn) str(node *node32) string {
	if node == nil || int(node.begin) < 0 || int(node.end) > len(r) {
		return ""
	}
	return string(r[node.begin:node.end])
}

type parser struct {
	r     rn
	input string
	out   string
	root  *node32
}

func newParser(input string) *parser {
	p := StockCodeParser{Buffer: input}
	p.Init()

	if err := p.Parse(); err != nil {
		return nil
	}

	return &parser{
		input: input,
		r:     rn([]rune(input)),
		root:  p.AST(),
	}
}

func (p *parser) consumeNode(node *node32, cb func(node, market, match string) string) (out string) {
	for node != nil {
		code := ""
		market := ""
		match := ""
		if node.pegRule == ruleStock {
			code, market, match = p.parseStock(node)
		} else if node.pegRule == ruleXLStock {
			code, market, match = p.parseXLStock(node)
		} else if node.pegRule == ruleFTStock {
			code, market, match = p.parseFTStock(node)
		}

		if code != "" {
			if market == "" {
				market = MarketUS
			}
			p.out = cb(code, market, match)
		}

		if node.up != nil {
			p.consumeNode(node.up, cb)
		}

		node = node.next
	}
	return
}

func (p *parser) parseFTStock(node *node32) (code, market, match string) {
	r := p.r

	match = r.str(node)
	subNode := node.up
	for subNode != nil {
		switch subNode.pegRule {
		case ruleCode:
			code = r.str(subNode)
		case ruleMarket:
			market = r.str(subNode)
		}
		subNode = subNode.next
	}

	return
}

func (p *parser) parseXLStock(node *node32) (code, market, match string) {
	r := p.r

	match = r.str(node)
	subNode := node.up
	for subNode != nil {
		switch subNode.pegRule {
		case ruleUSCode, ruleACode:
			code = r.str(subNode)
		case ruleHKCode, ruleHSCODE:
			code = r.str(subNode)
			market = MarketHK
		case ruleMarket, ruleCNMarket, ruleHKMarket:
			market = r.str(subNode)
		}
		subNode = subNode.next
	}
	return
}

func (p *parser) parseStock(node *node32) (code, market, match string) {
	r := p.r

	match = r.str(node)
	subNode := node.up
	for subNode != nil {
		switch subNode.pegRule {
		case ruleCode, ruleUSCode:
			code = r.str(subNode)
		case ruleMarket:
			market = r.str(subNode)
		case ruleSuffix:
			market = r.str(subNode.up)
		}
		subNode = subNode.next
	}

	return
}
