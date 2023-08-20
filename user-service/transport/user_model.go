package transport

import (
	"context"
	"errors"
	"final-design/pb"
	"final-design/user-service/endpoint"
)

// =============================================UserRequest编、解码========================================================

func EncodeGRPCUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.UserRequest)
	return &pb.UserRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}
func DecodeGRPCUserRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UserRequest)
	return endpoint.UserRequest{
		Username: req.Username,
		Password: req.Password,
	}, nil
}

// =============================================UserResponse编、解码=======================================================

func EncodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.UserResponse)
	if resp.Error != nil {
		return &pb.UserResponse{
			Result: resp.Result,
			Err:    "error",
		}, nil
	}

	return &pb.UserResponse{
		Result: resp.Result,
		UserId: resp.UserId,
		Err:    "",
	}, nil
}
func DecodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.UserResponse)
	return endpoint.UserResponse{
		Result: resp.Result,
		Error:  errors.New(resp.Err),
	}, nil
}

// =============================================CreateUser编、解码=======================================================

func DecodeGRPCCreateUserRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.CreateUserRequest)
	return endpoint.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
		UserId:   req.UserId,
		Age:      int(req.Age),
	}, nil
}
func EncodeGRPCCreateUserResponse(ctx context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.CreateUserResponse)
	if resp.Error != nil {
		return &pb.CreateUserResponse{
			Result: resp.Result,
			Err:    "error",
		}, nil
	}

	return &pb.CreateUserResponse{
		Result: resp.Result,
		Err:    "",
	}, nil
}
