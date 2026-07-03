package api

type Result struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type AccessTokenResult struct {
	Result
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
