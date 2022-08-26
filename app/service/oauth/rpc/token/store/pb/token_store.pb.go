// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: token_store.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OAuth2Token struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RefreshToken *OAuth2Token `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	TokenType    string       `protobuf:"bytes,2,opt,name=token_type,json=tokenType,proto3" json:"token_type,omitempty"`
	TokenValue   string       `protobuf:"bytes,3,opt,name=token_value,json=tokenValue,proto3" json:"token_value,omitempty"`
	ExpiresAt    int64        `protobuf:"varint,4,opt,name=expires_at,json=expiresAt,proto3" json:"expires_at,omitempty"`
}

func (x *OAuth2Token) Reset() {
	*x = OAuth2Token{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OAuth2Token) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OAuth2Token) ProtoMessage() {}

func (x *OAuth2Token) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OAuth2Token.ProtoReflect.Descriptor instead.
func (*OAuth2Token) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{0}
}

func (x *OAuth2Token) GetRefreshToken() *OAuth2Token {
	if x != nil {
		return x.RefreshToken
	}
	return nil
}

func (x *OAuth2Token) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *OAuth2Token) GetTokenValue() string {
	if x != nil {
		return x.TokenValue
	}
	return ""
}

func (x *OAuth2Token) GetExpiresAt() int64 {
	if x != nil {
		return x.ExpiresAt
	}
	return 0
}

type StoreTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      int64        `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	AccessToken *OAuth2Token `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *StoreTokenReq) Reset() {
	*x = StoreTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreTokenReq) ProtoMessage() {}

func (x *StoreTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreTokenReq.ProtoReflect.Descriptor instead.
func (*StoreTokenReq) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{1}
}

func (x *StoreTokenReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *StoreTokenReq) GetAccessToken() *OAuth2Token {
	if x != nil {
		return x.AccessToken
	}
	return nil
}

type StoreTokenRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok  bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Msg string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *StoreTokenRes) Reset() {
	*x = StoreTokenRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreTokenRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreTokenRes) ProtoMessage() {}

func (x *StoreTokenRes) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreTokenRes.ProtoReflect.Descriptor instead.
func (*StoreTokenRes) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{2}
}

func (x *StoreTokenRes) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *StoreTokenRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type CheckTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *CheckTokenReq) Reset() {
	*x = CheckTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckTokenReq) ProtoMessage() {}

func (x *CheckTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckTokenReq.ProtoReflect.Descriptor instead.
func (*CheckTokenReq) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{3}
}

func (x *CheckTokenReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type CheckTokenRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok      bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Msg     string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	IsExist bool   `protobuf:"varint,3,opt,name=is_exist,json=isExist,proto3" json:"is_exist,omitempty"`
}

func (x *CheckTokenRes) Reset() {
	*x = CheckTokenRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckTokenRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckTokenRes) ProtoMessage() {}

func (x *CheckTokenRes) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckTokenRes.ProtoReflect.Descriptor instead.
func (*CheckTokenRes) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{4}
}

func (x *CheckTokenRes) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *CheckTokenRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *CheckTokenRes) GetIsExist() bool {
	if x != nil {
		return x.IsExist
	}
	return false
}

type GetTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *GetTokenReq) Reset() {
	*x = GetTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTokenReq) ProtoMessage() {}

func (x *GetTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTokenReq.ProtoReflect.Descriptor instead.
func (*GetTokenReq) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{5}
}

func (x *GetTokenReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type GetTokenRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok   bool              `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Msg  string            `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *GetTokenRes_Data `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *GetTokenRes) Reset() {
	*x = GetTokenRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTokenRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTokenRes) ProtoMessage() {}

func (x *GetTokenRes) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTokenRes.ProtoReflect.Descriptor instead.
func (*GetTokenRes) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{6}
}

func (x *GetTokenRes) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *GetTokenRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *GetTokenRes) GetData() *GetTokenRes_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type RemoveTokenReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *RemoveTokenReq) Reset() {
	*x = RemoveTokenReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveTokenReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveTokenReq) ProtoMessage() {}

func (x *RemoveTokenReq) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveTokenReq.ProtoReflect.Descriptor instead.
func (*RemoveTokenReq) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{7}
}

func (x *RemoveTokenReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type RemoveTokenRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok  bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Msg string `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *RemoveTokenRes) Reset() {
	*x = RemoveTokenRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveTokenRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveTokenRes) ProtoMessage() {}

func (x *RemoveTokenRes) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveTokenRes.ProtoReflect.Descriptor instead.
func (*RemoveTokenRes) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{8}
}

func (x *RemoveTokenRes) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *RemoveTokenRes) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type GetTokenRes_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OauthToken *OAuth2Token `protobuf:"bytes,1,opt,name=oauth_token,json=oauthToken,proto3" json:"oauth_token,omitempty"`
}

func (x *GetTokenRes_Data) Reset() {
	*x = GetTokenRes_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_token_store_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTokenRes_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTokenRes_Data) ProtoMessage() {}

func (x *GetTokenRes_Data) ProtoReflect() protoreflect.Message {
	mi := &file_token_store_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTokenRes_Data.ProtoReflect.Descriptor instead.
func (*GetTokenRes_Data) Descriptor() ([]byte, []int) {
	return file_token_store_proto_rawDescGZIP(), []int{6, 0}
}

func (x *GetTokenRes_Data) GetOauthToken() *OAuth2Token {
	if x != nil {
		return x.OauthToken
	}
	return nil
}

var File_token_store_proto protoreflect.FileDescriptor

var file_token_store_proto_rawDesc = []byte{
	0x0a, 0x11, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0xa5, 0x01, 0x0a, 0x0b, 0x4f,
	0x41, 0x75, 0x74, 0x68, 0x32, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x37, 0x0a, 0x0d, 0x72, 0x65,
	0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x12, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x32,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73, 0x5f, 0x61,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x73,
	0x41, 0x74, 0x22, 0x5f, 0x0a, 0x0d, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x35, 0x0a, 0x0c,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4f, 0x41, 0x75, 0x74, 0x68,
	0x32, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x22, 0x31, 0x0a, 0x0d, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x02, 0x6f, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x28, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x4c, 0x0a, 0x0d, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x73, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f,
	0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x73, 0x67, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x5f, 0x65, 0x78, 0x69, 0x73, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x69, 0x73, 0x45, 0x78, 0x69, 0x73, 0x74, 0x22, 0x26,
	0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x99, 0x01, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47,
	0x65, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x3b, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x33, 0x0a,
	0x0b, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x4f, 0x41, 0x75, 0x74, 0x68,
	0x32, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x0a, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x29, 0x0a, 0x0e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x32, 0x0a,
	0x0e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x12,
	0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12,
	0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73,
	0x67, 0x32, 0xf1, 0x01, 0x0a, 0x0a, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x12, 0x38, 0x0a, 0x0a, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x53, 0x74, 0x6f,
	0x72, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x0a, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x14,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x73, 0x12, 0x32, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x12, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x47, 0x65, 0x74,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x0b, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x15, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x1a, 0x15,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x52, 0x65, 0x73, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_token_store_proto_rawDescOnce sync.Once
	file_token_store_proto_rawDescData = file_token_store_proto_rawDesc
)

func file_token_store_proto_rawDescGZIP() []byte {
	file_token_store_proto_rawDescOnce.Do(func() {
		file_token_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_token_store_proto_rawDescData)
	})
	return file_token_store_proto_rawDescData
}

var file_token_store_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_token_store_proto_goTypes = []interface{}{
	(*OAuth2Token)(nil),      // 0: store.OAuth2Token
	(*StoreTokenReq)(nil),    // 1: store.StoreTokenReq
	(*StoreTokenRes)(nil),    // 2: store.StoreTokenRes
	(*CheckTokenReq)(nil),    // 3: store.CheckTokenReq
	(*CheckTokenRes)(nil),    // 4: store.CheckTokenRes
	(*GetTokenReq)(nil),      // 5: store.GetTokenReq
	(*GetTokenRes)(nil),      // 6: store.GetTokenRes
	(*RemoveTokenReq)(nil),   // 7: store.RemoveTokenReq
	(*RemoveTokenRes)(nil),   // 8: store.RemoveTokenRes
	(*GetTokenRes_Data)(nil), // 9: store.GetTokenRes.Data
}
var file_token_store_proto_depIdxs = []int32{
	0, // 0: store.OAuth2Token.refresh_token:type_name -> store.OAuth2Token
	0, // 1: store.StoreTokenReq.access_token:type_name -> store.OAuth2Token
	9, // 2: store.GetTokenRes.data:type_name -> store.GetTokenRes.Data
	0, // 3: store.GetTokenRes.Data.oauth_token:type_name -> store.OAuth2Token
	1, // 4: store.TokenStore.StoreToken:input_type -> store.StoreTokenReq
	3, // 5: store.TokenStore.CheckToken:input_type -> store.CheckTokenReq
	5, // 6: store.TokenStore.GetToken:input_type -> store.GetTokenReq
	7, // 7: store.TokenStore.RemoveToken:input_type -> store.RemoveTokenReq
	2, // 8: store.TokenStore.StoreToken:output_type -> store.StoreTokenRes
	4, // 9: store.TokenStore.CheckToken:output_type -> store.CheckTokenRes
	6, // 10: store.TokenStore.GetToken:output_type -> store.GetTokenRes
	8, // 11: store.TokenStore.RemoveToken:output_type -> store.RemoveTokenRes
	8, // [8:12] is the sub-list for method output_type
	4, // [4:8] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_token_store_proto_init() }
func file_token_store_proto_init() {
	if File_token_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_token_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OAuth2Token); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreTokenReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreTokenRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckTokenReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckTokenRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTokenReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTokenRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveTokenReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveTokenRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_token_store_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTokenRes_Data); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_token_store_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_token_store_proto_goTypes,
		DependencyIndexes: file_token_store_proto_depIdxs,
		MessageInfos:      file_token_store_proto_msgTypes,
	}.Build()
	File_token_store_proto = out.File
	file_token_store_proto_rawDesc = nil
	file_token_store_proto_goTypes = nil
	file_token_store_proto_depIdxs = nil
}
