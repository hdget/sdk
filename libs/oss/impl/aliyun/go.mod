module github.com/hdget/sdk/providers/oss/aliyun

go 1.24.0

require (
	github.com/aliyun/alibabacloud-oss-go-sdk-v2 v1.4.1
	github.com/elliotchance/pie/v2 v2.9.1
	github.com/hdget/sdk/libs/oss v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
)

require (
	golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect
	golang.org/x/time v0.4.0 // indirect
)

replace github.com/hdget/sdk/libs/oss => ../..
