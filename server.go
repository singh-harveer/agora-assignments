package main

import (
	"agora/assignments/agora"
	"agora/assignments/gql"
	"agora/assignments/service"
	"agora/assignments/userpb"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"google.golang.org/grpc"
)

const (
	portGRPC = 8000

	graphQLAPIBaseURL = "https://dummyapi.io/data/v1/"
	appID             = "62538e33ac0033bc278e762a" // TODO - keep this in env.
)

func main() {
	var addr = fmt.Sprintf(":%d", portGRPC)
	log.Printf("starting grpc service on port %d\n", portGRPC)

	var listen, err = net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s:\n%v", addr, err)
	}

	var userMngr agora.UserManager
	userMngr, err = gql.NewGraphQLClientUserManager(graphQLAPIBaseURL, appID)
	if err != nil {
		log.Fatalf("failed to initialize new user manager: %v", err)
	}

	var server = grpc.NewServer()
	userpb.RegisterUserServiceServer(server, &service.UserService{Manager: userMngr})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := server.Serve(listen); err != nil {
			log.Fatalf("failed to serve grpc server over port %d : %v", portGRPC, err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)

		server.GracefulStop()
	}()

	wg.Wait()
}
