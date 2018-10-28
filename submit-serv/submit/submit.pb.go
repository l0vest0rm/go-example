// Code generated by protoc-gen-go. DO NOT EDIT.
// source: submit.proto

package submit

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SubmitRequest struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Pass                 string   `protobuf:"bytes,2,opt,name=pass,proto3" json:"pass,omitempty"`
	Url                  string   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	Title                string   `protobuf:"bytes,4,opt,name=title,proto3" json:"title,omitempty"`
	Body                 string   `protobuf:"bytes,5,opt,name=body,proto3" json:"body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubmitRequest) Reset()         { *m = SubmitRequest{} }
func (m *SubmitRequest) String() string { return proto.CompactTextString(m) }
func (*SubmitRequest) ProtoMessage()    {}
func (*SubmitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f5b4b0271dafc7e, []int{0}
}

func (m *SubmitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubmitRequest.Unmarshal(m, b)
}
func (m *SubmitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubmitRequest.Marshal(b, m, deterministic)
}
func (m *SubmitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitRequest.Merge(m, src)
}
func (m *SubmitRequest) XXX_Size() int {
	return xxx_messageInfo_SubmitRequest.Size(m)
}
func (m *SubmitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitRequest proto.InternalMessageInfo

func (m *SubmitRequest) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *SubmitRequest) GetPass() string {
	if m != nil {
		return m.Pass
	}
	return ""
}

func (m *SubmitRequest) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *SubmitRequest) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *SubmitRequest) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

type SubmitResponse struct {
	Code                 int32    `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubmitResponse) Reset()         { *m = SubmitResponse{} }
func (m *SubmitResponse) String() string { return proto.CompactTextString(m) }
func (*SubmitResponse) ProtoMessage()    {}
func (*SubmitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_0f5b4b0271dafc7e, []int{1}
}

func (m *SubmitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubmitResponse.Unmarshal(m, b)
}
func (m *SubmitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubmitResponse.Marshal(b, m, deterministic)
}
func (m *SubmitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitResponse.Merge(m, src)
}
func (m *SubmitResponse) XXX_Size() int {
	return xxx_messageInfo_SubmitResponse.Size(m)
}
func (m *SubmitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitResponse proto.InternalMessageInfo

func (m *SubmitResponse) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *SubmitResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*SubmitRequest)(nil), "submit.SubmitRequest")
	proto.RegisterType((*SubmitResponse)(nil), "submit.SubmitResponse")
}

func init() { proto.RegisterFile("submit.proto", fileDescriptor_0f5b4b0271dafc7e) }

var fileDescriptor_0f5b4b0271dafc7e = []byte{
	// 195 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x3f, 0x8b, 0xc2, 0x40,
	0x10, 0xc5, 0xc9, 0xe5, 0xcf, 0x71, 0xc3, 0xdd, 0x21, 0x83, 0xca, 0x62, 0x25, 0xa9, 0xac, 0x52,
	0x68, 0x61, 0x67, 0x6d, 0x9d, 0x7c, 0x82, 0xfc, 0x19, 0x64, 0x21, 0x71, 0x93, 0x9d, 0x8d, 0xe0,
	0xb7, 0x97, 0xcc, 0x26, 0x85, 0x76, 0xbf, 0xf7, 0x78, 0xc3, 0xdb, 0xb7, 0xf0, 0xcb, 0x63, 0xd5,
	0x69, 0x97, 0xf5, 0xd6, 0x38, 0x83, 0x89, 0x57, 0xe9, 0x00, 0x7f, 0x85, 0x50, 0x4e, 0xc3, 0x48,
	0xec, 0x70, 0x05, 0xe1, 0xa8, 0x1b, 0x15, 0xec, 0x83, 0xc3, 0x4f, 0x3e, 0x21, 0x22, 0x44, 0x7d,
	0xc9, 0xac, 0xbe, 0xc4, 0x12, 0x96, 0x94, 0x6d, 0x55, 0x38, 0xa7, 0x6c, 0x8b, 0x6b, 0x88, 0x9d,
	0x76, 0x2d, 0xa9, 0x48, 0x3c, 0x2f, 0xa6, 0xdb, 0xca, 0x34, 0x4f, 0x15, 0xfb, 0xdb, 0x89, 0xd3,
	0x0b, 0xfc, 0x2f, 0x95, 0xdc, 0x9b, 0x3b, 0x4b, 0xaa, 0x36, 0x0d, 0x49, 0x69, 0x9c, 0x0b, 0xa3,
	0x82, 0xef, 0x8e, 0x98, 0xcb, 0x1b, 0xcd, 0xc5, 0x8b, 0x3c, 0x5e, 0x97, 0x27, 0x17, 0x64, 0x1f,
	0xba, 0x26, 0x3c, 0x43, 0xe2, 0x0d, 0xdc, 0x64, 0xf3, 0xc8, 0xb7, 0x4d, 0xbb, 0xed, 0xa7, 0xed,
	0x7b, 0xab, 0x44, 0xfe, 0xe2, 0xf4, 0x0a, 0x00, 0x00, 0xff, 0xff, 0x22, 0xac, 0xf2, 0xaa, 0x1b,
	0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SubmitServiceClient is the client API for SubmitService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SubmitServiceClient interface {
	// One request followed by one response
	// The server returns the client message as-is.
	Submit(ctx context.Context, in *SubmitRequest, opts ...grpc.CallOption) (*SubmitResponse, error)
}

type submitServiceClient struct {
	cc *grpc.ClientConn
}

func NewSubmitServiceClient(cc *grpc.ClientConn) SubmitServiceClient {
	return &submitServiceClient{cc}
}

func (c *submitServiceClient) Submit(ctx context.Context, in *SubmitRequest, opts ...grpc.CallOption) (*SubmitResponse, error) {
	out := new(SubmitResponse)
	err := c.cc.Invoke(ctx, "/submit.SubmitService/Submit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubmitServiceServer is the server API for SubmitService service.
type SubmitServiceServer interface {
	// One request followed by one response
	// The server returns the client message as-is.
	Submit(context.Context, *SubmitRequest) (*SubmitResponse, error)
}

func RegisterSubmitServiceServer(s *grpc.Server, srv SubmitServiceServer) {
	s.RegisterService(&_SubmitService_serviceDesc, srv)
}

func _SubmitService_Submit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubmitServiceServer).Submit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/submit.SubmitService/Submit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubmitServiceServer).Submit(ctx, req.(*SubmitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SubmitService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "submit.SubmitService",
	HandlerType: (*SubmitServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Submit",
			Handler:    _SubmitService_Submit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "submit.proto",
}
