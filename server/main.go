package main

import (
	"context"
	"fmt"
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
    tickets map[string]*pb.PurchaseResponse 
    seatAvailabilityA [10]bool
	seatAvailabilityB [10]bool
	pb.UnimplementedTicketServiceServer
}

// Helper function to extract section and seat number from PurchaseResponse
func extractSectionAndSeat(pr *pb.PurchaseResponse) (string, int32) {
	if pr == nil || pr.SeatNumber < 0 {
		return "", -1
	}

	// Convert seat number to int32
	seatNumber := pr.SeatNumber

	// Extract section from the combined section and seat number string
	section := pr.Section

	return section, seatNumber
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

	var section string
	var seatNumber int

	// Check if a ticket with the same email already exists
    _, exists := s.tickets[req.User.Email]
    if exists {
        return nil, status.Errorf(codes.AlreadyExists, "Ticket already purchased for the provided email: %s", req.User.Email)
    }


	// Allocate seat
	var seatAvailability *[10]bool
	if seat, available := findNextAvailableSeat(&s.seatAvailabilityA); available {
		section = "A"
		seatAvailability = &s.seatAvailabilityA
        seatNumber = seat
	} else if seat, available := findNextAvailableSeat(&s.seatAvailabilityB); available {
		section = "B"
		seatAvailability = &s.seatAvailabilityB
        seatNumber = seat
	} else {
		return nil, status.Errorf(codes.ResourceExhausted, "No more seats available in both sections")
	}

	// Mark the seat as unavailable
	(*seatAvailability)[seatNumber] = true

	// Create a PurchaseResponse
	purchaseResponse := &pb.PurchaseResponse{
		From:       req.From,
		To:         req.To,
		User:       req.User,
		PricePaid:  price,
		PurchaseId: purchaseID,
		Section:    section,
		SeatNumber: int32(seatNumber+1),
	}

	// Store the purchaseResponse in the Server's in-memory storage
	s.tickets[req.User.Email] = purchaseResponse

	return purchaseResponse, nil
}

func (s *Server) AllocateSeat(ctx context.Context, req *pb.AllocateSeatRequest) (*pb.AllocateSeatResponse, error) {
	// Validate the request
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request")
	}

	purchaseResponse, exists := s.tickets[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Purchase not found for the provided email")
	}

	// Extract section and seat number from the stored purchase response
	section, seatNumber := extractSectionAndSeat(purchaseResponse)
	if section == "" || seatNumber == -1 {
		return nil, status.Errorf(codes.Internal, "Error extracting section and seat number")
	}

	// Create an AllocateSeatResponse with the allocated seat information
	allocateSeatResponse := &pb.AllocateSeatResponse{
		Section: section,
		SeatNumber: seatNumber,
	}

	return allocateSeatResponse, nil
}

func (s *Server) ShowReceipt(ctx context.Context, req *pb.ShowReceiptRequest) (*pb.ShowReceiptResponse, error) {
	// Validate the request
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Email cannot be empty")
	}

	// Retrieve the purchase response based on the user's email
	purchaseResponse, exists := s.tickets[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "Purchase not found for the provided email")
	}

	// Extract section and seat number from the stored purchase response
	section, seatNumber := extractSectionAndSeat(purchaseResponse)
	if section == "" || seatNumber == -1 {
		return nil, status.Errorf(codes.Internal, "Error extracting section and seat number")
	}

	// Create a ShowReceiptResponse
	showReceiptResponse := &pb.ShowReceiptResponse{
		PurchaseResponse: purchaseResponse,
		AllocatedSeat:     fmt.Sprintf("%s%d", section, seatNumber),
	}

	return showReceiptResponse, nil
}

func (s *Server) GetUsersBySection(ctx context.Context, req *pb.GetUsersBySectionRequest) (*pb.GetUsersBySectionResponse, error) {
	// Validate the request
	if req == nil || req.Section == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Section cannot be empty")
	}

	// Initialize a list to store UserSeatInfo for the requested section
	usersBySection := []*pb.UserSeatInfo{}

	// Iterate through stored tickets and collect users with the requested section
	for _, purchaseResponse := range s.tickets {
		if purchaseResponse.Section == req.Section {
			userSeatInfo := &pb.UserSeatInfo{
				User:       purchaseResponse.User,
				Section:    purchaseResponse.Section,
				SeatNumber: purchaseResponse.SeatNumber,
			}
			usersBySection = append(usersBySection, userSeatInfo)
		}
	}

	// Create a GetUsersBySectionResponse
	getUsersBySectionResponse := &pb.GetUsersBySectionResponse{
		UserSeatInfo: usersBySection,
	}

	return getUsersBySectionResponse, nil
}

func (s *Server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	// Validate the request
	if req == nil || req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: Email cannot be empty")
	}

	// Check if the user exists in the stored tickets
	purchaseResponse, exists := s.tickets[req.Email]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "User removed or not present")
	}

	// Mark the current seat and seat number as available
    var currSeat *[10]bool

        if purchaseResponse.Section == "A" {
            currSeat = &s.seatAvailabilityA
        } else {
            currSeat = &s.seatAvailabilityB
        }

    	currentSeatNumber := int32(purchaseResponse.SeatNumber)
    	(*currSeat)[currentSeatNumber-1] = false

    	// Remove the user from the stored tickets
    	delete(s.tickets, req.Email)

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
	purchaseResponse, exists := s.tickets[req.Email]
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

    if purchaseResponse.Section == "A" {
        currSeat = &s.seatAvailabilityA
    } else {
        currSeat = &s.seatAvailabilityB
    }

	currentSeatNumber := int32(purchaseResponse.SeatNumber)
	(*currSeat)[currentSeatNumber-1] = false

	// Mark the new seat and seat number as unavailable
	(*seatAvailability)[req.NewSeatNumber-1] = true

	// Update the seat number in the purchase response
	purchaseResponse.SeatNumber = int32(req.NewSeatNumber)

    // Update the section in the purchase response
    purchaseResponse.Section = section

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
        tickets: make(map[string]*pb.PurchaseResponse),
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