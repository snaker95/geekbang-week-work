package conf

import (
	"testing"
)

/*
@Time : 2021/6/17 下午10:06
@Author : snaker95
@File : boot_test.go
@Software: GoLand
*/

func TestInit(t *testing.T) {
	tests := []struct {
		name     string
		flagconf string
	}{
		{
			name:     "release",
			flagconf: "../../configs/release.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got *Config
			err := Scan(tt.flagconf, got)
			t.Logf("got = %v, err=%+v", got, err)
		})
	}
}
