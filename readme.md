

# building the docker image

`docker build -t ticket-service .`

# run docker image

`docker run -p 8081:8080 ticket-service`




# testing grpc using grpcurl

`cd proto` 

note section-id and seat-id are 0 indexed

```
./grpcurl --plaintext -d '{
    "email" : "user@gmail.com"
}' -proto ./tickets.proto  localhost:8080 ticket.v1.TicketService/RemoveUser
```

```
./grpcurl --plaintext -d '{
    "email" : "user@gmail.com",
    "new_section": "0",
    "new_seat_number" : "3"
}' -proto ./tickets.proto  localhost:8080 ticket.v1.TicketService/ModifyUserSeat
```

```
./grpcurl --plaintext -d '{
    "email" : "user@gmail.com"
}' -proto ./tickets.proto  localhost:8080 ticket.v1.TicketService/GetReceipt
```

```
./grpcurl --plaintext -d '{
    "user": {
        "email" : "user@gmail.com",
        "first_name": "first",
        "last_name": "last"
    },
    "from": "london",
    "to" : "france",
    "price_paid": 12
}' -proto ./tickets.proto  localhost:8080 ticket.v1.TicketService/PurchaseTicket
```

```
./grpcurl --plaintext -d '{
    "section" : "section2"
}' -proto ./tickets.proto  localhost:8080 ticket.v1.TicketService/ViewUsersBySection
```