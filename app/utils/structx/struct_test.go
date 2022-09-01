package structx

import (
	"reflect"
	"testing"
	"time"
)

type Payload struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Follower   int32     `json:"follower"`
	CreateTime time.Time `json:"create_time"`
}
type Payload2 struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Follower   int32     `json:"follower"`
	CreateTime time.Time `json:"create_time"`
}

func TestSyncWithNoZero(t *testing.T) {
	type args struct {
		src interface{}
		dst interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStruct interface{}
	}{
		{
			name: "1",
			args: args{
				src: Payload{
					Id:         0,
					Username:   "StellarisW",
					Follower:   3,
					CreateTime: time.Unix(1662046824, 0),
				},
				dst: &Payload2{
					Id:         2,
					Username:   "ww",
					Follower:   5,
					CreateTime: time.Time{},
				},
			},
			wantErr: false,
			wantStruct: &Payload2{
				Id:         2,
				Username:   "StellarisW",
				Follower:   3,
				CreateTime: time.Unix(1662046824, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SyncWithNoZero(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("SyncWithNoZero() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.args.dst, tt.wantStruct) {
				t.Errorf("SyncWithNoZero() dst = %v, wantStruct %v", tt.args.dst, tt.wantStruct)
			}
		})
	}
}
