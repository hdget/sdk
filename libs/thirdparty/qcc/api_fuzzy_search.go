package qcc

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hdget/sdk/common/protobuf"
	"github.com/pkg/errors"
)

const (
	urlFuzzySearch = "/FuzzySearch/GetList" // 企业模糊搜索接口地址
	networkTimeout = 10 * time.Second       // 网络请求超时时间
)

// 有效请求状态码
const (
	StatusOK                      = "200" // 查询成功
	StatusNoResult                = "201" // 查询无结果
	StatusParamError              = "202" // 查询参数错误，请检查
	StatusPending                 = "205" // 等待处理中
	StatusTooManyRecords          = "207" // 请求数据的条目数超过上限（5000）
	StatusUnsupportedCompanyType  = "208" // 此接口不支持此公司类型查询
	StatusTooManyCompanies        = "209" // 企业数量超过上限
	StatusParamTooShort           = "213" // 参数长度不能小于 2
	StatusUnsupportedKeyword      = "215" // 不支持的查询关键字
	StatusUnsupportedShellScan    = "218" // 该企业暂不支持空壳扫描
	StatusUnsupportedDueDiligence = "219" // 该企业暂不支持准入尽调
	StatusAPIOffline              = "105" // 接口已下线停用
	StatusRepeatedError           = "110" // 当前相同查询连续出错，请等 2 小时后重试
)

// 无效请求状态码
const (
	StatusInvalidKey               = "101" // 当前的 KEY 无效或者还未生效中
	StatusKeyArrears               = "102" // 当前 KEY 已欠费
	StatusKeySuspended             = "103" // 当前 KEY 被暂停使用
	StatusKeyError                 = "104" // 请求 KEY 异常，请联系管理员
	StatusTooManyIllegalRequests   = "106" // 非法请求过多，请联系管理员
	StatusForbiddenIPOrSignError   = "107" // 被禁止的 IP 或者签名错误
	StatusTooManyAbnormalRequests  = "108" // 异常请求过多，请联系管理员
	StatusExceedDailyLimit         = "109" // 请求超过每日系统限制
	StatusPermissionDenied         = "111" // 接口权限未开通，请联系管理员
	StatusQuotaExhausted           = "112" // 您的账号剩余使用量已不足或已过期
	StatusAPIDeleted               = "113" // 当前接口已被删除，请重新申请
	StatusAPIDisabled              = "114" // 当前接口已被禁用，请联系管理员
	StatusAuthExpired              = "115" // 身份验证错误或者已过期
	StatusExceedTotalLimit         = "116" // 请求超过每日（服务期）调用总量限制
	StatusTooManyUnsupportedParams = "117" // 当前不支持的请求参数调用量过多
	StatusUnsupportedCallMethod    = "118" // 当前接口不支持此方式的调用
	StatusAccountAbnormal          = "119" // 您的账号出现异常，请联系管理员
	StatusTrafficAbnormal          = "120" // 系统流量异常，请稍后再试
	StatusDataCrossBorder          = "121" // 数据不能出境
	StatusExceedConcurrencyLimit   = "122" // 请求超过系统并发限制
	StatusUnverifiedLimit          = "123" // 您的请求已达未认证权益上限，请及时认证
	StatusSystemMaintenance        = "124" // 系统维护中，暂时无法下单
	StatusInvalidQueryParam        = "125" // 查询参数错误
	StatusInvalidParamType         = "126" // 查询参数类型或结构有误
	StatusUnknownError             = "199" // 系统未知错误，请联系技术客服
	StatusQueryAbnormal            = "203" // 系统查询有异常，请联系技术人员
	StatusAPINotPurchased          = "214" // 您还未购买过该接口，请先购买
	StatusAmbiguousParam           = "223" // 查询参数模糊，无法获取结果
	StatusInvalidParam             = "224" // 查询参数无效
)

// Paging 分页信息
type Paging struct {
	PageSize     int `json:"PageSize"`     // 每页记录数
	PageIndex    int `json:"PageIndex"`    // 当前页码
	TotalRecords int `json:"TotalRecords"` // 总记录数
}

// Enterprise 企业基本信息
type Enterprise struct {
	KeyNo      string `json:"KeyNo"`      // 主键
	Name       string `json:"Name"`       // 企业名称
	CreditCode string `json:"CreditCode"` // 统一社会信用代码（中国香港企业返回商业登记号码）
	StartDate  string `json:"StartDate"`  // 成立日期
	OperName   string `json:"OperName"`   // 法定代表人姓名
	Status     string `json:"Status"`     // 企业状态
	No         string `json:"No"`         // 注册号（中国香港企业返回企业编号）
	Address    string `json:"Address"`    // 注册地址
}

// FuzzySearchResponse 企业模糊搜索响应
type FuzzySearchResponse struct {
	Paging      Paging       `json:"Paging"`      // 分页信息
	Result      []Enterprise `json:"Result"`      // 企业列表
	Status      string       `json:"Status"`      // 请求状态码
	Message     string       `json:"Message"`     // 状态说明
	OrderNumber string       `json:"OrderNumber"` // 订单号
}

// FuzzySearch 企业模糊搜索
// searchKey: 搜索关键词（支持企业名、人名、产品名、地址、电话、经营范围等）
// pageIndex: 页码，可选，默认第 1 页
// 仅当状态码为 200（查询成功）或 201（查询无结果）时正常返回，其余状态码返回 error
func (impl *apiImpl) FuzzySearch(searchKey string, list ...*protobuf.ListParam) (int64, []Enterprise, error) {
	request := resty.New().SetTimeout(networkTimeout).R().
		SetQueryParams(map[string]string{
			"key":       impl.key,
			"searchKey": searchKey,
		})

	pageIndex := 1
	if len(list) > 0 && list[0] != nil {
		pageIndex = int(list[0].Page)
	}

	request.SetQueryParam("pageIndex", strconv.Itoa(pageIndex))
	for name, value := range impl.authHeaders() {
		request.SetHeader(name, value)
	}

	resp, err := request.Get(impl.baseURL + urlFuzzySearch)
	if err != nil {
		return 0, nil, errors.Wrap(err, "fuzzy search request")
	}

	var ret FuzzySearchResponse
	if err = json.Unmarshal(resp.Body(), &ret); err != nil {
		return 0, nil, errors.Wrapf(err, "parse fuzzy search result, ret: %s", string(resp.Body()))
	}

	// 仅 200/201 视为成功，其余状态码返回错误
	if ret.Status != StatusOK && ret.Status != StatusNoResult {
		return 0, nil, errors.Errorf("qcc fuzzy search error, status: %s, message: %s, orderNumber: %s",
			ret.Status, ret.Message, ret.OrderNumber)
	}

	return int64(ret.Paging.TotalRecords), ret.Result, nil
}
