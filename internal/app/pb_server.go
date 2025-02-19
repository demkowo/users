package app

import (
	"net"

	pb "github.com/demkowo/users/internal/generated"
	handler "github.com/demkowo/users/internal/handlers/grpc"
	service "github.com/demkowo/users/internal/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const pbPortNumber = ":50000"

func pbServerStart(serv service.Users) {

	lis, err := net.Listen("tcp", pbPortNumber)
	if err != nil {
		log.Fatal("failed to start gRPC server", err)
	}

	log.Println("gRPC server listen on", pbPortNumber)

	s := grpc.NewServer()
	pb.RegisterUsersServer(s, &handler.UsersServer{Service: serv})

	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to start gRPC server", err)
	}
}
