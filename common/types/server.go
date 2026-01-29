package types

type HookFunction func() error

type AppServer interface {
	Start() error                               // 开始
	Stop(forced ...bool) error                  // 默认为优雅关闭, 如果forced为true, 则强制关闭
	HookPreStart(fns ...HookFunction) AppServer // 添加启动前的钩子函数
	HookPreStop(fns ...HookFunction) AppServer  // 添加停止前的钩子函数
}
