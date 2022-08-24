package nsq

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
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
			_, err := GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMustGetNSQDAddr(t *testing.T) {
	if got := MustGetNSQDAddr(); got == "" {
		t.Errorf("MustGetNSQDAddr() = %v", got)
	}
}

func TestGetNSQDAddr(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNSQDAddr()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNSQDAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNSQDAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMustGetNSQLookupAddrs(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MustGetNSQLookupAddrs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MustGetNSQLookupAddrs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNSQLookupdAddrs(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNSQLookupdAddrs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNSQLookupdAddrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNSQLookupdAddrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
