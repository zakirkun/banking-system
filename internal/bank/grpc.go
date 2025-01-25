package bank

import (
	"context"

	pb "github.com/zakirkun/banking-microservices/api/proto"
)

type GRPCServer struct {
	pb.UnimplementedBankServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) CreateBank(ctx context.Context, req *pb.CreateBankRequest) (*pb.BankResponse, error) {
	resp, err := s.svc.CreateBank(
		req.BankName,
		req.BankCode,
		req.SwiftCode,
		req.Country,
		req.Currency,
	)
	if err != nil {
		return nil, err
	}

	return s.bankToProto(resp), nil
}

func (s *GRPCServer) GetBank(ctx context.Context, req *pb.GetBankRequest) (*pb.BankResponse, error) {
	resp, err := s.svc.GetBank(req.BankId)
	if err != nil {
		return nil, err
	}

	return s.bankToProto(resp), nil
}

func (s *GRPCServer) UpdateBank(ctx context.Context, req *pb.UpdateBankRequest) (*pb.BankResponse, error) {
	resp, err := s.svc.UpdateBank(req.BankId, req.BankName, req.SwiftCode, req.Status)
	if err != nil {
		return nil, err
	}

	return s.bankToProto(resp), nil
}

func (s *GRPCServer) DeleteBank(ctx context.Context, req *pb.DeleteBankRequest) (*pb.DeleteBankResponse, error) {
	err := s.svc.DeleteBank(req.BankId)
	if err != nil {
		return &pb.DeleteBankResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteBankResponse{
		Success: true,
		Message: "Bank deleted successfully",
	}, nil
}

func (s *GRPCServer) ListBanks(ctx context.Context, req *pb.ListBanksRequest) (*pb.ListBanksResponse, error) {
	resp, err := s.svc.ListBanks(int(req.Page), int(req.Limit), req.Search, req.Country, req.Status)
	if err != nil {
		return nil, err
	}

	var banks []*pb.BankResponse
	for _, bank := range resp.Banks {
		banks = append(banks, s.bankToProto(&bank))
	}

	return &pb.ListBanksResponse{
		Banks: banks,
		Total: int32(resp.Total),
	}, nil
}

func (s *GRPCServer) bankToProto(b *BankResponse) *pb.BankResponse {
	return &pb.BankResponse{
		BankId:    b.ID,
		BankName:  b.BankName,
		BankCode:  b.BankCode,
		SwiftCode: b.SwiftCode,
		Country:   b.Country,
		Currency:  b.Currency,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
