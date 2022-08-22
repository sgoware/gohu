package utils

import (
	"testing"
)

func TestGetServiceFullName(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "api",
			args: args{serviceName: "oauth.api"},
			want: "gohu-oauth-api",
		},
		{
			name: "rpc",
			args: args{serviceName: "oauth.rpc.tokenEnhancer"},
			want: "gohu-oauth-rpc-token-enhancer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceFullName(tt.args.serviceName); got != tt.want {
				t.Errorf("GetServiceFullName() = %v, wantNameSpace %v", got, tt.want)
			}
		})
	}
}

func TestGetNamespace(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "api",
			args: args{serviceName: "oauth.api"},
			want: "oauth.yaml",
		},
		{
			name: "rpc",
			args: args{serviceName: "oauth.rpc.tokenEnhancer"},
			want: "oauth.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNamespace(tt.args.serviceName); got != tt.want {
				t.Errorf("GetNamespace() = %v, wantNameSpace %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceType(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "api",
			args: args{serviceName: "oauth.api"},
			want: "api",
		},
		{
			name: "rpc",
			args: args{serviceName: "oauth.rpc.tokenEnhancer"},
			want: "rpc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceType(tt.args.serviceName); got != tt.want {
				t.Errorf("GetServiceType() = %v, wantNameSpace %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceSingleName(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "api",
			args: args{serviceName: "oauth.api"},
			want: "oauth",
		},
		{
			name: "rpc",
			args: args{serviceName: "oauth.rpc.tokenEnhancer"},
			want: "tokenEnhancer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceSingleName(tt.args.serviceName); got != tt.want {
				t.Errorf("GetServiceSingleName() = %v, wantNameSpace %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceDetails(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name                  string
		args                  args
		wantNameSpace         string
		wantServiceType       string
		wantServiceSingleName string
	}{
		{
			name:                  "api",
			args:                  args{serviceName: "oauth.api"},
			wantNameSpace:         "oauth.yaml",
			wantServiceType:       "api",
			wantServiceSingleName: "oauth",
		},
		{
			name:                  "rpc",
			args:                  args{serviceName: "oauth.rpc.tokenEnhancer"},
			wantNameSpace:         "oauth.yaml",
			wantServiceType:       "rpc",
			wantServiceSingleName: "tokenEnhancer",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := GetServiceDetails(tt.args.serviceName)
			if got != tt.wantNameSpace {
				t.Errorf("GetServiceDetails() got = %v, wantNameSpace %v", got, tt.wantNameSpace)
			}
			if got1 != tt.wantServiceType {
				t.Errorf("GetServiceDetails() got1 = %v, wantNameSpace %v", got1, tt.wantServiceType)
			}
			if got2 != tt.wantServiceSingleName {
				t.Errorf("GetServiceDetails() got2 = %v, wantNameSpace %v", got2, tt.wantServiceSingleName)
			}
		})
	}
}
