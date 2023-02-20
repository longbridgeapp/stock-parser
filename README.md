# Stock Parser

This is [Ticker symbol](https://en.wikipedia.org/wiki/Ticker_symbol) parse example to show how to write a parser in Go and Rust (Pest, Nom, rust-peg).
And run benchmarks to compare the performance.

## Usage

```go
// social-service example
package main

import (
	"fmt"

	stock_parser "github.com/longbridgeapp/stock-parser"
)

func main() {
	out := Parse("看好 $BABA 和 $XPEV 未来有好的增长")
	fmt.Println(out)
	// 看好 <span type="security-tag" counter_id="ST/US/BABA" name="阿里巴巴">$阿里巴巴.US</span> 和 <span type="security-tag" counter_id="ST/US/XPEV" name="小鹏汽车">$X小鹏汽车.US</span> 未来有好的增长
}


func Parse(ctx context.Context, body string) string {
  // preproccess
  out := body
  body = preprocess(body)

  // Use with StockInfo SDK
	_ = stock_parser.Parse(body, func(code, market, match string) string {
		if stock, ok := StockInfoSDK.GetCounterId(code, market); ok {
			s := fmt.Sprintf(`<span type="security-tag" counter_id="%s" name="%s">$%s.%s</span>`, stock.CounterId, stock.Name, stock.Name, stock.Market)
			// 替换
			out = strings.ReplaceAll(out, match, s)
		}
		return out
	})

	return out
}
```

## Parsers

- [Pest](https://pest.rs)
- [Nom](https://github.com/rust-bakery/nom)
- [rust-peg](https://github.com/kevinmehall/rust-peg)
- [tdewolff/parse](github.com/tdewolff/parse) - (Go)

### Benchmark in Go

```
Benchmark_tdewolff_parse-8         	  489860	      2729 ns/op
Benchmark_tdewolff_parse_long-8    	   64664	     21624 ns/op
Benchmark_tdewolff_parse_large-8   	    1910	    589914 ns/op
```

### Benchmark in Rust

```
pest_parse              time:   [1.6484 µs 1.6730 µs 1.7062 µs]
pest_parse_long         time:   [16.631 µs 16.912 µs 17.265 µs]
pest_parse_large        time:   [571.65 µs 590.45 µs 616.38 µs]

nom_parse               time:   [420.42 ns 423.97 ns 428.22 ns]
nom_parse_long          time:   [3.7042 µs 3.7315 µs 3.7676 µs]
nom_parse_large         time:   [107.76 µs 109.42 µs 112.08 µs]

peg_parse               time:   [1.0119 µs 1.0189 µs 1.0268 µs]
peg_parse_long          time:   [9.7554 µs 9.9865 µs 10.388 µs]
peg_parse_large         time:   [331.86 µs 342.67 µs 358.57 µs]
```

## Development in Go

Use [https://github.com/pointlander/peg](https://github.com/pointlander/peg)

```
go install github.com/pointlander/peg
```

And then run `make` to generate `grammar.peg` into `grammar.go`.

> NOTE: Please do not change `grammar.go`.
