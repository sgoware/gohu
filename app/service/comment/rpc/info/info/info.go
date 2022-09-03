// Code generated by goctl. DO NOT EDIT!
// Source: info.proto

package info

import (
	"context"

	"main/app/service/comment/rpc/info/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CommentContent                 = pb.CommentContent
	CommentIndex                   = pb.CommentIndex
	CommentSubject                 = pb.CommentSubject
	GetCommentInfoReq              = pb.GetCommentInfoReq
	GetCommentInfoRes              = pb.GetCommentInfoRes
	GetCommentInfoRes_Data         = pb.GetCommentInfoRes_Data
	GetCommentSubjectIndexReq      = pb.GetCommentSubjectIndexReq
	GetCommentSubjectIndexRes      = pb.GetCommentSubjectIndexRes
	GetCommentSubjectIndexRes_Data = pb.GetCommentSubjectIndexRes_Data
	GetCommentSubjectReq           = pb.GetCommentSubjectReq
	GetCommentSubjectRes           = pb.GetCommentSubjectRes
	GetCommentSubjectRes_Data      = pb.GetCommentSubjectRes_Data

	Info interface {
		GetCommentSubject(ctx context.Context, in *GetCommentSubjectReq, opts ...grpc.CallOption) (*GetCommentSubjectRes, error)
		GetCommentInfo(ctx context.Context, in *GetCommentInfoReq, opts ...grpc.CallOption) (*GetCommentInfoRes, error)
		GetCommentSubjectIndex(ctx context.Context, in *GetCommentSubjectIndexReq, opts ...grpc.CallOption) (*GetCommentSubjectIndexRes, error)
	}

	defaultInfo struct {
		cli zrpc.Client
	}
)

func NewInfo(cli zrpc.Client) Info {
	return &defaultInfo{
		cli: cli,
	}
}

func (m *defaultInfo) GetCommentSubject(ctx context.Context, in *GetCommentSubjectReq, opts ...grpc.CallOption) (*GetCommentSubjectRes, error) {
	client := pb.NewInfoClient(m.cli.Conn())
	return client.GetCommentSubject(ctx, in, opts...)
}

func (m *defaultInfo) GetCommentInfo(ctx context.Context, in *GetCommentInfoReq, opts ...grpc.CallOption) (*GetCommentInfoRes, error) {
	client := pb.NewInfoClient(m.cli.Conn())
	return client.GetCommentInfo(ctx, in, opts...)
}

func (m *defaultInfo) GetCommentSubjectIndex(ctx context.Context, in *GetCommentSubjectIndexReq, opts ...grpc.CallOption) (*GetCommentSubjectIndexRes, error) {
	client := pb.NewInfoClient(m.cli.Conn())
	return client.GetCommentSubjectIndex(ctx, in, opts...)
}
