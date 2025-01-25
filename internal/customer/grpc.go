package customer

import (
	"context"

	pb "github.com/zakirkun/banking-microservices/api/proto"
)

type GRPCServer struct {
	pb.UnimplementedCustomerServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CustomerResponse, error) {
	resp, err := s.svc.CreateCustomer(
		req.UserId,
		req.FullName,
		req.DateOfBirth,
		req.Address,
		req.IdNumber,
		req.IdType,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CustomerResponse{
		CustomerId:  resp.ID,
		UserId:      resp.UserID,
		FullName:    resp.FullName,
		DateOfBirth: resp.DateOfBirth,
		Address:     resp.Address,
		IdNumber:    resp.IDNumber,
		IdType:      resp.IDType,
		Status:      resp.Status,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

func (s *GRPCServer) GetCustomer(ctx context.Context, req *pb.GetCustomerRequest) (*pb.CustomerResponse, error) {
	resp, err := s.svc.GetCustomer(req.CustomerId)
	if err != nil {
		return nil, err
	}

	return &pb.CustomerResponse{
		CustomerId:  resp.ID,
		UserId:      resp.UserID,
		FullName:    resp.FullName,
		DateOfBirth: resp.DateOfBirth,
		Address:     resp.Address,
		IdNumber:    resp.IDNumber,
		IdType:      resp.IDType,
		Status:      resp.Status,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

func (s *GRPCServer) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (*pb.CustomerResponse, error) {
	resp, err := s.svc.UpdateCustomer(req.CustomerId, req.FullName, req.Address)
	if err != nil {
		return nil, err
	}

	return &pb.CustomerResponse{
		CustomerId:  resp.ID,
		UserId:      resp.UserID,
		FullName:    resp.FullName,
		DateOfBirth: resp.DateOfBirth,
		Address:     resp.Address,
		IdNumber:    resp.IDNumber,
		IdType:      resp.IDType,
		Status:      resp.Status,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

func (s *GRPCServer) DeleteCustomer(ctx context.Context, req *pb.DeleteCustomerRequest) (*pb.DeleteCustomerResponse, error) {
	err := s.svc.DeleteCustomer(req.CustomerId)
	if err != nil {
		return &pb.DeleteCustomerResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteCustomerResponse{
		Success: true,
		Message: "Customer deleted successfully",
	}, nil
}

func (s *GRPCServer) ListCustomers(ctx context.Context, req *pb.ListCustomersRequest) (*pb.ListCustomersResponse, error) {
	resp, err := s.svc.ListCustomers(int(req.Page), int(req.Limit), req.Search)
	if err != nil {
		return nil, err
	}

	var customers []*pb.CustomerResponse
	for _, c := range resp.Customers {
		customers = append(customers, &pb.CustomerResponse{
			CustomerId:  c.ID,
			UserId:      c.UserID,
			FullName:    c.FullName,
			DateOfBirth: c.DateOfBirth,
			Address:     c.Address,
			IdNumber:    c.IDNumber,
			IdType:      c.IDType,
			Status:      c.Status,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	return &pb.ListCustomersResponse{
		Customers: customers,
		Total:     int32(resp.Total),
	}, nil
}
