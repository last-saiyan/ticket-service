package main

import (
	"fmt"
	"log"
	"net"

	ticket_v1 "cloudbees.com/ticket/proto/ticket.v1"
	"cloudbees.com/ticket/src/ticket"
	"google.golang.org/grpc"
)

func main() {

	port := 8080
	// todo use proper loggers
	fmt.Printf("starting the grpc server %d", port )
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	ticket_v1.RegisterTicketServiceServer(s, ticket.NewServer())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
