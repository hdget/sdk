package viper

import (
	"github.com/hdget/providers/config/viper/param"
)

type Option func(*param.Param)

func WithConfigFile(configFile string) Option {
	return func(param *param.Param) {
		param.File.File = configFile
	}
}

func WithConfigContent(configContent []byte) Option {
	return func(param *param.Param) {
		param.Cli.Content = configContent
	}
}

func WithDefaultRemote() Option {
	return func(p *param.Param) {
		p.Remote = param.NewRemoteDefaultParam()
		if p.Remote.WatchCallback == nil {
			p.Remote.WatchCallback = p.DefaultRemoteWatcher
		}
	}
}

func WithRemote(remoteParam *param.Remote) Option {
	return func(p *param.Param) {
		p.Remote = remoteParam
		if p.Remote.WatchCallback == nil {
			p.Remote.WatchCallback = p.DefaultRemoteWatcher
		}
	}
}

func WithRemoteWatcher(watcher func()) Option {
	return func(p *param.Param) {
		p.DefaultRemoteWatcher = watcher
	}
}
