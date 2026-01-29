# How to generate protobuf files
```
protoc --go_out=./protobuf --go_opt=module=github.com/hdget/sdk/common/protobuf proto/*.proto
```