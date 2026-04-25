package sdk

import (
	"github.com/hdget/sdk/common/provider"
)

func Logger() provider.Logger {
	return _instance.loggerProvider
}

func Db() provider.Database {
	return _instance.dbProvider
}

func Redis() provider.Redis {
	return _instance.redisProvider
}

func Config() provider.Config {
	return _instance.configProvider
}

func Mq() provider.MessageQueue {
	return _instance.mqProvider
}
