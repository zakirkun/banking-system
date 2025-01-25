package auth

import (
	"context"

	pb "github.com/zakirkun/banking-microservices/api/proto"
)

type GRPCServer struct {
	pb.UnimplementedAuthServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	resp, err := s.svc.Register(req.Username, req.Password, req.Email, req.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{
		Token:     resp.Token,
		UserId:    resp.UserID,
		ExpiresAt: resp.ExpiresAt,
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	resp, err := s.svc.Login(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{
		Token:     resp.Token,
		UserId:    resp.UserID,
		ExpiresAt: resp.ExpiresAt,
	}, nil
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	resp, err := s.svc.ValidateToken(req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.ValidateTokenResponse{
		IsValid: resp.IsValid,
		UserId:  resp.UserID,
	}, nil
}
