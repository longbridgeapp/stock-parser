package stock_parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/longbridgeapp/assert"
)

func assert_matches_code(t *testing.T, expected string, input string) {
	t.Helper()

	res := input

	cb := func(code, market, match string) string {
		counterId := fmt.Sprintf("ST/%s/%s", market, code)
		name := code
		s := fmt.Sprintf(`<span type="security-tag" counter_id="%s" name="%s">$%s.%s</span>`, counterId, name, name, market)
		res = strings.ReplaceAll(res, match, s)
		return res
	}

	out := Parse(input, cb)
	assert.Equal(t, expected, out)
}
func TestParse(t *testing.T) {
	assert_matches_code(t, `Alibaba <span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span> published its Q2 results`, "Alibaba BABA.US published its Q2 results")
	assert_matches_code(t, `Alibaba <span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span> published its Q2 results`, "Alibaba $BABA$ published its Q2 results")
	assert_matches_code(t, `阿里巴巴<span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span> 发布财报`, "阿里巴巴$BABA.US 发布财报")
	assert_matches_code(t, `阿里巴巴<span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span>发布财报`, "阿里巴巴$BABA.US$发布财报")
	assert_matches_code(t, `阿里巴巴 <span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span>$发布财报`, "阿里巴巴 BABA.US$发布财报")
	assert_matches_code(t, `阿里巴巴 <span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span> 发布财报`, "阿里巴巴 BABA.US 发布财报")
	assert_matches_code(t, "阿里巴巴 BABA$发布财报", "阿里巴巴 BABA$发布财报")
	assert_matches_code(t, "腾讯 700 发布财报", "腾讯 700 发布财报")
	assert_matches_code(t, "阿里巴巴 [BABA] 发布财报", "阿里巴巴 [BABA] 发布财报")
	assert_matches_code(t, "腾讯 (700) 发布财报", "腾讯 (700) 发布财报")
	assert_matches_code(t, `腾讯 <span type="security-tag" counter_id="ST/HK/00700" name="00700">$00700.HK</span> 发布财报`, "腾讯 00700.HK 发布财报")
	assert_matches_code(t, `Tesla Inc (<span type="security-tag" counter_id="ST/US/TSLA" name="TSLA">$TSLA.US</span>) will finalise a deal to invest in a production facility in his country`, "Tesla Inc (TSLA.O) will finalise a deal to invest in a production facility in his country")
	assert_matches_code(t, "Only the fortune of Tesla's (TSLA)", "Only the fortune of Tesla's (TSLA)")
	assert_matches_code(t, `阿里巴巴<span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span>和腾讯<span type="security-tag" counter_id="ST/HK/700" name="700">$700.HK</span>发布财报`, "阿里巴巴$BABA.US$和腾讯$700.HK$发布财报")

	assert_matches_code(t, `阿里巴巴 <span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span> 发布财报`, "阿里巴巴 $BABA 发布财报")

	assert_matches_code(t, `吉利交付 <span type="security-tag" counter_id="ST/HK/00175" name="00175">$00175.HK</span> 股票, 哈哈 <span type="security-tag" counter_id="ST/SH/603200" name="603200">$603200.SH</span> 哈哈`, `吉利交付 $00175.HK 股票, 哈哈 $上海洗霸(SH603200)$ 哈哈`)

}

func TestXueqiuLaohuFutu(t *testing.T) {
	// 老虎雪球
	assert_matches_code(t, `啊哈哈哈 <span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span> <span type="security-tag" counter_id="ST/SH/601318" name="601318">$601318.SH</span> <span type="security-tag" counter_id="ST/HK/00700" name="00700">$00700.HK</span> 这些天，江苏省100多家外贸企业在境外参展招商，开拓市场。这家大型显示设备的制造企业前期团队已在国外开展对接，这次他们带了20多项新技术。`, "啊哈哈哈 $阿里巴巴(BABA)$ $中国平安(SH601318)$ $腾讯控股(00700)$ 这些天，江苏省100多家外贸企业在境外参展招商，开拓市场。这家大型显示设备的制造企业前期团队已在国外开展对接，这次他们带了20多项新技术。")

	// 富途
	assert_matches_code(t, `<span type="security-tag" counter_id="ST/US/BABA" name="BABA">$BABA.US</span> 不错的哈哈哈  <span type="security-tag" counter_id="ST/HK/00700" name="00700">$00700.HK</span> 看好  <span type="security-tag" counter_id="ST/SZ/002241" name="002241">$002241.SZ</span> 也不错`, "$阿里巴巴(BABA.US)$ 不错的哈哈哈  $腾讯控股(00700.HK)$ 看好  $歌尔股份(002241.SZ)$ 也不错")
}

func TestSpecialMarket(t *testing.T) {
	assert_matches_code(t, `美股简短的股票代码 <span type="security-tag" counter_id="ST/US/Q" name="Q">$Q.US</span>, <span type="security-tag" counter_id="ST/US/AB" name="AB">$AB.US</span>`, "美股简短的股票代码 Q.US, AB.US")
}
