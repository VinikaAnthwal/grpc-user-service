// server\main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"

	pb "grpc-user-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedUserServiceServer
	users map[int32]*pb.User
	mu    sync.RWMutex
}

func newServer() *server {
	return &server{
		users: map[int32]*pb.User{
			1: {Id: 1, Fname: "Steve", City: "LA", Phone: 1234567890, Height: 5.8, Married: true},
			2: {Id: 2, Fname: "Jane", City: "NY", Phone: 1234567891, Height: 5.5, Married: false},
			3: {Id: 3, Fname: "Alice", City: "LA", Phone: 1234567892, Height: 5.6, Married: true},
		},
		mu: sync.RWMutex{},
	}
}

func (s *server) GetUserByID(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user ID")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("user with ID %d not found", req.Id))
	}
	return &pb.GetUserResponse{User: user}, nil
}

func (s *server) GetUsersByIDs(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	if len(req.Ids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "user IDs cannot be empty")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	var users []*pb.User
	var foundIds []int32
	for _, id := range req.Ids {
		if id <= 0 {
			return nil, status.Error(codes.InvalidArgument, "invalid user ID")
		}
		if user, exists := s.users[id]; exists {
			users = append(users, user)
			foundIds = append(foundIds, id)
		}
	}
	notFoundIds := findMissingIDs(req.Ids, foundIds)
	if len(notFoundIds) > 0 {
		log.Printf("Some user IDs not found: %v", notFoundIds)
	}

	return &pb.GetUsersResponse{Users: users, NotFoundIds: notFoundIds}, nil
}

func (s *server) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    var users []*pb.User
    for _, user := range s.users {
        users = append(users, user)
    }

    return &pb.GetAllUsersResponse{Users: users}, nil
}

func (s *server) SearchUsers(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
    // Validate the search query
    if req.Query == "" {
        return nil, status.Error(codes.InvalidArgument, "search query cannot be empty")
    }

    // Validate search input data
    if len(req.Query) < 1 || len(req.Query) > 100 {
        return nil, status.Error(codes.InvalidArgument, "search query should be between 1 and 100 characters")
    }

    s.mu.RLock()
    defer s.mu.RUnlock()

    var matchedUsers []*pb.User
    query := strings.ToLower(req.Query)
    for _, user := range s.users {
        if strings.Contains(strings.ToLower(user.Fname), query) ||
            strings.Contains(strings.ToLower(user.City), query) ||
            strings.Contains(strings.ToLower(fmt.Sprint(user.Phone)), query) ||
            (user.Married && query == "true") || (!user.Married && query == "false") {
            matchedUsers = append(matchedUsers, user)
        }
    }

    if len(matchedUsers) == 0 {
        return nil, status.Error(codes.NotFound, "no users found matching the search criteria")
    }

    return &pb.SearchUserResponse{Users: matchedUsers}, nil
}

func findMissingIDs(allIDs []int32, foundIDs []int32) []int32 {
	found := make(map[int32]bool)
	for _, id := range foundIDs {
		found[id] = true
	}

	var missing []int32
	for _, id := range allIDs {
		if !found[id] {
			missing = append(missing, id)
		}
	}
	return missing
}

func matchesSearchQuery(user *pb.User, query string) bool {
	query = strings.ToLower(query)
	return strings.Contains(strings.ToLower(user.Fname), query) ||
		strings.Contains(strings.ToLower(user.City), query) ||
		strings.Contains(strings.ToLower(fmt.Sprint(user.Phone)), query) ||
		(user.Married && strings.Contains(strings.ToLower("married"), query))
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, newServer())

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
