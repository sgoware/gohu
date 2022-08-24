package nsq

import "testing"

func TestNewProducer(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "addr1",
			args:    args{addr: "101.132.79.73:9093"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewProducer(); (err != nil) != tt.wantErr {
				t.Errorf("NewProducer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
