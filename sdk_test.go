package sdk

import (
	"fmt"
	viper "github.com/hdget/provider-config-viper"
	"testing"
	"time"
)

type Conf struct {
	Global any `json:"global"`
}

func TestSdkInstance_Initialize(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config Conf

			err := New("base", WithConfigOptions(
				viper.WithEnableRemote(),
			)).UseConfig(&config).Initialize()
			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(60 * time.Second)
		})
	}
}
