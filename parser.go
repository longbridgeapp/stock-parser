package stock_parser

func Parse(input string, cb func(code, market, match string) string) (out string) {
	parser := StockCodeParser{Buffer: input}
	parser.Init()

	runeBuff := []rune(input)

	str := func(node *node32) string {
		if node == nil || int(node.begin) < 0 || int(node.end) > len(runeBuff) {
			return ""
		}
		return string(runeBuff[node.begin:node.end])
	}

	if err := parser.Parse(); err != nil {
		return
	}

	node := parser.AST()

	var cunsumeNode func(node *node32)

	cunsumeNode = func(node *node32) {
		for node != nil {
			code := ""
			market := ""
			match := ""

			if node.pegRule == ruleStock {
				// fmt.Println("ruleStock", node.begin, node.end, str(node))

				match = str(node)

				sub_node := node.up

				for sub_node != nil {
					switch sub_node.pegRule {
					case ruleCode, ruleUSCode:
						code = str(sub_node)
					case ruleMarket:
						market = str(sub_node)
					case ruleSuffix:
						market = str(sub_node.up)
					}
					sub_node = sub_node.next
				}
			}

			if code != "" {
				if market == "" {
					market = "US"
				}
				out = cb(code, market, match)
			}

			if node.up != nil {
				cunsumeNode(node.up)
			}

			node = node.next
		}
	}

	cunsumeNode(node)

	if out == "" {
		return input
	}

	return
}
