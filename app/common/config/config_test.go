package config

import (
	"fmt"
	"testing"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfigClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfigClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAgollo_GetMysqlDsn(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test for user",
			args: args{namespace: "user.yaml"},
		},
	}
	client, _ := NewConfigClient()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.GetMysqlDsn(tt.args.namespace)
			fmt.Println(got)
			if got == "" {
				t.Errorf("get mysql dsn failed")
			}
		})
	}
}

func TestAgollo_NewRedisOptions(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test for user",
			args: args{namespace: "user.yaml"},
		},
	}
	client, _ := NewConfigClient()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := client.NewRedisOptions(tt.args.namespace)
			fmt.Println(got)
			if got == nil {
				t.Errorf("get redisOptions failed")
			}
		})
	}
}
