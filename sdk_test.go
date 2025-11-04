package sdk

import (
	"fmt"
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

			err := New("core", WithDefaultRemote()).UseConfig(&config).Initialize()
			if err != nil {
				fmt.Println(err)
			}

			time.Sleep(60 * time.Second)
		})
	}
}
