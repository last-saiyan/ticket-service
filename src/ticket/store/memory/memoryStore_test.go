package memory

import (
	"testing"

	ticket_v1 "cloudbees.com/ticket/proto/ticket.v1"
)

func TestAdd(t *testing.T) {

	store := NewInMemoryStore(10)
	user := &ticket_v1.User{
		Email:     "user2@gmail.com",
		FirstName: "first",
		LastName:  "last",
	}
	store.PurchaseTicket(&ticket_v1.TicketPurchaseRequest{
		From:      "london",
		To:        "france",
		PricePaid: 12,
		User:      user,
	})

	if (store.firstSectionSize + store.secondSectionSize) != 1 {
		t.Errorf("store.firstSectionSize + store.secondSectionSize is not 1 ")
	}
	_, ok := store.data["user2@gmail.com"]
	if !ok {
		t.Errorf("item missing in store data")
	}
}

func TestGetReceipt(t *testing.T) {
	store := NewInMemoryStore(10)
	user := &ticket_v1.User{
		Email:     "user2@gmail.com",
		FirstName: "first",
		LastName:  "last",
	}
	store.PurchaseTicket(&ticket_v1.TicketPurchaseRequest{
		From:      "london",
		To:        "france",
		PricePaid: 12,
		User:      user,
	})

	ticketInfo, err := store.GetReceipt(&ticket_v1.ReceiptRequest{
		Email: "user2@gmail.com",
	})
	if err != nil {
		t.Errorf("received error on GetReceipt %v", err)
	}

	if ticketInfo.User.Email != "user2@gmail.com" {
		t.Errorf("received invalid email %s", err)
	}
	// todo assert other properties on this test
}

func TestDeleteUser(t *testing.T) {
	store := NewInMemoryStore(10)
	user := &ticket_v1.User{
		Email:     "user2@gmail.com",
		FirstName: "first",
		LastName:  "last",
	}
	store.PurchaseTicket(&ticket_v1.TicketPurchaseRequest{
		From:      "london",
		To:        "france",
		PricePaid: 12,
		User:      user,
	})
	err := store.DeleteUser(&ticket_v1.RemoveUserRequest{
		Email: "user2@gmail.com",
	})

	if err != nil {
		t.Errorf("received error on DeleteUser %v", err)
	}

	if (store.firstSectionSize + store.secondSectionSize) != 0 {
		t.Errorf("store.firstSectionSize + store.secondSectionSize is not 0 ")
	}

	_, ok := store.data["user2@gmail.com"]
	if ok {
		t.Errorf("item exists in store data after delete")
	}
}

func TestUpdateSeat(t *testing.T) {

	store := NewInMemoryStore(10)
	user := &ticket_v1.User{
		Email:     "user2@gmail.com",
		FirstName: "first",
		LastName:  "last",
	}
	store.PurchaseTicket(&ticket_v1.TicketPurchaseRequest{
		From:      "london",
		To:        "france",
		PricePaid: 12,
		User:      user,
	})

	store.UpdateSeat(&ticket_v1.ModifySeatRequest{
		Email:         "user2@gmail.com",
		NewSection:    1,
		NewSeatNumber: 3,
	})

	ticketInfo, _ := store.GetReceipt(&ticket_v1.ReceiptRequest{
		Email: "user2@gmail.com",
	})

	if ticketInfo.Section != 1 {
		t.Errorf("person secton did not change")
	}
	if ticketInfo.SeatNumber != 3 {
		t.Errorf("person seat number did not change")
	}
}
