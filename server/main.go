package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/google/uuid"
	pb "github.com/harshithvh/go_gRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
    userInfo map[string]*pb.Receipt 
    seatAvailabilityA [10]bool
	seatAvailabilityB [10]bool
	pb.UnimplementedTicketServiceServer
}

// Helper function to find the next available seat in a section
func findNextAvailableSeat(seatAvailability *[10]bool) (int, bool) {
	for seatNumber, available := range seatAvailability {
		if !available {
			return seatNumber, true
		}
	}
	return -1, false
}

// gRPC methods:
func (s *Server) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	// Validate the request
	if req == nil || req.User == nil || req.User.Email == "" || req.User.FirstName == "" || req.User.LastName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: First name, last name, and email cannot be empty")
	}

	price := 20.0

	purchaseID := uuid.New().String()

	// Check if a ticket with the same email already exists
    _, exists := s.userInfo[req.User.Email]
    if exists {
        return nil, status.Errorf(codes.AlreadyExists, "Ticket already purchased for the provided email: %s", req.User.Email)
    }

	// Create a PurchaseResponse
	purchaseResponse := &pb.PurchaseResponse{
		From:       req.From,
		To:         req.To,
		User:       req.User,
		PricePaid:  price,
		PurchaseId: purchaseID,
	}

	ticketInfo := &pb.Receipt{
		From:       req.From,
		To:         req.To,
		User:       req.User,
		PricePaid:  float32(price),
		PurchaseId: purchaseID,
		Seat:       &pb.Seat{},
	}

	// Store the purchaseResponse in the Server's in-memory storage
	s.userInfo[req.User.Email] = ticketInfo

	return purchaseResponse, nil
}

func (s *Server) AllocateSeat(ctx context.Context, req *pb.AllocateSeatRequest) (*pb.AllocateSeatResponse, error) {
	// Validate the request
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	purchaseInfo, exists := s.userInfo[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Purchase not found for the provided email")
	}

	// Check if the seat is already allocated for the user
	if purchaseInfo.Seat.Section != "" && purchaseInfo.Seat.SeatNumber > 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "Seat already allocated for the user with email: %s", req.Email)
	}

	if req == nil || req.Section == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Section cannot be empty")
	}

	if req.Section != "A" && req.Section != "B" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid new section: %s", req.Section)
	}

	var seatAvailability *[10]bool
	var seatNumber int

	switch req.Section {
	case "A":
		seatAvailability = &s.seatAvailabilityA
		seat, available := findNextAvailableSeat(&s.seatAvailabilityA)
		if !available {
			return nil, status.Errorf(codes.ResourceExhausted, "No more seats available in section A")
		}
		seatNumber = seat
	case "B":
		seatAvailability = &s.seatAvailabilityB
		seat, available := findNextAvailableSeat(&s.seatAvailabilityB)
		if !available {
			return nil, status.Errorf(codes.ResourceExhausted, "No more seats available in section B")
		}
		seatNumber = seat
	default:
		return nil, status.Errorf(codes.InvalidArgument, "Invalid section: %s", req.Section)
	}

	// Mark the seat as unavailable
	(*seatAvailability)[seatNumber] = true

	// Update the PurchaseResponse with the allocated seat information
	purchaseInfo.Seat.Section = req.Section
	purchaseInfo.Seat.SeatNumber = int32(seatNumber + 1)

	s.userInfo[req.Email] = purchaseInfo

	// Create an AllocateSeatResponse with the allocated seat information
	allocateSeatResponse := &pb.AllocateSeatResponse{
		Email:      req.Email,
		Section:    req.Section,
		SeatNumber: int32(seatNumber+1),
	}

	return allocateSeatResponse, nil
}

func (s *Server) ShowReceipt(ctx context.Context, req *pb.ShowReceiptRequest) (*pb.ShowReceiptResponse, error) {
	// Validate the request
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Email cannot be empty")
	}

	// Retrieve the purchase response based on the user's email
	receiptInfo, exists := s.userInfo[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Purchase not found for the provided email")
	}

	// Check if the section and seat number is allocated for the user
	if receiptInfo.Seat.Section == "" && receiptInfo.Seat.SeatNumber == 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "No section and seat number allocated for the user with email: %s", req.Email)
	}

	// Create a ShowReceiptResponse
	showReceiptResponse := &pb.ShowReceiptResponse{
		UserInfo: receiptInfo,
	}

	return showReceiptResponse, nil
}

func (s *Server) GetUsersBySection(ctx context.Context, req *pb.GetUsersBySectionRequest) (*pb.GetUsersBySectionResponse, error) {
	// Validate the request
	if req == nil || req.Section == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Section cannot be empty")
	}

	// Initialize a list to store UserSeatInfo for the requested section
	usersBySection := []*pb.Receipt{}

	// Iterate through stored tickets and collect users with the requested section
	for _, receiptInfo := range s.userInfo {
		if receiptInfo.Seat.Section == req.Section {
			usersBySection = append(usersBySection, receiptInfo)
		}
	}

	// Create a GetUsersBySectionResponse
	getUsersBySectionResponse := &pb.GetUsersBySectionResponse{
		UserInfo: usersBySection,
	}

	return getUsersBySectionResponse, nil
}

func (s *Server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	// Validate the request
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Email cannot be empty")
	}

	// Check if the user exists in the stored tickets
	purchaseResponse, exists := s.userInfo[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "User removed or not present")
	}

	// Mark the current seat and seat number as available
    var currSeat *[10]bool

        if purchaseResponse.Seat.Section == "A" {
            currSeat = &s.seatAvailabilityA
        } else {
            currSeat = &s.seatAvailabilityB
        }

    	currentSeatNumber := int32(purchaseResponse.Seat.SeatNumber)
    	(*currSeat)[currentSeatNumber-1] = false

    	// Remove the user from the stored tickets
    	delete(s.userInfo, req.Email)

        // Create a RemoveUserResponse indicating success
        removeUserResponse := &pb.RemoveUserResponse{
            Res: "User removed successfully",
        }

    return removeUserResponse, nil
}


func (s *Server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	// Validate the request
	if req == nil || req.Email == "" || req.NewSection == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Email and new section cannot be empty")
	}

	// Check if a purchase for the given email exists
	purchaseResponse, exists := s.userInfo[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "No purchase found for the provided email")
	}

    // Check if the requested new seat number is within the valid range (1 to 10)
	if req.NewSeatNumber < 1 || req.NewSeatNumber > 10 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid new seat number. Must be between 1 and 10")
	}

    section := req.NewSection

	var seatAvailability *[10]bool
	if section == "A" {
		seatAvailability = &s.seatAvailabilityA
	} else if section == "B" {
		seatAvailability = &s.seatAvailabilityB
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid new section: %s", req.NewSection)
	}

	if (*seatAvailability)[req.NewSeatNumber-1] {
		return nil, status.Errorf(codes.ResourceExhausted, "Requested seat is not available in the specified section")
	}

	// Mark the current seat and seat number as available
    var currSeat *[10]bool

    if purchaseResponse.Seat.Section == "A" {
        currSeat = &s.seatAvailabilityA
    } else {
        currSeat = &s.seatAvailabilityB
    }

	currentSeatNumber := int32(purchaseResponse.Seat.SeatNumber)
	(*currSeat)[currentSeatNumber-1] = false

	// Mark the new seat and seat number as unavailable
	(*seatAvailability)[req.NewSeatNumber-1] = true

	// Update the seat number in the purchase response
	purchaseResponse.Seat.SeatNumber = int32(req.NewSeatNumber)

    // Update the section in the purchase response
    purchaseResponse.Seat.Section = section

	// Create a ModifySeatResponse indicating success
	modifySeatResponse := &pb.ModifySeatResponse{Res: "Seat modified successfully"}

	return modifySeatResponse, nil
}

func main() {
    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    service := &Server{
        userInfo: make(map[string]*pb.Receipt),
    }
    pb.RegisterTicketServiceServer(s, service)

    go func() {
        if err := s.Serve(lis); err != nil {
            log.Fatalf("failed to serve: %v", err)
        }
    }()

    log.Println("Server is running on :8080")

    // Ctrl+C to stop the server
    ch := make(chan os.Signal, 1)
    signal.Notify(ch, os.Interrupt)
    <-ch

    log.Println("Stopping the Server...")
    s.GracefulStop()
    log.Println("Server stopped")
}