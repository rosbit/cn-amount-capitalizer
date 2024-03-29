## 数字金额转中文大写金额

 - 可以转换到`万万亿`级别
 - 如有需要，可以在单位中追加更高级别
 - 该项目已停止更新，所有功能已经合并到[cn-capitalizer](https://github.com/rosbit/cn-capitalizer).

### 使用方法
 - 见测试程序`a_test.go`
   ```bash
   $ go test
   435235324532.00 => 肆仟叁佰伍拾贰亿叁仟伍佰叁拾贰万肆仟伍佰叁拾贰元整
   100.02 => 壹佰元零贰分
   100.20 => 壹佰元零贰角
   10.23 => 壹拾元零贰角叁分
   101.23 => 壹佰零壹元贰角叁分
   101.03 => 壹佰零壹元零叁分
   10100.02 => 壹万零壹佰元零贰分
   340210100.02 => 叁亿肆仟零贰拾壹万零壹佰元零贰分
   -340210100.02 => 负叁亿肆仟零贰拾壹万零壹佰元零贰分
   3400000000.02 => 叁拾肆亿万元零贰分
   4352352343400000000.00 => 肆佰叁拾伍万万亿贰仟叁佰伍拾贰万亿叁仟肆佰叁拾肆亿万元整
   14352352343399999488.00 => 壹仟肆佰叁拾伍万万亿贰仟叁佰伍拾贰万亿叁仟肆佰叁拾叁亿玖仟玖佰玖拾玖万玖仟肆佰捌拾捌元整
   214352352343399989248.00 => 超大数额
   -214352352343399989248.00 => 超小数额
   0.00 => 零元整
   0.12 => 壹角贰分
   9999 => 玖仟玖佰玖拾玖元整
   19800 => 壹万玖仟捌佰元整
   2980 => 贰仟玖佰捌拾元整
   500200 => 伍拾万零贰佰元整
   103 => 壹佰零叁元整
   32766 => 叁万贰仟柒佰陆拾陆元整
   -100.23 => 负壹佰元零贰角叁分
   0 => 零元整
   PASS
   ok  	github.com/rosbit/cn-amount-capitalizer	0.101s
   ```
