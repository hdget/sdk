package param

type Param struct {
	*File
	*Env
	*Remote
	*Cli
	DefaultRemoteWatcher func() // 默认的远程配置监听函数，从SDK传入
}

func GetDefaultParam() *Param {
	return &Param{
		File:   NewFileDefaultParam(),
		Env:    NewEnvDefaultParam(),
		Cli:    NewCliDefaultParam(),
		Remote: nil, // 暂时禁用remote
	}
}
