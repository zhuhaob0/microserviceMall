package transport

import (
	"context"
	"final-design/pb"

	endpts "final-design/user-service/endpoint"

	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	check      grpc.Handler
	create     grpc.Handler
	adminCheck grpc.Handler
}

func (s *grpcServer) Check(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error) {
	_, resp, err := s.check.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserResponse), nil
}

func (s *grpcServer) AdminCheck(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error) {
	_, resp, err := s.adminCheck.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserResponse), nil
}

func (s *grpcServer) Create(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, resp, err := s.create.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CreateUserResponse), nil
}

func NewGRPCServer(ctx context.Context, endpoints endpts.UserEndpoints, serverTracer grpc.ServerOption) pb.UserServiceServer {
	return &grpcServer{
		check: grpc.NewServer(
			endpoints.UserEndpoint,
			DecodeGRPCUserRequest,
			EncodeGRPCUserResponse,
			serverTracer,
		),
		create: grpc.NewServer(
			endpoints.CreateUserEndpoint,
			DecodeGRPCCreateUserRequest,
			EncodeGRPCCreateUserResponse,
			serverTracer,
		),
		adminCheck: grpc.NewServer(
			endpoints.AdminUserEndpoint,
			DecodeGRPCUserRequest,
			EncodeGRPCUserResponse,
			serverTracer,
		),
	}
}
