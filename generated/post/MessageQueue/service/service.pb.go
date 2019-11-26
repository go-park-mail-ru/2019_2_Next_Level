// Code generated by protoc-gen-go. DO NOT EDIT.
// source: service.proto

package service

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Email struct {
	From                 string   `protobuf:"bytes,1,opt,name=From,proto3" json:"From,omitempty"`
	To                   string   `protobuf:"bytes,2,opt,name=To,proto3" json:"To,omitempty"`
	Subject              string   `protobuf:"bytes,3,opt,name=Subject,proto3" json:"Subject,omitempty"`
	Body                 string   `protobuf:"bytes,4,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Email) Reset()         { *m = Email{} }
func (m *Email) String() string { return proto.CompactTextString(m) }
func (*Email) ProtoMessage()    {}
func (*Email) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{0}
}

func (m *Email) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Email.Unmarshal(m, b)
}
func (m *Email) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Email.Marshal(b, m, deterministic)
}
func (m *Email) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Email.Merge(m, src)
}
func (m *Email) XXX_Size() int {
	return xxx_messageInfo_Email.Size(m)
}
func (m *Email) XXX_DiscardUnknown() {
	xxx_messageInfo_Email.DiscardUnknown(m)
}

var xxx_messageInfo_Email proto.InternalMessageInfo

func (m *Email) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Email) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *Email) GetSubject() string {
	if m != nil {
		return m.Subject
	}
	return ""
}

func (m *Email) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

type Empty struct {
	S                    bool     `protobuf:"varint,1,opt,name=s,proto3" json:"s,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_a0b84a42fa06f626, []int{1}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func (m *Empty) GetS() bool {
	if m != nil {
		return m.S
	}
	return false
}

func init() {
	proto.RegisterType((*Email)(nil), "service.Email")
	proto.RegisterType((*Empty)(nil), "service.Empty")
}

func init() { proto.RegisterFile("service.proto", fileDescriptor_a0b84a42fa06f626) }

var fileDescriptor_a0b84a42fa06f626 = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2d, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x22, 0xb9,
	0x58, 0x5d, 0x73, 0x13, 0x33, 0x73, 0x84, 0x84, 0xb8, 0x58, 0xdc, 0x8a, 0xf2, 0x73, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x21, 0x3e, 0x2e, 0xa6, 0x90, 0x7c, 0x09, 0x26, 0xb0,
	0x08, 0x53, 0x48, 0xbe, 0x90, 0x04, 0x17, 0x7b, 0x70, 0x69, 0x52, 0x56, 0x6a, 0x72, 0x89, 0x04,
	0x33, 0x58, 0x10, 0xc6, 0x05, 0xe9, 0x76, 0xca, 0x4f, 0xa9, 0x94, 0x60, 0x81, 0xe8, 0x06, 0xb1,
	0x95, 0x44, 0x41, 0x46, 0x17, 0x94, 0x54, 0x0a, 0xf1, 0x70, 0x31, 0x16, 0x83, 0xcd, 0xe5, 0x08,
	0x62, 0x2c, 0x36, 0xca, 0xe0, 0xe2, 0xf1, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0x0d, 0x2c, 0x4d,
	0x2d, 0x4d, 0x15, 0xd2, 0xe6, 0x62, 0x77, 0xcd, 0x2b, 0x04, 0x33, 0xf9, 0xf4, 0x60, 0xae, 0x04,
	0xbb, 0x49, 0x0a, 0x99, 0x5f, 0x50, 0x52, 0xa9, 0xc4, 0x00, 0x52, 0xec, 0x92, 0x8a, 0xa9, 0xb8,
	0xa0, 0xa4, 0x52, 0x0a, 0x4d, 0xb3, 0x12, 0x43, 0x12, 0x1b, 0xd8, 0xaf, 0xc6, 0x80, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x22, 0xe5, 0x0f, 0x94, 0xfc, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MessageQueueClient is the client API for MessageQueue service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MessageQueueClient interface {
	Enqueue(ctx context.Context, in *Email, opts ...grpc.CallOption) (*Empty, error)
	Dequeue(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Email, error)
}

type messageQueueClient struct {
	cc *grpc.ClientConn
}

func NewMessageQueueClient(cc *grpc.ClientConn) MessageQueueClient {
	return &messageQueueClient{cc}
}

func (c *messageQueueClient) Enqueue(ctx context.Context, in *Email, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/service.MessageQueue/Enqueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageQueueClient) Dequeue(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Email, error) {
	out := new(Email)
	err := c.cc.Invoke(ctx, "/service.MessageQueue/Dequeue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageQueueServer is the server API for MessageQueue service.
type MessageQueueServer interface {
	Enqueue(context.Context, *Email) (*Empty, error)
	Dequeue(context.Context, *Empty) (*Email, error)
}

// UnimplementedMessageQueueServer can be embedded to have forward compatible implementations.
type UnimplementedMessageQueueServer struct {
}

func (*UnimplementedMessageQueueServer) Enqueue(ctx context.Context, req *Email) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enqueue not implemented")
}
func (*UnimplementedMessageQueueServer) Dequeue(ctx context.Context, req *Empty) (*Email, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Dequeue not implemented")
}

func RegisterMessageQueueServer(s *grpc.Server, srv MessageQueueServer) {
	s.RegisterService(&_MessageQueue_serviceDesc, srv)
}

func _MessageQueue_Enqueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageQueueServer).Enqueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.MessageQueue/Enqueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageQueueServer).Enqueue(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageQueue_Dequeue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageQueueServer).Dequeue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.MessageQueue/Dequeue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageQueueServer).Dequeue(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _MessageQueue_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.MessageQueue",
	HandlerType: (*MessageQueueServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Enqueue",
			Handler:    _MessageQueue_Enqueue_Handler,
		},
		{
			MethodName: "Dequeue",
			Handler:    _MessageQueue_Dequeue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}