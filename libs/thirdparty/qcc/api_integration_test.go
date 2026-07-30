package qcc

import (
	"os"
	"testing"
)

// TestFuzzySearchReal 集成测试: 调用真实企查查接口，验证实际返回结果
// 运行方式: QCC_KEY=<AppKey> QCC_SECRET_KEY=<SecretKey> go test -run TestFuzzySearchReal -v
// 未设置环境变量时自动跳过，不影响常规单元测试
func TestFuzzySearchReal(t *testing.T) {
	if testing.Short() {
		t.Skip("short 模式下跳过网络集成测试")
	}

	key := os.Getenv("QCC_KEY")
	secretKey := os.Getenv("QCC_SECRET_KEY")
	if key == "" || secretKey == "" {
		t.Skip("跳过: 未设置 QCC_KEY / QCC_SECRET_KEY 环境变量")
	}

	api := New(key, secretKey)
	total, list, err := api.FuzzySearch("阿里巴巴")
	if err != nil {
		t.Fatalf("FuzzySearch() 错误: %v", err)
	}

	t.Logf("总记录数: %d, 本页返回: %d", total, len(list))
	for i, ent := range list {
		t.Logf("[%d] %s | 法人: %s | 状态: %s | 信用代码: %s | 成立日期: %s",
			i+1, ent.Name, ent.OperName, ent.Status, ent.CreditCode, ent.StartDate)
	}

	if total == 0 {
		t.Error("total 应大于 0")
	}
	if len(list) == 0 {
		t.Fatal("结果列表不应为空")
	}
	// 每次请求最多返回 5 条记录
	if len(list) > 5 {
		t.Errorf("len(list) = %d, 超过单页上限 5", len(list))
	}
	for i, ent := range list {
		if ent.Name == "" {
			t.Errorf("[%d] 企业名称不应为空", i)
		}
		if ent.KeyNo == "" {
			t.Errorf("[%d] 企业 KeyNo 不应为空", i)
		}
	}
}
