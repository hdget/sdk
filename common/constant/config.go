package constant

const (
	ConfigPathRemote         = "/config/%s" // 全局配置 /config/<environment>
	ConfigKeySdk             = "sdk"
	ConfigKeyRemoteEndpoints = ConfigKeySdk + ".%s.endpoints" // 远程配置的configKey, e,g: sdk.etcd3.endpoints
)
