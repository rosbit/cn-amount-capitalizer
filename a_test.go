package cnamount

import (
	"testing"
	"fmt"
)

func testAmount(f interface{}) {
	sAmount, _, ok := FormatAmount(f)
	if !ok {
		fmt.Printf("%v => %s\n", f, sAmount)
	} else {
		fmt.Printf("%s => %s\n", sAmount, ToCNAmount(f))
	}
}

func TestCNAmount(t *testing.T) {
	testAmount(435235324532.0) // 元整
	testAmount(100.02)         // 佰元零x分
	testAmount(100.2)          // 佰元零x角
	testAmount(10.23)          // 拾元零x角x分
	testAmount(101.23)         // 元x角x分
	testAmount(101.03)         // 零x分
	testAmount(float32(10100.02))
	testAmount(340210100.02)
	testAmount(-340210100.02)  // 负
	testAmount(3400000000.02)
	testAmount(4352352343400000000.02)    // 万万亿
	testAmount(14352352343400000000.02)   // float64已不精准
	testAmount(214352352343400000000.02)  // 超大数额
	testAmount(-214352352343400000000.02) // 超小数额
	testAmount(.00)    // 零元
	testAmount(.12)    // x角x分
	testAmount(9999)   // 元整
	testAmount(19800)  // 佰元整
	testAmount(2980)   // 拾元整
	testAmount(500200)
	testAmount(103)    // 佰零x元整
	testAmount(int16(32766))
	testAmount(-100.23)
	testAmount(0)
}

