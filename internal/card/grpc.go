package card

import (
	"context"

	pb "github.com/zakirkun/banking-microservices/api/proto"
)

type GRPCServer struct {
	pb.UnimplementedCardServiceServer
	svc Service
}

func NewGRPCServer(svc Service) *GRPCServer {
	return &GRPCServer{svc: svc}
}

func (s *GRPCServer) IssueCard(ctx context.Context, req *pb.IssueCardRequest) (*pb.CardResponse, error) {
	resp, err := s.svc.IssueCard(
		req.CustomerId,
		req.AccountId,
		req.CardType,
		req.CardNetwork,
	)
	if err != nil {
		return nil, err
	}

	return s.cardToProto(resp), nil
}

func (s *GRPCServer) GetCard(ctx context.Context, req *pb.GetCardRequest) (*pb.CardResponse, error) {
	resp, err := s.svc.GetCard(req.CardId)
	if err != nil {
		return nil, err
	}

	return s.cardToProto(resp), nil
}

func (s *GRPCServer) BlockCard(ctx context.Context, req *pb.BlockCardRequest) (*pb.CardResponse, error) {
	resp, err := s.svc.BlockCard(req.CardId, req.Reason)
	if err != nil {
		return nil, err
	}

	return s.cardToProto(resp), nil
}

func (s *GRPCServer) UnblockCard(ctx context.Context, req *pb.UnblockCardRequest) (*pb.CardResponse, error) {
	resp, err := s.svc.UnblockCard(req.CardId)
	if err != nil {
		return nil, err
	}

	return s.cardToProto(resp), nil
}

func (s *GRPCServer) ListCards(ctx context.Context, req *pb.ListCardsRequest) (*pb.ListCardsResponse, error) {
	resp, err := s.svc.ListCards(req.CustomerId, int(req.Page), int(req.Limit), req.Status)
	if err != nil {
		return nil, err
	}

	var cards []*pb.CardResponse
	for _, card := range resp.Cards {
		cards = append(cards, s.cardToProto(&card))
	}

	return &pb.ListCardsResponse{
		Cards: cards,
		Total: int32(resp.Total),
	}, nil
}

func (s *GRPCServer) cardToProto(c *CardResponse) *pb.CardResponse {
	return &pb.CardResponse{
		CardId:      c.ID,
		CustomerId:  c.CustomerID,
		AccountId:   c.AccountID,
		CardNumber:  c.CardNumber,
		CardType:    c.CardType,
		CardNetwork: c.CardNetwork,
		ExpiryDate:  c.ExpiryDate,
		Cvv:         c.CVV,
		Status:      c.Status,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
