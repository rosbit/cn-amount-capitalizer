// 大写数字生成器
// 单位:
//   "",  拾,   佰,   仟
//   元
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
	"reflect"
	"fmt"
	"strings"
)

var (
	units = []string{"元", "万", "亿", "万亿", "万万亿"/*, ...*/}  // 更大的数字单位虽然可以往后增加，但float64已经不精准了
	intBases = []string{"", "拾", "佰", "仟"}
	digits = []string{"零","壹","贰","叁","肆","伍","陆","柒","捌","玖"}
	fracBases = []string{/*..., */"分", "角"} // 更小的金额单位可往前面增加
	ending = "整"
	zero = "零元"
	sep = "零"
	neg = "负"
	tooLarge = "超大数额"
	tooSmall = "超小数额"
	unsupporting = "数据类型不支持"

	intCount = len(intBases)
)

// 把数字金额转换成中文大写金额
func ToCNAmount(amount interface{}) string {
	sAmount, isNeg, ok := FormatAmount(amount)
	if !ok {
		return sAmount
	}

	res := convert(sAmount, isNeg)
	r := &strings.Builder{}
	for d := range res {
		r.WriteString(d)
	}
	return r.String()
}

func FormatAmount(amount interface{}) (sAmount string, isNeg, ok bool) {
	switch amount.(type) {
	case int, int8, int16, int32, int64:
		a := reflect.ValueOf(amount).Int()
		isNeg = a < 0
		sAmount = fmt.Sprintf("%d", a)
	case uint, uint8, uint16, uint32, uint64:
		sAmount = fmt.Sprintf("%v", amount)
	case float64, float32:
		a := reflect.ValueOf(amount).Float()
		isNeg = a < 0.0
		floatFormat := fmt.Sprintf("%%.%df", len(fracBases))
		sAmount = fmt.Sprintf(floatFormat, a)
	default:
		sAmount = unsupporting
		return
	}
	ok = true
	return
}

func convert(sAmount string, isNeg bool) (<-chan string) {
	if isNeg {
		sAmount = sAmount[1:]
	}
	length := len(sAmount)
	intLength := strings.IndexByte(sAmount, '.')
	if intLength < 0 {
		intLength = length
	}

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
			if isNeg {
				res <- tooSmall
			} else {
				res <- tooLarge
			}
			return
		}
		if nGroup >= 0 {
			unit = units[nGroup]
		}

		if isNeg {
			res <- neg
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
				prevZero, gAllZero = false, false
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
			}
			if gAllZero || intLength == length {
				res <- ending
			}
			break
		}
	}()

	return res
}

