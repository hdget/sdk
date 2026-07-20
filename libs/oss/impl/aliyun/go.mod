module github.com/hdget/sdk/libs/oss/impl/aliyun

go 1.25.3

require (
	github.com/aliyun/alibabacloud-oss-go-sdk-v2 v1.4.1
	github.com/elliotchance/pie/v2 v2.9.1
	github.com/hdget/sdk/libs/oss v0.0.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20220321173239-a90fa8a75705 // indirect
	golang.org/x/time v0.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/hdget/sdk/libs/oss => ../..
