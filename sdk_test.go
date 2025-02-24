package sdk

import (
	"fmt"
	"testing"
)

func TestSdkInstance_Initialize(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := New("ddd").Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
			var a any

			i := GetInstance()
			fmt.Println(i)

			err := Config().Unmarshal(&a, "app")
			if err != nil {
				panic(err)
			}

			Logger().Debug("xxxxxxxx")
		})
	}
}
