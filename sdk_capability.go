package sdk

import "github.com/hdget/common/intf"

func Logger() intf.LoggerProvider {
	return _instance.loggerProvider
}

func Db() intf.DbProvider {
	return _instance.dbProvider
}

func Redis() intf.RedisProvider {
	return _instance.redisProvider
}

func Config() intf.ConfigProvider {
	return _instance.configProvider
}

func Mq() intf.MessageQueueProvider {
	return _instance.mqProvider
}
