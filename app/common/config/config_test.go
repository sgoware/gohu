package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/apolloconfig/agollo/v4"
	"github.com/spf13/viper"
)

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

func TestNewConfigClient(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test",
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

func TestGetConfigClient(t *testing.T) {
	tests := []struct {
		name    string
		want    *Agollo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfigClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfigClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgollo_GetViper(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "yaml",
			args:    args{namespace: "user.yaml"},
			wantErr: false,
		},
		{
			name:    "properties",
			args:    args{namespace: "NSQ"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewConfigClient()
			_, err = c.GetViper(tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Agollo.GetViper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestAgollo_UnmarshalServiceConfig(t *testing.T) {
	type fields struct {
		client agollo.Client
		vipers map[string]*viper.Viper
	}
	type args struct {
		namespace   string
		serviceType string
		serviceName string
		dst         interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Agollo{
				client: tt.fields.client,
				vipers: tt.fields.vipers,
			}
			if err := c.UnmarshalServiceConfig(tt.args.namespace, tt.args.serviceType, tt.args.serviceName, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Agollo.UnmarshalServiceConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgollo_UnmarshalKey(t *testing.T) {
	type fields struct {
		client agollo.Client
		vipers map[string]*viper.Viper
	}
	type args struct {
		namespace string
		key       string
		dst       interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Agollo{
				client: tt.fields.client,
				vipers: tt.fields.vipers,
			}
			if err := c.UnmarshalKey(tt.args.namespace, tt.args.key, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("Agollo.UnmarshalKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAgollo_GetClientDetails(t *testing.T) {
	type fields struct {
		client agollo.Client
		vipers map[string]*viper.Viper
	}
	tests := []struct {
		name            string
		fields          fields
		wantClientAuths map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Agollo{
				client: tt.fields.client,
				vipers: tt.fields.vipers,
			}
			if gotClientAuths := c.GetClientDetails(); !reflect.DeepEqual(gotClientAuths, tt.wantClientAuths) {
				t.Errorf("Agollo.GetClientDetails() = %v, want %v", gotClientAuths, tt.wantClientAuths)
			}
		})
	}
}
