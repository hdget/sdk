package kdniao

// shipperCodeMap 快递鸟常用快递公司编码映射（名称 -> 编码）
var shipperCodeMap = map[string]string{
	"顺丰速运":  "SF",
	"百世快递":  "HTKY",
	"中通快递":  "ZTO",
	"申通快递":  "STO",
	"圆通速递":  "YTO",
	"韵达速递":  "YD",
	"邮政快递包裹": "YZPY",
	"EMS":     "EMS",
	"京东快递":  "JD",
	"优速快递":  "UC",
	"德邦快递":  "DBL",
	"极兔速递":  "JTSD",
	"众邮快递":  "ZYE",
	"宅急送":   "ZJS",
}

// GetShipperCode 根据快递公司名称获取编码
func (a *api) GetShipperCode(name string) string {
	return shipperCodeMap[name]
}