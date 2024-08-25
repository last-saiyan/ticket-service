package memory

import (
	"errors"
	"sync"

	ticket_v1 "cloudbees.com/ticket/proto/ticket.v1"
)

type InMemoryStore struct {
	mutex             sync.RWMutex
	data              map[string]*ticket_v1.TicketPurchaseRequest
	sections          [2][]string
	sectionSize       int
	firstSectionSize  int
	secondSectionSize int
}

func NewInMemoryStore(sectionSize int) *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]*ticket_v1.TicketPurchaseRequest),
		sections: [2][]string{
			make([]string, sectionSize),
			make([]string, sectionSize),
		},
		sectionSize:       sectionSize,
		firstSectionSize:  0, // todo tie this up with the sections
		secondSectionSize: 0,
	}
}

// compile time interface check
// var _ store.TicketInterface = InMemoryStore{}

func (store *InMemoryStore) PurchaseTicket(ticketPurchase *ticket_v1.TicketPurchaseRequest) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, ok := store.data[ticketPurchase.User.Email]
	if ok {
		return errors.New("user with this email has a ticket")
	}
	if store.firstSectionSize+store.secondSectionSize >= 2*store.sectionSize {
		return errors.New("no space left")
	}

	store.data[ticketPurchase.User.Email] = ticketPurchase

	var smallerSection []string
	if store.firstSectionSize > store.secondSectionSize {
		store.secondSectionSize++
		smallerSection = store.sections[1]
	} else {
		store.firstSectionSize++
		smallerSection = store.sections[0]
	}

	for i, seat := range smallerSection {
		if seat == "" {
			smallerSection[i] = ticketPurchase.User.Email
			break
		}
	}

	return nil
}

func (store *InMemoryStore) DeleteUser(deleteUserReq *ticket_v1.RemoveUserRequest) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, ok := store.data[deleteUserReq.Email]
	if !ok {
		return errors.New("user does not exists")
	}

	delete(store.data, deleteUserReq.Email)

	for i, section := range store.sections {
		for j, seat := range section {
			if seat == deleteUserReq.Email {
				section[j] = ""
				if i == 1 {
					store.secondSectionSize--
				} else {
					store.firstSectionSize--
				}
			}
		}
	}
	return nil
}

func (store *InMemoryStore) UpdateSeat(seatUpdateReq *ticket_v1.ModifySeatRequest) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, ok := store.data[seatUpdateReq.Email]
	if !ok {
		return errors.New("user does not exists")
	}

	var requiredSection []string

	if seatUpdateReq.NewSection < int32(len(store.sections)) {
		requiredSection = store.sections[seatUpdateReq.NewSection]
	} else {
		return errors.New("section does not exists")
	}

	if seatUpdateReq.NewSeatNumber < int32(len(requiredSection)) {
		if requiredSection[seatUpdateReq.NewSeatNumber] == "" {

			for i, section := range store.sections {
				for j, seat := range section {
					if seat == seatUpdateReq.Email {
						section[j] = ""
						if i == 0 {
							store.firstSectionSize--
						} else {
							store.secondSectionSize--
						}
					}
				}
			}
			requiredSection[seatUpdateReq.NewSeatNumber] = seatUpdateReq.Email

			if seatUpdateReq.NewSection == 0 {
				store.firstSectionSize++
			} else {
				store.secondSectionSize++
			}
		} else {
			return errors.New("seat occupied")
		}
	} else {
		return errors.New("seat does not exists")
	}

	return nil
}

func (store *InMemoryStore) GetReceipt(receiptReq *ticket_v1.ReceiptRequest) (*ticket_v1.ReceiptResponse, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	ticketInfo, ok := store.data[receiptReq.Email]
	if !ok {
		return nil, errors.New("user does not exists")
	}

	var sectionId, seatId int
	for i, section := range store.sections {
		for j, seat := range section {
			if seat == receiptReq.Email {
				seatId = j
				sectionId = i
			}
		}
	}

	return &ticket_v1.ReceiptResponse{
		From:       ticketInfo.From,
		To:         ticketInfo.To,
		PricePaid:  ticketInfo.PricePaid,
		User:       ticketInfo.User,
		Section:    int64(sectionId),
		SeatNumber: int64(seatId),
	}, nil

}

func (store *InMemoryStore) GetUserBySection(receiptReq *ticket_v1.SectionRequest) (*ticket_v1.SectionResponse, error) {

	if receiptReq.Section >= int64(len(store.sections)) {
		return nil, errors.New("section does not exist")
	} else {

		var response []*ticket_v1.UserSeat
		for i, seatInfo := range store.sections[receiptReq.Section] {
			if seatInfo != "" {
				userSeat := &ticket_v1.UserSeat{
					User:        seatInfo,
					Section:     receiptReq.Section,
					Seat_Number: int64(i),
				}
				response = append(response, userSeat)
			}
		}
		return &ticket_v1.SectionResponse{
			Users: response,
		}, nil
	}
}
