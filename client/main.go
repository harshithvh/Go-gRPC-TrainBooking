package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/harshithvh/go_gRPC/proto"
	"google.golang.org/grpc"
)

const serverAddress = "localhost:8080"

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := proto.NewTicketServiceClient(conn)

	purchaseTicket(client)
	time.Sleep(3 * time.Second)
	allocateSeat(client)
	showReceipt(client)
	getUsersBySection(client)
	removeUser(client)
	modifySeat(client)

}

func purchaseTicket(client proto.TicketServiceClient) {
	// Example of using PurchaseTicket function
	purchaseRequest := &proto.PurchaseRequest{
		From: "London",
		To:   "France",
		User: &proto.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe198@gmail.com",
		},
	}

	purchaseResponse, err := client.PurchaseTicket(context.Background(), purchaseRequest)
	if err != nil {
		log.Fatalf("Error calling PurchaseTicket: %v", err)
	}

	// Convert the PurchaseResponse to indented JSON format
	jsonResponse, err := json.MarshalIndent(purchaseResponse, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling PurchaseResponse to JSON: %v", err)
	}

	// Print
	fmt.Printf("PurchaseTicket Response:\n%s\n", jsonResponse)
}

func allocateSeat(client proto.TicketServiceClient) {

	// Example of using AllocateSeat function
	allocateSeatRequest := &proto.AllocateSeatRequest{
		Email: "john.doe198@gmail.com",
		Section: "A",
	}

	allocateSeatResponse, err := client.AllocateSeat(context.Background(), allocateSeatRequest)
	if err != nil {
		log.Fatalf("Error calling AllocateSeat: %v", err)
	}

	// Convert the AllocateSeatResponse to indented JSON format
	allocateSeatJSON, err := json.MarshalIndent(allocateSeatResponse, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling AllocateSeatResponse to JSON: %v", err)
	}

	// Print
	fmt.Printf("AllocateSeat Response:\n%s\n", allocateSeatJSON)
}

func showReceipt(client proto.TicketServiceClient) {
	// Example of using ShowReceipt function
	showReceiptRequest := &proto.ShowReceiptRequest{
		Email: "john.doe198@gmail.com",
	}

	showReceiptResponse, err := client.ShowReceipt(context.Background(), showReceiptRequest)
	if err != nil {
		log.Fatalf("Error calling ShowReceipt: %v", err)
	}

	// Convert the ShowReceiptResponse to indented JSON format
	showReceiptJSON, err := json.MarshalIndent(showReceiptResponse, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling ShowReceiptResponse to JSON: %v", err)
	}
	// Print
	fmt.Printf("ShowReceipt Response:\n%s\n", showReceiptJSON)
}

func getUsersBySection(client proto.TicketServiceClient) {

	// Example of using GetUsersBySection function
	getUsersBySectionRequest := &proto.GetUsersBySectionRequest{
		Section: "A",
	}

	getUsersBySectionResponse, err := client.GetUsersBySection(context.Background(), getUsersBySectionRequest)
	if err != nil {
		log.Fatalf("Error calling GetUsersBySection: %v", err)
	}

	// Convert the GetUsersBySectionResponse to indented JSON format
	getUsersBySectionJSON, err := json.MarshalIndent(getUsersBySectionResponse, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling GetUsersBySectionResponse to JSON: %v", err)
	}
	fmt.Printf("GetUsersBySection Response:\n%s\n", getUsersBySectionJSON)
}

func removeUser(client proto.TicketServiceClient) {

	// Example of using RemoveUser function
	removeUserRequest := &proto.RemoveUserRequest{
		Email: "john.doe198@gmail.com",
	}

	removeUserResponse, err := client.RemoveUser(context.Background(), removeUserRequest)
	if err != nil {
		log.Fatalf("Error calling RemoveUser: %v", err)
	}

	// Convert the RemoveUserResponse to indented JSON format
	removeUserJSON, err := json.MarshalIndent(removeUserResponse, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling RemoveUserResponse to JSON: %v", err)
	}
	fmt.Printf("RemoveUser Response:\n%s\n", removeUserJSON)
}

func modifySeat(client proto.TicketServiceClient) {

	// Example of using ModifySeat function
	modifySeatRequest := &proto.ModifySeatRequest{
		Email:            "john.doe@gmail.com",
		NewSection:       "B",
		NewSeatNumber:    5,
	}

	modifySeatResponse, err := client.ModifySeat(context.Background(), modifySeatRequest)
	if err != nil {
		log.Fatalf("Error calling ModifySeat: %v", err)
	}

	// Convert the ModifySeatResponse to indented JSON format
	modifySeatJSON, err := json.MarshalIndent(modifySeatResponse, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling ModifySeatResponse to JSON: %v", err)
	}
	fmt.Printf("ModifySeat Response (Indented JSON):\n%s\n", modifySeatJSON)

}
