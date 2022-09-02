package structx

import (
	"google.golang.org/protobuf/runtime/protoimpl"
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

type UserSubjectProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Username   string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Password   string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	Nickname   string `protobuf:"bytes,4,opt,name=nickname,proto3" json:"nickname,omitempty"`
	Email      string `protobuf:"bytes,5,opt,name=email,proto3" json:"email,omitempty"`
	Phone      string `protobuf:"bytes,6,opt,name=phone,proto3" json:"phone,omitempty"`
	LastIp     string `protobuf:"bytes,7,opt,name=last_ip,json=lastIp,proto3" json:"last_ip,omitempty"`
	Vip        int32  `protobuf:"varint,8,opt,name=vip,proto3" json:"vip,omitempty"`
	Follower   int32  `protobuf:"varint,9,opt,name=follower,proto3" json:"follower,omitempty"`
	State      int32  `protobuf:"varint,10,opt,name=state,proto3" json:"state,omitempty"`
	CreateTime string `protobuf:"bytes,11,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	UpdateTime string `protobuf:"bytes,12,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

type UserSubjectModel struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`             // 主键
	Username   string    `gorm:"column:username;not null" json:"username"`                      // 用户名 (登陆用)
	Password   string    `gorm:"column:password;not null" json:"password"`                      // 密码
	Nickname   string    `gorm:"column:nickname;not null" json:"nickname"`                      // 昵称
	Email      string    `gorm:"column:email" json:"email"`                                     // 邮箱
	Phone      string    `gorm:"column:phone" json:"phone"`                                     // 手机号
	LastIP     string    `gorm:"column:last_ip;not null" json:"last_ip"`                        // 最近登录 ip 地址
	Vip        int32     `gorm:"column:vip;not null" json:"vip"`                                // vip 等级
	Follower   int32     `gorm:"column:follower;not null" json:"follower"`                      // 被关注数
	State      int32     `gorm:"column:state;not null" json:"state"`                            // 状态 (0-正常 1-冻结 2-封禁)
	CreateTime time.Time `gorm:"autoCreateTime;column:create_time;not null" json:"create_time"` // 创建时间
	UpdateTime time.Time `gorm:"autoUpdateTime;column:update_time;not null" json:"update_time"` // 修改时间
}

type UserSubjectPayload struct {
	Id         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Nickname   string    `json:"nickname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	LastIp     string    `json:"last_ip"`
	Vip        int32     `json:"vip"`
	Follower   int32     `json:"follower"`
	State      int32     `json:"state"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
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
		{
			name: "2",
			args: args{
				src: UserSubjectModel{
					ID:         10001,
					Username:   "StellarisW",
					Password:   "ad13ce12",
					Nickname:   "24cf24c",
					Email:      "awdawd",
					Phone:      "awdwa",
					LastIP:     "dawdw",
					Vip:        0,
					Follower:   0,
					State:      0,
					CreateTime: time.Unix(1662046824, 0),
					UpdateTime: time.Unix(1662046824, 0),
				},
				dst: &UserSubjectProto{
					state:         protoimpl.MessageState{},
					sizeCache:     0,
					unknownFields: nil,
					Id:            0,
					Username:      "",
					Password:      "",
					Nickname:      "",
					Email:         "",
					Phone:         "",
					LastIp:        "",
					Vip:           0,
					Follower:      0,
					State:         0,
					CreateTime:    "",
					UpdateTime:    "",
				},
			},
			wantErr: false,
			wantStruct: &UserSubjectProto{
				Id:         10001,
				Username:   "StellarisW",
				Password:   "ad13ce12",
				Nickname:   "24cf24c",
				Email:      "awdawd",
				Phone:      "awdwa",
				LastIp:     "dawdw",
				Vip:        0,
				Follower:   0,
				State:      0,
				CreateTime: time.Unix(1662046824, 0).String(),
				UpdateTime: time.Unix(1662046824, 0).String(),
			},
		},
		{
			name: "3",
			args: args{
				src: UserSubjectPayload{
					Id:         10006,
					Username:   "StellarisW",
					Password:   "asd",
					Nickname:   "qwd",
					Email:      "",
					Phone:      "",
					LastIp:     "",
					Vip:        0,
					Follower:   0,
					State:      0,
					CreateTime: time.Unix(1662046824, 0),
					UpdateTime: time.Unix(1662046824, 0),
				},
				dst: &UserSubjectProto{},
			},
			wantErr: false,
			wantStruct: &UserSubjectProto{
				state:         protoimpl.MessageState{},
				sizeCache:     0,
				unknownFields: nil,
				Id:            10006,
				Username:      "StellarisW",
				Password:      "asd",
				Nickname:      "qwd",
				Email:         "",
				Phone:         "",
				LastIp:        "",
				Vip:           0,
				Follower:      0,
				State:         0,
				CreateTime:    time.Unix(1662046824, 0).String(),
				UpdateTime:    time.Unix(1662046824, 0).String(),
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
