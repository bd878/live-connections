// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// AreaManagerClient is the client API for AreaManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AreaManagerClient interface {
	Create(ctx context.Context, in *CreateAreaRequest, opts ...grpc.CallOption) (*CreateAreaResponse, error)
	ListUsers(ctx context.Context, in *ListAreaUsersRequest, opts ...grpc.CallOption) (*ListAreaUsersResponse, error)
	HasUser(ctx context.Context, in *HasUserRequest, opts ...grpc.CallOption) (*HasUserResponse, error)
}

type areaManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewAreaManagerClient(cc grpc.ClientConnInterface) AreaManagerClient {
	return &areaManagerClient{cc}
}

func (c *areaManagerClient) Create(ctx context.Context, in *CreateAreaRequest, opts ...grpc.CallOption) (*CreateAreaResponse, error) {
	out := new(CreateAreaResponse)
	err := c.cc.Invoke(ctx, "/disk.AreaManager/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *areaManagerClient) ListUsers(ctx context.Context, in *ListAreaUsersRequest, opts ...grpc.CallOption) (*ListAreaUsersResponse, error) {
	out := new(ListAreaUsersResponse)
	err := c.cc.Invoke(ctx, "/disk.AreaManager/ListUsers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *areaManagerClient) HasUser(ctx context.Context, in *HasUserRequest, opts ...grpc.CallOption) (*HasUserResponse, error) {
	out := new(HasUserResponse)
	err := c.cc.Invoke(ctx, "/disk.AreaManager/HasUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AreaManagerServer is the server API for AreaManager service.
// All implementations must embed UnimplementedAreaManagerServer
// for forward compatibility
type AreaManagerServer interface {
	Create(context.Context, *CreateAreaRequest) (*CreateAreaResponse, error)
	ListUsers(context.Context, *ListAreaUsersRequest) (*ListAreaUsersResponse, error)
	HasUser(context.Context, *HasUserRequest) (*HasUserResponse, error)
	mustEmbedUnimplementedAreaManagerServer()
}

// UnimplementedAreaManagerServer must be embedded to have forward compatible implementations.
type UnimplementedAreaManagerServer struct {
}

func (UnimplementedAreaManagerServer) Create(context.Context, *CreateAreaRequest) (*CreateAreaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAreaManagerServer) ListUsers(context.Context, *ListAreaUsersRequest) (*ListAreaUsersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListUsers not implemented")
}
func (UnimplementedAreaManagerServer) HasUser(context.Context, *HasUserRequest) (*HasUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasUser not implemented")
}
func (UnimplementedAreaManagerServer) mustEmbedUnimplementedAreaManagerServer() {}

// UnsafeAreaManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AreaManagerServer will
// result in compilation errors.
type UnsafeAreaManagerServer interface {
	mustEmbedUnimplementedAreaManagerServer()
}

func RegisterAreaManagerServer(s *grpc.Server, srv AreaManagerServer) {
	s.RegisterService(&_AreaManager_serviceDesc, srv)
}

func _AreaManager_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAreaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AreaManagerServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.AreaManager/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AreaManagerServer).Create(ctx, req.(*CreateAreaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AreaManager_ListUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAreaUsersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AreaManagerServer).ListUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.AreaManager/ListUsers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AreaManagerServer).ListUsers(ctx, req.(*ListAreaUsersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AreaManager_HasUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HasUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AreaManagerServer).HasUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.AreaManager/HasUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AreaManagerServer).HasUser(ctx, req.(*HasUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _AreaManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "disk.AreaManager",
	HandlerType: (*AreaManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _AreaManager_Create_Handler,
		},
		{
			MethodName: "ListUsers",
			Handler:    _AreaManager_ListUsers_Handler,
		},
		{
			MethodName: "HasUser",
			Handler:    _AreaManager_HasUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "disk.proto",
}

// UserManagerClient is the client API for UserManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserManagerClient interface {
	Add(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserResponse, error)
	Read(ctx context.Context, in *ReadUserRequest, opts ...grpc.CallOption) (*ReadUserResponse, error)
}

type userManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewUserManagerClient(cc grpc.ClientConnInterface) UserManagerClient {
	return &userManagerClient{cc}
}

func (c *userManagerClient) Add(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserResponse, error) {
	out := new(AddUserResponse)
	err := c.cc.Invoke(ctx, "/disk.UserManager/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagerClient) Read(ctx context.Context, in *ReadUserRequest, opts ...grpc.CallOption) (*ReadUserResponse, error) {
	out := new(ReadUserResponse)
	err := c.cc.Invoke(ctx, "/disk.UserManager/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserManagerServer is the server API for UserManager service.
// All implementations must embed UnimplementedUserManagerServer
// for forward compatibility
type UserManagerServer interface {
	Add(context.Context, *AddUserRequest) (*AddUserResponse, error)
	Read(context.Context, *ReadUserRequest) (*ReadUserResponse, error)
	mustEmbedUnimplementedUserManagerServer()
}

// UnimplementedUserManagerServer must be embedded to have forward compatible implementations.
type UnimplementedUserManagerServer struct {
}

func (UnimplementedUserManagerServer) Add(context.Context, *AddUserRequest) (*AddUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
}
func (UnimplementedUserManagerServer) Read(context.Context, *ReadUserRequest) (*ReadUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedUserManagerServer) mustEmbedUnimplementedUserManagerServer() {}

// UnsafeUserManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserManagerServer will
// result in compilation errors.
type UnsafeUserManagerServer interface {
	mustEmbedUnimplementedUserManagerServer()
}

func RegisterUserManagerServer(s *grpc.Server, srv UserManagerServer) {
	s.RegisterService(&_UserManager_serviceDesc, srv)
}

func _UserManager_Add_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.UserManager/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).Add(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManager_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagerServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.UserManager/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagerServer).Read(ctx, req.(*ReadUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _UserManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "disk.UserManager",
	HandlerType: (*UserManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Add",
			Handler:    _UserManager_Add_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _UserManager_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "disk.proto",
}

// SquareManagerClient is the client API for SquareManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SquareManagerClient interface {
	Write(ctx context.Context, in *WriteSquareRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*Coords, error)
}

type squareManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewSquareManagerClient(cc grpc.ClientConnInterface) SquareManagerClient {
	return &squareManagerClient{cc}
}

func (c *squareManagerClient) Write(ctx context.Context, in *WriteSquareRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/disk.SquareManager/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *squareManagerClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*Coords, error) {
	out := new(Coords)
	err := c.cc.Invoke(ctx, "/disk.SquareManager/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SquareManagerServer is the server API for SquareManager service.
// All implementations must embed UnimplementedSquareManagerServer
// for forward compatibility
type SquareManagerServer interface {
	Write(context.Context, *WriteSquareRequest) (*EmptyResponse, error)
	Read(context.Context, *ReadRequest) (*Coords, error)
	mustEmbedUnimplementedSquareManagerServer()
}

// UnimplementedSquareManagerServer must be embedded to have forward compatible implementations.
type UnimplementedSquareManagerServer struct {
}

func (UnimplementedSquareManagerServer) Write(context.Context, *WriteSquareRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedSquareManagerServer) Read(context.Context, *ReadRequest) (*Coords, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedSquareManagerServer) mustEmbedUnimplementedSquareManagerServer() {}

// UnsafeSquareManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SquareManagerServer will
// result in compilation errors.
type UnsafeSquareManagerServer interface {
	mustEmbedUnimplementedSquareManagerServer()
}

func RegisterSquareManagerServer(s *grpc.Server, srv SquareManagerServer) {
	s.RegisterService(&_SquareManager_serviceDesc, srv)
}

func _SquareManager_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteSquareRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SquareManagerServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.SquareManager/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SquareManagerServer).Write(ctx, req.(*WriteSquareRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SquareManager_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SquareManagerServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.SquareManager/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SquareManagerServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SquareManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "disk.SquareManager",
	HandlerType: (*SquareManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _SquareManager_Write_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _SquareManager_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "disk.proto",
}

// TextsManagerClient is the client API for TextsManager service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TextsManagerClient interface {
	Write(ctx context.Context, in *WriteTextRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*Text, error)
}

type textsManagerClient struct {
	cc grpc.ClientConnInterface
}

func NewTextsManagerClient(cc grpc.ClientConnInterface) TextsManagerClient {
	return &textsManagerClient{cc}
}

func (c *textsManagerClient) Write(ctx context.Context, in *WriteTextRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/disk.TextsManager/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *textsManagerClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*Text, error) {
	out := new(Text)
	err := c.cc.Invoke(ctx, "/disk.TextsManager/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TextsManagerServer is the server API for TextsManager service.
// All implementations must embed UnimplementedTextsManagerServer
// for forward compatibility
type TextsManagerServer interface {
	Write(context.Context, *WriteTextRequest) (*EmptyResponse, error)
	Read(context.Context, *ReadRequest) (*Text, error)
	mustEmbedUnimplementedTextsManagerServer()
}

// UnimplementedTextsManagerServer must be embedded to have forward compatible implementations.
type UnimplementedTextsManagerServer struct {
}

func (UnimplementedTextsManagerServer) Write(context.Context, *WriteTextRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedTextsManagerServer) Read(context.Context, *ReadRequest) (*Text, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedTextsManagerServer) mustEmbedUnimplementedTextsManagerServer() {}

// UnsafeTextsManagerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TextsManagerServer will
// result in compilation errors.
type UnsafeTextsManagerServer interface {
	mustEmbedUnimplementedTextsManagerServer()
}

func RegisterTextsManagerServer(s *grpc.Server, srv TextsManagerServer) {
	s.RegisterService(&_TextsManager_serviceDesc, srv)
}

func _TextsManager_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteTextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TextsManagerServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.TextsManager/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TextsManagerServer).Write(ctx, req.(*WriteTextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TextsManager_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TextsManagerServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/disk.TextsManager/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TextsManagerServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TextsManager_serviceDesc = grpc.ServiceDesc{
	ServiceName: "disk.TextsManager",
	HandlerType: (*TextsManagerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _TextsManager_Write_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _TextsManager_Read_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "disk.proto",
}
