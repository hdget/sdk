package param

type Remote struct {
	Provider         string
	Endpoints        []string
	Secret           string
	RemoteConfigType string
	WatchInterval    int // 单位：秒
	WatchPath        string
	WatchCallback    func()
}

var (
	DefaultRemoteEndpoints = []string{"http://127.0.0.1:2379"}
)

func NewRemoteDefaultParam() *Remote {
	return &Remote{
		Provider:         "etcd3",
		RemoteConfigType: "json",
		WatchInterval:    10,
		WatchPath:        "",
	}
}
