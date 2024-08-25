package ticket

import (
	"context"

	ticket_v1 "cloudbees.com/ticket/proto/ticket.v1"
	"cloudbees.com/ticket/src/ticket/store/memory"
)

type server struct {
	ticket_v1.UnimplementedTicketServiceServer
	store memory.InMemoryStore
}

func NewServer() *server {
	return &server{
		// this 50 is the max store of the section
		// get it from config
		store: *memory.NewInMemoryStore(10),
	}
}

// todo add debug level logging to all the handler
func (s *server) RemoveUser(ctx context.Context, req *ticket_v1.RemoveUserRequest) (*ticket_v1.UserDeleteResponse, error) {
	err := s.store.DeleteUser(req)
	if err != nil {
		return nil, err
	}
	return &ticket_v1.UserDeleteResponse{Message: "User removed successfully"}, nil
}

func (s *server) ModifyUserSeat(ctx context.Context, req *ticket_v1.ModifySeatRequest) (*ticket_v1.SeatUpdateResponse, error) {

	err := s.store.UpdateSeat(req)
	if err != nil {
		return nil, err
	}
	return &ticket_v1.SeatUpdateResponse{Message: "User modified successfully"}, nil
}

func (s *server) GetReceipt(ctx context.Context, req *ticket_v1.ReceiptRequest) (*ticket_v1.ReceiptResponse, error) {
	ticketInfo, err := s.store.GetReceipt(req)
	if err != nil {
		return nil, err
	}
	return ticketInfo, nil
}

func (s *server) PurchaseTicket(ctx context.Context, req *ticket_v1.TicketPurchaseRequest) (*ticket_v1.PaymentResponse, error) {
	err := s.store.PurchaseTicket(req)
	if err != nil {
		return nil, err
	}
	return &ticket_v1.PaymentResponse{}, nil
}

func (s *server) ViewUsersBySection(ctx context.Context, req *ticket_v1.SectionRequest) (*ticket_v1.SectionResponse, error) {
	resp, err := s.store.GetUserBySection(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
