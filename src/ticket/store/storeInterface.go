package store

import ticket_v1 "cloudbees.com/ticket/proto/ticket.v1"

type TicketInterface interface {
	PurchaseTicket(ticketPurchase *ticket_v1.TicketPurchaseRequest) error

	UpdateSeat(seatUpdateReq *ticket_v1.ModifySeatRequest) error
	GetReceipt(ReceiptReq *ticket_v1.ReceiptRequest) (*ticket_v1.ReceiptResponse, error)
	DeleteUser(DeleteUserReq *ticket_v1.RemoveUserRequest) error
	GetUserBySection(req *ticket_v1.SectionRequest) (*ticket_v1.SectionResponse, error)
}
