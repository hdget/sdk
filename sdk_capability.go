package sdk

import (
	"github.com/hdget/common/types"
)

func Logger() types.LoggerProvider {
	return _instance.loggerProvider
}

func Db() types.DbProvider {
	return _instance.dbProvider
}

func Redis() types.RedisProvider {
	return _instance.redisProvider
}

func Config() types.ConfigProvider {
	return _instance.configProvider
}

func Mq() types.MessageQueueProvider {
	return _instance.mqProvider
}

func Oss() types.OssProvider {
	return _instance.ossProvider
}
