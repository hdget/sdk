package param

type Env struct {
	Prefix string // 环境变量前缀
}

const (
	defaultEnvPrefix = "HD"
)

func NewEnvDefaultParam() *Env {
	return &Env{
		Prefix: defaultEnvPrefix,
	}
}
