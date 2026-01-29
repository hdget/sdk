package param

type Cli struct {
	Content []byte // 如果用WithConfigContent指定了配置内容，则这里不为空
}

func NewCliDefaultParam() *Cli {
	return &Cli{}
}
