package ip

import "testing"

func TestGetIPLocFromApi(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		wantLoc string
		wantErr bool
	}{
		{
			name:    "国内直辖市",
			args:    args{ip: "125.86.164.146"},
			wantLoc: "重庆重庆市",
			wantErr: false,
		},
		{
			name:    "国内地级市1",
			args:    args{ip: "222.209.14.199"},
			wantLoc: "四川省成都市",
			wantErr: false,
		},
		{
			name:    "国内地级市2",
			args:    args{ip: "218.72.111.105"},
			wantLoc: "浙江省杭州市",
			wantErr: false,
		},

		{
			name:    "国外",
			args:    args{ip: "188.166.171.114"},
			wantLoc: "英国",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLoc, err := GetIpLocFromApi(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIPLoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLoc != tt.wantLoc {
				t.Errorf("GetIPLoc() = %v, want %v", gotLoc, tt.wantLoc)
			}
		})
	}
}
