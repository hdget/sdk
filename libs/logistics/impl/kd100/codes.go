package kd100

// shipperCodeMap 快递100常用快递公司编码映射（名称 -> 编码）
var shipperCodeMap = map[string]string{
	"顺丰速运":  "shunfeng",
	"百世快递":  "baishiwuliu",
	"中通快递":  "zhongtong",
	"申通快递":  "shentong",
	"圆通速递":  "yuantong",
	"韵达速递":  "yunda",
	"邮政快递包裹": "youzhengguonei",
	"EMS":     "ems",
	"京东快递":  "jd",
	"优速快递":  "youshuwuliu",
	"德邦快递":  "debangkuaidi",
	"极兔速递":  "jtexpress",
	"众邮快递":  "zhongyoukuaidi",
	"宅急送":   "zhaijisong",
}

// GetShipperCode 根据快递公司名称获取编码
func (a *api) GetShipperCode(name string) string {
	return shipperCodeMap[name]
}