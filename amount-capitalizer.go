// 大写数字生成器
// 单位:
//   "",  拾,   佰,   仟
//   圆
//   万 
//   亿
//   万亿
//   万万亿
//
// 小数点后单位：角 分 
// 
// 数字: 零,壹,贰,叁,肆,伍,陆,柒,捌,玖

package cnamount

import (
	"fmt"
	"strings"
)

var (
	units = []string{"圆", "万", "亿", "万亿", "万万亿"/*, ...*/}  // 更大的数字单位虽然可以往后增加，但float64已经不精准了
	intBases = []string{"", "拾", "佰", "仟"}
	digits = []string{"零","壹","贰","叁","肆","伍","陆","柒","捌","玖"}
	fracBases = []string{/*..., */"分", "角"} // 更小的金额单位可往前面增加
	ending = "整"
	zero = "零圆"
	sep = "零"
	tooLarge = "超大数额"

	intCount = len(intBases)
)

// 把数字金额转换成中文大写金额
func ToCNAmount(amount float64) string {
	floatFormat := fmt.Sprintf("%%.%df", len(fracBases))
	sAmount := fmt.Sprintf(floatFormat, amount)
	res := convert(sAmount)
	r := &strings.Builder{}
	for d := range res {
		r.WriteString(d)
	}
	return r.String()
}

func convert(sAmount string) (<-chan string) {
	length := len(sAmount)
	intLength := strings.IndexByte(sAmount, '.')

	var gStart, gEnd int
	var unit string
	bases := intBases
	nGroup := intLength / intCount
	remLen := intLength % intCount
	if remLen > 0 {
		nGroup += 1
		gEnd = remLen
	} else if nGroup == 0 {
		gStart = intLength + 1
		gEnd = length
		bases = fracBases
	} else {
		gEnd = intCount
	}
	nGroup -= 1

	res := make(chan string)
	go func() {
		defer close(res)
		if nGroup >= len(units) {
			res <- tooLarge
			return
		}
		if nGroup >= 0 {
			unit = units[nGroup]
		}

		prevZero := false
		allZero := true
		for gStart < length {
			gAllZero := true
			for i, idx := gStart, gEnd - gStart - 1; i<gEnd; i, idx = i+1, idx-1 {
				d := sAmount[i]
				if d == '0' {
					if !prevZero {
						prevZero = true
					}
					continue
				}

				if allZero {
					allZero = false
				} else if prevZero {
					res <- sep
				}

				res <- digits[d - '0']
				res <- bases[idx]
				prevZero = false
				gAllZero = false
			}
			if !allZero {
				res <- unit
			}

			gStart = gEnd
			if gStart < intLength {
				gEnd += intCount
				nGroup -= 1
				unit = units[nGroup]
				continue
			}

			gStart += 1
			if gStart < length {
				gEnd = length
				bases = fracBases
				unit = ""
				continue
			}

			if allZero {
				res <- zero
			} else if gAllZero {
				res <- ending
			}
			break
		}
	}()

	return res
}

