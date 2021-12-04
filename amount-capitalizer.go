// 大写数字生成器
// 单位:
//       "",  拾,   佰,   仟
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

// 数字型值转为字符串
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

type dataGroup struct {
	start, end int
	bases []string
	unit    string
	lastInt bool
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

	nGroup := intLength / intCount
	gEnd := intLength % intCount
	if gEnd > 0 {
		nGroup += 1
	} else if nGroup > 0 {
		gEnd = intCount
	}
	if nGroup > len(units) {
		return makeTooLarge(isNeg)
	}

	groups := make(chan *dataGroup)
	res := convertGroups(groups, sAmount, isNeg)

	go func() {
		defer close(groups)

		nGroup -= 1
		for gStart := 0; nGroup >= 0; nGroup-- {
			groups <- &dataGroup {
				start: gStart,
				end: gEnd,
				bases: intBases,
				unit: units[nGroup],
				lastInt: gEnd>=intLength,
			}
			gStart = gEnd
			gEnd += intCount
		}
		groups <- &dataGroup{
			start: intLength + 1,
			end: length,
			bases: fracBases,
		}
	}()

	return res
}

func convertGroups(c <-chan *dataGroup, sAmount string, isNeg bool) (<-chan string) {
	res := make(chan string)

	go func() {
		defer close(res)
		if isNeg {
			res <- neg
		}

		prevZero := false
		allZero := true
		prevGroupAllZero := true

		for dg := range c {
			prevGroupAllZero = true
			gStart, gEnd, bases, unit := dg.start, dg.end, dg.bases, dg.unit
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
				prevZero, prevGroupAllZero = false, false
			}
			if !prevGroupAllZero || dg.lastInt {
				res <- unit
			}
		}

		if allZero {
			res <- zero
		}
		if prevGroupAllZero {
			res <- ending
		}
	}()

	return res
}

func makeTooLarge(isNeg bool) (<-chan string) {
	res := make(chan string)
	go func() {
		defer close(res)
		if isNeg {
			res <- tooSmall
		} else {
			res <- tooLarge
		}
	}()
	return res
}
