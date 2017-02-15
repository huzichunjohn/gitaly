// Code generated by protoc-gen-go.
// source: gitaly.proto
// DO NOT EDIT!

/*
Package gitaly is a generated protocol buffer package.

It is generated from these files:
	gitaly.proto

It has these top-level messages:
	InfoRefsRequest
	InfoRefsResponse
	Repository
	PostReceiveRequest
	PostReceiveResponse
*/
package gitaly

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type InfoRefsRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
}

func (m *InfoRefsRequest) Reset()                    { *m = InfoRefsRequest{} }
func (m *InfoRefsRequest) String() string            { return proto.CompactTextString(m) }
func (*InfoRefsRequest) ProtoMessage()               {}
func (*InfoRefsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *InfoRefsRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

type InfoRefsResponse struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *InfoRefsResponse) Reset()                    { *m = InfoRefsResponse{} }
func (m *InfoRefsResponse) String() string            { return proto.CompactTextString(m) }
func (*InfoRefsResponse) ProtoMessage()               {}
func (*InfoRefsResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *InfoRefsResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Repository struct {
	Path string `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
}

func (m *Repository) Reset()                    { *m = Repository{} }
func (m *Repository) String() string            { return proto.CompactTextString(m) }
func (*Repository) ProtoMessage()               {}
func (*Repository) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Repository) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

type PostReceiveRequest struct {
	Repository *Repository `protobuf:"bytes,1,opt,name=repository" json:"repository,omitempty"`
}

func (m *PostReceiveRequest) Reset()                    { *m = PostReceiveRequest{} }
func (m *PostReceiveRequest) String() string            { return proto.CompactTextString(m) }
func (*PostReceiveRequest) ProtoMessage()               {}
func (*PostReceiveRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *PostReceiveRequest) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

type PostReceiveResponse struct {
}

func (m *PostReceiveResponse) Reset()                    { *m = PostReceiveResponse{} }
func (m *PostReceiveResponse) String() string            { return proto.CompactTextString(m) }
func (*PostReceiveResponse) ProtoMessage()               {}
func (*PostReceiveResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto.RegisterType((*InfoRefsRequest)(nil), "gitaly.InfoRefsRequest")
	proto.RegisterType((*InfoRefsResponse)(nil), "gitaly.InfoRefsResponse")
	proto.RegisterType((*Repository)(nil), "gitaly.Repository")
	proto.RegisterType((*PostReceiveRequest)(nil), "gitaly.PostReceiveRequest")
	proto.RegisterType((*PostReceiveResponse)(nil), "gitaly.PostReceiveResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SmartHTTP service

type SmartHTTPClient interface {
	// The response body for GET /info/refs?service=git-upload-pack
	InfoRefsUploadPack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTP_InfoRefsUploadPackClient, error)
	// The response body for GET /info/refs?service=git-receive-pack
	InfoRefsReceivePack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTP_InfoRefsReceivePackClient, error)
}

type smartHTTPClient struct {
	cc *grpc.ClientConn
}

func NewSmartHTTPClient(cc *grpc.ClientConn) SmartHTTPClient {
	return &smartHTTPClient{cc}
}

func (c *smartHTTPClient) InfoRefsUploadPack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTP_InfoRefsUploadPackClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SmartHTTP_serviceDesc.Streams[0], c.cc, "/gitaly.SmartHTTP/InfoRefsUploadPack", opts...)
	if err != nil {
		return nil, err
	}
	x := &smartHTTPInfoRefsUploadPackClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SmartHTTP_InfoRefsUploadPackClient interface {
	Recv() (*InfoRefsResponse, error)
	grpc.ClientStream
}

type smartHTTPInfoRefsUploadPackClient struct {
	grpc.ClientStream
}

func (x *smartHTTPInfoRefsUploadPackClient) Recv() (*InfoRefsResponse, error) {
	m := new(InfoRefsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *smartHTTPClient) InfoRefsReceivePack(ctx context.Context, in *InfoRefsRequest, opts ...grpc.CallOption) (SmartHTTP_InfoRefsReceivePackClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_SmartHTTP_serviceDesc.Streams[1], c.cc, "/gitaly.SmartHTTP/InfoRefsReceivePack", opts...)
	if err != nil {
		return nil, err
	}
	x := &smartHTTPInfoRefsReceivePackClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SmartHTTP_InfoRefsReceivePackClient interface {
	Recv() (*InfoRefsResponse, error)
	grpc.ClientStream
}

type smartHTTPInfoRefsReceivePackClient struct {
	grpc.ClientStream
}

func (x *smartHTTPInfoRefsReceivePackClient) Recv() (*InfoRefsResponse, error) {
	m := new(InfoRefsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for SmartHTTP service

type SmartHTTPServer interface {
	// The response body for GET /info/refs?service=git-upload-pack
	InfoRefsUploadPack(*InfoRefsRequest, SmartHTTP_InfoRefsUploadPackServer) error
	// The response body for GET /info/refs?service=git-receive-pack
	InfoRefsReceivePack(*InfoRefsRequest, SmartHTTP_InfoRefsReceivePackServer) error
}

func RegisterSmartHTTPServer(s *grpc.Server, srv SmartHTTPServer) {
	s.RegisterService(&_SmartHTTP_serviceDesc, srv)
}

func _SmartHTTP_InfoRefsUploadPack_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InfoRefsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SmartHTTPServer).InfoRefsUploadPack(m, &smartHTTPInfoRefsUploadPackServer{stream})
}

type SmartHTTP_InfoRefsUploadPackServer interface {
	Send(*InfoRefsResponse) error
	grpc.ServerStream
}

type smartHTTPInfoRefsUploadPackServer struct {
	grpc.ServerStream
}

func (x *smartHTTPInfoRefsUploadPackServer) Send(m *InfoRefsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _SmartHTTP_InfoRefsReceivePack_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(InfoRefsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SmartHTTPServer).InfoRefsReceivePack(m, &smartHTTPInfoRefsReceivePackServer{stream})
}

type SmartHTTP_InfoRefsReceivePackServer interface {
	Send(*InfoRefsResponse) error
	grpc.ServerStream
}

type smartHTTPInfoRefsReceivePackServer struct {
	grpc.ServerStream
}

func (x *smartHTTPInfoRefsReceivePackServer) Send(m *InfoRefsResponse) error {
	return x.ServerStream.SendMsg(m)
}

var _SmartHTTP_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gitaly.SmartHTTP",
	HandlerType: (*SmartHTTPServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "InfoRefsUploadPack",
			Handler:       _SmartHTTP_InfoRefsUploadPack_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "InfoRefsReceivePack",
			Handler:       _SmartHTTP_InfoRefsReceivePack_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "gitaly.proto",
}

// Client API for Notifications service

type NotificationsClient interface {
	PostReceive(ctx context.Context, in *PostReceiveRequest, opts ...grpc.CallOption) (*PostReceiveResponse, error)
}

type notificationsClient struct {
	cc *grpc.ClientConn
}

func NewNotificationsClient(cc *grpc.ClientConn) NotificationsClient {
	return &notificationsClient{cc}
}

func (c *notificationsClient) PostReceive(ctx context.Context, in *PostReceiveRequest, opts ...grpc.CallOption) (*PostReceiveResponse, error) {
	out := new(PostReceiveResponse)
	err := grpc.Invoke(ctx, "/gitaly.Notifications/PostReceive", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Notifications service

type NotificationsServer interface {
	PostReceive(context.Context, *PostReceiveRequest) (*PostReceiveResponse, error)
}

func RegisterNotificationsServer(s *grpc.Server, srv NotificationsServer) {
	s.RegisterService(&_Notifications_serviceDesc, srv)
}

func _Notifications_PostReceive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostReceiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationsServer).PostReceive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gitaly.Notifications/PostReceive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationsServer).PostReceive(ctx, req.(*PostReceiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Notifications_serviceDesc = grpc.ServiceDesc{
	ServiceName: "gitaly.Notifications",
	HandlerType: (*NotificationsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PostReceive",
			Handler:    _Notifications_PostReceive_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gitaly.proto",
}

func init() { proto.RegisterFile("gitaly.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 256 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x92, 0xc1, 0x4a, 0xc3, 0x40,
	0x10, 0x86, 0x0d, 0x48, 0xa1, 0xd3, 0x8a, 0x32, 0x45, 0x2c, 0xf1, 0x52, 0xf6, 0x20, 0x9e, 0x8a,
	0xc4, 0x67, 0x10, 0x22, 0x8a, 0x84, 0xb5, 0x1e, 0x3c, 0xae, 0xe9, 0x44, 0x17, 0x6b, 0x66, 0xcd,
	0x8e, 0x42, 0x5f, 0xc8, 0xe7, 0x14, 0x37, 0xae, 0xa9, 0xd6, 0x5b, 0x6e, 0xb3, 0xfb, 0xff, 0xf3,
	0xed, 0xbf, 0xc3, 0xc0, 0xf8, 0xd1, 0x8a, 0x59, 0xad, 0xe7, 0xae, 0x61, 0x61, 0x1c, 0xb4, 0x27,
	0x75, 0x01, 0xfb, 0x97, 0x75, 0xc5, 0x9a, 0x2a, 0xaf, 0xe9, 0xf5, 0x8d, 0xbc, 0x60, 0x06, 0xd0,
	0x90, 0x63, 0x6f, 0x85, 0x9b, 0xf5, 0x34, 0x99, 0x25, 0xa7, 0xa3, 0x0c, 0xe7, 0xdf, 0xdd, 0xfa,
	0x47, 0xd1, 0x1b, 0x2e, 0x75, 0x02, 0x07, 0x1d, 0xc6, 0x3b, 0xae, 0x3d, 0x21, 0xc2, 0xee, 0xd2,
	0x88, 0x09, 0x84, 0xb1, 0x0e, 0xb5, 0x9a, 0x01, 0x74, 0x84, 0x2f, 0x87, 0x33, 0xf2, 0x14, 0x1c,
	0x43, 0x1d, 0x6a, 0x95, 0x03, 0x16, 0xec, 0x45, 0x53, 0x49, 0xf6, 0x9d, 0xfa, 0x64, 0x3a, 0x84,
	0xc9, 0x2f, 0x52, 0x1b, 0x2b, 0xfb, 0x48, 0x60, 0x78, 0xfb, 0x62, 0x1a, 0xc9, 0x17, 0x8b, 0x02,
	0xaf, 0x00, 0x63, 0xf0, 0x3b, 0xb7, 0x62, 0xb3, 0x2c, 0x4c, 0xf9, 0x8c, 0x47, 0x11, 0xfd, 0x67,
	0x36, 0xe9, 0x74, 0x5b, 0x68, 0xb1, 0x6a, 0xe7, 0x2c, 0xc1, 0x6b, 0x98, 0x74, 0xf7, 0xe1, 0xd5,
	0x1e, 0xb4, 0xec, 0x1e, 0xf6, 0x6e, 0x58, 0x6c, 0x65, 0x4b, 0x23, 0x96, 0x6b, 0x8f, 0x39, 0x8c,
	0x36, 0x3e, 0x84, 0x69, 0xec, 0xde, 0x9e, 0x57, 0x7a, 0xfc, 0xaf, 0x16, 0xe1, 0x0f, 0x83, 0xb0,
	0x04, 0xe7, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x86, 0x0b, 0x65, 0xc1, 0x14, 0x02, 0x00, 0x00,
}