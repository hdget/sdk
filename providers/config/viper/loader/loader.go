package loader

type Loader interface {
	Load() error
}

//// Load 从各个配置源获取配置数据, 并加载到configVar中， 同名变量配置高的覆盖低的
//// - remote: kvstore配置（低）
//// - File: 文件配置(中）
//// - env: 环境变量配置(高)
//func (p *viperConfigProvider) Load() error {
//	// 如果指定了配置内容，则合并
//	if p.config.Content != nil {
//		_ = p.localViper.MergeConfig(bytes.NewReader(p.config.Content))
//	}
//
//	// 如果环境变量为空，则加载最小基本配置
//	if p.env == "" {
//		return p.loadMinimal()
//	}
//
//	// 尝试从环境变量中获取配置信息
//	p.loadFromEnv()
//
//	// 尝试从配置文件中获取配置信息
//	return p.loadFromFile()
//}
