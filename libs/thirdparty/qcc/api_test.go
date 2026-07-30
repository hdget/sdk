package qcc

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/hdget/sdk/common/protobuf"
)

// TestGenToken 验证 Token 生成规则: key + timespan + secretKey 的 32 位大写 md5
func TestGenToken(t *testing.T) {
	got := genToken("appkey", "1609459200", "secret")
	want := "D15874248F52D05C42E8D2EC06BB27BE" // 预计算的固定值
	if got != want {
		t.Errorf("genToken() = %s, want %s", got, want)
	}
	if len(got) != 32 {
		t.Errorf("genToken() 长度 = %d, want 32", len(got))
	}
}

// TestAuthHeaders 验证鉴权请求头格式
func TestAuthHeaders(t *testing.T) {
	impl := &apiImpl{key: "testkey", secretKey: "testsecret"}
	headers := impl.authHeaders()

	timespan := headers["Timespan"]
	if _, err := strconv.ParseInt(timespan, 10, 64); err != nil {
		t.Errorf("Timespan 应为 unix 秒级时间戳, got %s", timespan)
	}
	if want := genToken("testkey", timespan, "testsecret"); headers["Token"] != want {
		t.Errorf("Token = %s, want %s", headers["Token"], want)
	}
}

// TestFuzzySearch 通过 httptest 模拟企查查服务，验证请求构造与响应解析
func TestFuzzySearch(t *testing.T) {
	const (
		testKey    = "testkey"
		testSecret = "testsecret"
	)

	tests := []struct {
		name          string
		searchKey     string
		listParam     *protobuf.ListParam // nil 表示不传分页参数
		respBody      string
		wantErr       bool
		wantErrStatus string // 错误信息中应包含的状态码
		wantTotal     int64
		wantCount     int
		checkRequest  func(t *testing.T, r *http.Request)
	}{
		{
			name:      "查询成功返回企业列表，默认第 1 页",
			searchKey: "阿里巴巴",
			respBody: `{
				"Paging": {"PageSize": 5, "PageIndex": 1, "TotalRecords": 1},
				"Result": [
					{
						"KeyNo": "1234567890",
						"Name": "阿里巴巴（中国）有限公司",
						"CreditCode": "91330100799655058B",
						"StartDate": "1999-09-10",
						"OperName": "张三",
						"Status": "存续",
						"No": "330100400015575",
						"Address": "杭州市滨江区网商路699号"
					}
				],
				"Status": "200",
				"Message": "查询成功",
				"OrderNumber": "FUZZYSEARCH2021012016353715836099"
			}`,
			wantTotal: 1,
			wantCount: 1,
			checkRequest: func(t *testing.T, r *http.Request) {
				query := r.URL.Query()
				if query.Get("key") != testKey {
					t.Errorf("query key = %s, want %s", query.Get("key"), testKey)
				}
				if query.Get("searchKey") != "阿里巴巴" {
					t.Errorf("query searchKey = %s, want 阿里巴巴", query.Get("searchKey"))
				}
				// 不传分页参数时默认第 1 页
				if query.Get("pageIndex") != "1" {
					t.Errorf("query pageIndex = %s, want 1", query.Get("pageIndex"))
				}
			},
		},
		{
			name:      "指定页码",
			searchKey: "腾讯",
			listParam: &protobuf.ListParam{Page: 2, PageSize: 5},
			respBody:  `{"Paging": {"PageSize": 5, "PageIndex": 2, "TotalRecords": 12}, "Result": [], "Status": "200", "Message": "查询成功", "OrderNumber": "test"}`,
			wantTotal: 12,
			wantCount: 0,
			checkRequest: func(t *testing.T, r *http.Request) {
				if r.URL.Query().Get("pageIndex") != "2" {
					t.Errorf("query pageIndex = %s, want 2", r.URL.Query().Get("pageIndex"))
				}
			},
		},
		{
			name:      "查询无结果 201 不报错",
			searchKey: "不存在的公司",
			respBody:  `{"Paging": {"PageSize": 5, "PageIndex": 1, "TotalRecords": 0}, "Result": [], "Status": "201", "Message": "查询无结果", "OrderNumber": "test"}`,
			wantTotal: 0,
			wantCount: 0,
		},
		{
			name:          "无效请求 102 返回错误",
			searchKey:     "阿里巴巴",
			respBody:      `{"Paging": {}, "Result": [], "Status": "102", "Message": "当前 KEY 已欠费", "OrderNumber": "test"}`,
			wantErr:       true,
			wantErrStatus: "102",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRequest *http.Request
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotRequest = r
				// 验证鉴权头: Timespan 为秒级时间戳, Token 为对应签名
				timespan := r.Header.Get("Timespan")
				if _, err := strconv.ParseInt(timespan, 10, 64); err != nil {
					t.Errorf("Timespan 头不是 unix 时间戳: %s", timespan)
				}
				if want := genToken(testKey, timespan, testSecret); r.Header.Get("Token") != want {
					t.Errorf("Token 头 = %s, want %s", r.Header.Get("Token"), want)
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tt.respBody))
			}))
			defer server.Close()

			api := New(testKey, testSecret).(*apiImpl)
			api.baseURL = server.URL
			var (
				total int64
				list  []Enterprise
				err   error
			)
			if tt.listParam != nil {
				total, list, err = api.FuzzySearch(tt.searchKey, tt.listParam)
			} else {
				total, list, err = api.FuzzySearch(tt.searchKey)
			}

			if tt.wantErr {
				if err == nil {
					t.Fatal("FuzzySearch() 预期返回错误, 实际为 nil")
				}
				if !strings.Contains(err.Error(), tt.wantErrStatus) {
					t.Errorf("错误信息应包含状态码 %s: %v", tt.wantErrStatus, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("FuzzySearch() 意外错误: %v", err)
			}
			if total != tt.wantTotal {
				t.Errorf("total = %d, want %d", total, tt.wantTotal)
			}
			if len(list) != tt.wantCount {
				t.Errorf("len(list) = %d, want %d", len(list), tt.wantCount)
			}
			if tt.checkRequest != nil {
				tt.checkRequest(t, gotRequest)
			}
		})
	}
}

// TestFuzzySearchFieldMapping 验证成功响应各字段解析正确
func TestFuzzySearchFieldMapping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"Paging": {"PageSize": 5, "PageIndex": 1, "TotalRecords": 1},
			"Result": [
				{
					"KeyNo": "1234567890",
					"Name": "阿里巴巴（中国）有限公司",
					"CreditCode": "91330100799655058B",
					"StartDate": "1999-09-10",
					"OperName": "张三",
					"Status": "存续",
					"No": "330100400015575",
					"Address": "杭州市滨江区网商路699号"
				}
			],
			"Status": "200",
			"Message": "查询成功",
			"OrderNumber": "FUZZYSEARCH2021012016353715836099"
		}`))
	}))
	defer server.Close()

	api := New("testkey", "testsecret").(*apiImpl)
	api.baseURL = server.URL
	total, list, err := api.FuzzySearch("阿里巴巴")
	if err != nil {
		t.Fatalf("FuzzySearch() 意外错误: %v", err)
	}

	// 总记录数来自 Paging.TotalRecords
	if total != 1 {
		t.Errorf("total = %d, want 1", total)
	}
	if len(list) != 1 {
		t.Fatalf("len(list) = %d, want 1", len(list))
	}

	// 企业字段
	want := Enterprise{
		KeyNo:      "1234567890",
		Name:       "阿里巴巴（中国）有限公司",
		CreditCode: "91330100799655058B",
		StartDate:  "1999-09-10",
		OperName:   "张三",
		Status:     "存续",
		No:         "330100400015575",
		Address:    "杭州市滨江区网商路699号",
	}
	if list[0] != want {
		t.Errorf("Enterprise 解析错误:\ngot  %+v\nwant %+v", list[0], want)
	}
}
