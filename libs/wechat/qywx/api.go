package qywx

import (
	"github.com/hdget/sdk/common/provider"
)

type ApiCommon interface {
	GetAccessToken() (string, error) // 企业内部应用access token
	VerifyCallbackServer(token, signature, timestamp, nonce, echostr string) (string, error)
	// GetProviderAccessToken() (string, error) // 获取服务商access token
	// GetSuiteAccessToken() (string, error)    // 获取第三方应用access token
	GetCorpId() string
	GetCorpSecret() string
}

type apiImpl struct {
	corpId        string
	corpSecret    string
	redisProvider provider.Redis
}

func NewApiCommon(corpId, corpSecret string, redisProvider ...provider.Redis) ApiCommon {
	var redis provider.Redis
	if len(redisProvider) > 0 {
		redis = redisProvider[0]
	}

	return &apiImpl{
		corpId:        corpId,
		corpSecret:    corpSecret,
		redisProvider: redis,
	}
}

func (impl apiImpl) GetCorpId() string {
	return impl.corpId
}

func (impl apiImpl) GetCorpSecret() string {
	return impl.corpSecret
}
