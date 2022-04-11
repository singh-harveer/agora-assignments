package main

import (
	"agora/assignments/agora"
	"agora/assignments/service"
	"context"
	"fmt"
	"log"
)

const (
	port       = 8000
	testUserID = "60d0fe4f5311236168a109ca"
)

func main() {

	log.Println("--------creating new GRPC users service client---------")
	var userClient, err = service.NewUserClient(context.Background(), fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to create client :%v", err)
	}
	log.Println("------------retrieving user details by given user id--------------")

	var user agora.User
	user, err = userClient.GetUserByID(context.Background(), testUserID)
	if err != nil {
		log.Fatalf("failed to get user by id : %v", err)
	}

	log.Printf("user details for userID: %s\n", user.ID)
	fmt.Println(user)

	log.Println("------------retrieving all users with pagination--------------")

	var users []agora.User
	users, err = userClient.GetUsers(context.Background(), 1, 20)
	if err != nil {
		log.Fatalf("failed to get users : %v", err)
	}

	log.Printf("users details:\n %v", users)
}
