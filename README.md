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

### Benchmark in Go

```
Benchmark-8         	  489860	      2729 ns/op
Benchmark_long-8    	   64664	     21624 ns/op
Benchmark_large-8   	    1910	    589914 ns/op
```

## Development in Go

Use [https://github.com/pointlander/peg](https://github.com/pointlander/peg)

```sh
go install github.com/pointlander/peg
```

And then run `make` to generate `grammar.peg` into `grammar.go`.

> NOTE: Please do not change `grammar.go`.
