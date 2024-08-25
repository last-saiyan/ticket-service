FROM golang:1.20-alpine AS builder

RUN apk add --no-cache  protobuf

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

    
WORKDIR /app

COPY . .

RUN go mod download

RUN cd proto && protoc --go_out=. --go-grpc_out=. ./tickets.proto 

RUN cd src && go build -o tickets .

ENTRYPOINT cd src && ./tickets

# todo use multistage builds