# How to generate protobuf files
```
protoc -I ./proto --go_out=./protobuf --go_opt=module=github.com/hdget/sdk/common/protobuf proto/*.proto
```