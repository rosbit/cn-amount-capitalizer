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
	units = []string{"圆", "万", "亿", "万亿", "万万亿"}
	intBases = []string{"", "拾", "佰", "仟"}
	digits = []string{"零","壹","贰","叁","肆","伍","陆","柒","捌","玖"}
	fracBases = []string{"分", "角"}
	ending = "整"
	zero = "零圆"
	sep = "零"
	tooLarge = "超大数额"

	intCount = len(intBases)
)

// 把数字金额转换成中文大写金额
func ToCNAmount(amount float64) string {
	sAmount := fmt.Sprintf("%.2f", amount)
	pointPos := strings.IndexByte(sAmount, '.')

	switch {
	case pointPos > 0:
		// 有整有零
		intRes, intAllZero := convertInteger(sAmount[:pointPos])
		fracRes, fracAllZero, fracStartingZero := convertGroup(sAmount[pointPos+1:], fracBases)
		if intAllZero {
			if fracAllZero {
				return zero
			}
			return fracRes
		}
		if fracAllZero {
			return fmt.Sprintf("%s%s", intRes, ending)
		}
		if fracStartingZero {
			return fmt.Sprintf("%s%s%s", intRes, sep, fracRes)
		}
		return fmt.Sprintf("%s%s", intRes, fracRes)
	case pointPos == 0:
		// 全是小数
		res, _, _ := convertGroup(sAmount[1:], fracBases)
		return res
	default:
		// 全是整数
		res, intAllZero := convertInteger(sAmount)
		if intAllZero {
			return zero
		}
		return fmt.Sprintf("%s%s", res, ending)
	}
}

// 把整数转为中文
// @param intAmount 金额的整数部分
// @return res 中文结果
//         allZero 是否全零?
func convertInteger(intAmount string) (res string, allZero bool) {
	var gStart, gEnd int

	length := len(intAmount)
	nGroup := length / intCount
	remLen := length % intCount
	if remLen > 0 {
		nGroup += 1
		gEnd = remLen
	} else {
		gEnd = intCount
	}
	if nGroup > len(units) {
		res = tooLarge
		return
	}

	r := &strings.Builder{}
	allZero = true
	endingZero := false
	groupIndex := nGroup - 1
	for {
		gRes, gAllZero, gStartingZero := convertGroup(intAmount[gStart:gEnd], intBases)
		if !gAllZero {
			if !allZero && (gStartingZero || endingZero) {
				r.WriteString(sep)
			}
			allZero = false
			r.WriteString(gRes)
			r.WriteString(units[groupIndex])
		}
		if groupIndex == 0 {
			if !allZero && gAllZero {
				r.WriteString(units[groupIndex])
			}
			break
		}
		endingZero = intAmount[gEnd-1] == '0'
		groupIndex -= 1
		gStart = gEnd
		gEnd += intCount
	}
	res = r.String()
	return
}

// 基于某个进制对一节数转中文
// @param dGroup 一节数字，整数部分最多4位，小数点后数字最多2位
// @param bases  进制
// @return res 结果
//         allZero 是否全零?
//         startingZero 是否零打头
func convertGroup(dGroup string, bases []string) (res string, allZero bool, startingZero bool) {
	l := len(dGroup)
	switch {
	case l == len(bases):
		startingZero = (dGroup[0] == '0')
	case l == 0:
		return
	default:
		startingZero = true
	}

	allZero = true
	prevZero := true
	r := &strings.Builder{}
	for i, d := range dGroup {
		if d == '0' {
			if !prevZero {
				prevZero = true
			}
			continue
		}

		p := l - i - 1
		if allZero {
			allZero = false
		} else if prevZero {
			r.WriteString(sep)
		}
		r.WriteString(digits[d - '0'])
		r.WriteString(bases[p])
		prevZero = false
	}
	res = r.String()
	return
}
