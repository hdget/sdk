# lib-dapr
dapr library

## New Invocation Module
```go
type exampleModule struct {
    module.InvocationModule
}
```

> when the suffix of module struct name is module, hd will take `example` as module name 

## Initialize module

```go
func init() {
    v := &exampleModule{}
    err := module.NewInvocationModule(v, g.App, map[string]module.InvocationFunction{
        "hello": v.helloHandler, // handler adds here
	})
    if err != nil {
        panic(err)
    }
}
```

## Add Invocation Handler
```go
func (*exampleModule) helloHandler(ctx biz.Context, data []byte) (any, error) {
	var req ExampleRequest
	if err := json.Unmarshal(data, &req); err != nil {
		return nil, err
	}

	err := xxx.New().Hello(req.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "hello world", "req", req)
	}

	return &pb.SimpleMessage{"Message": "ok"}, nil
}
```


