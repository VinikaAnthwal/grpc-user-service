// client/main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	pb "grpc-user-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for {
		fmt.Println("Select an option:")
		fmt.Println("1. Get user by ID")
		fmt.Println("2. Get users by IDs")
		fmt.Println("3. Get all users")
		fmt.Println("4. Search users by criteria")
		fmt.Println("5. Exit")

		choice := readIntInput("Enter your choice (1-5): ", 1, 5)

		switch choice {
		case 1:
			getUserByID(client)
		case 2:
			getUsersByIDs(client)
		case 3:
			getAllUsers(client)
		case 4:
			searchUsers(client)
		case 5:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 5.")
		}
	}
}

func getUserByID(client pb.UserServiceClient) {
	id := readIntInput("Enter user ID: ", 1, 1000)
	req := &pb.GetUserRequest{
		Id: int32(id),
	}
	resp, err := client.GetUserByID(context.Background(), req)
	if err != nil {
		handleRPCError(err)
		return
	}
	fmt.Printf("User details: %+v\n", resp.GetUser())
}

func getUsersByIDs(client pb.UserServiceClient) {
	ids := readIDsInput()

	req := &pb.GetUsersRequest{
		Ids: ids,
	}
	resp, err := client.GetUsersByIDs(context.Background(), req)
	if err != nil {
		handleRPCError(err)
		return
	}

	fmt.Println("Users found:")
	for _, user := range resp.GetUsers() {
		fmt.Printf("%+v\n", user)
	}
	if len(resp.GetNotFoundIds()) > 0 {
		fmt.Printf("User IDs not found: %v\n", resp.GetNotFoundIds())
	}
}

func getAllUsers(client pb.UserServiceClient) {
    req := &pb.GetAllUsersRequest{}
    resp, err := client.GetAllUsers(context.Background(), req)
    if err != nil {
        handleRPCError(err)
        return
    }

    fmt.Println("All users:")
    for _, user := range resp.GetUsers() {
        fmt.Printf("id:%d fname:%q city:%q phone:%d height:%.1f married:%t\n",
            user.Id, user.Fname, user.City, user.Phone, user.Height, user.Married)
    }
}


func searchUsers(client pb.UserServiceClient) {
    fmt.Println("Select search criteria:")
    fmt.Println("1. By First Name")
    fmt.Println("2. By City")
    fmt.Println("3. By Phone Number")
    fmt.Println("4. By Marital Status")

    searchChoice := readIntInput("Enter your choice (1-5): ", 1, 5)
    query := readStringInput("Enter search query: ")

    // Validate user input based on selected search criteria
    switch searchChoice {
    case 1:
        if len(query) < 3 || len(query) > 20 {
            fmt.Println("First name should be between 3 and 20 characters.")
            return
        }
    case 2:
        if len(query) < 1 || len(query) > 20 {
            fmt.Println("City should be between 1 and 20 characters.")
            return
        }
    case 3:
        if len(query) != 10 {
            fmt.Println("Phone number should be 10 digits.")
            return
        }
    case 4:
        if query != "true" && query != "false" {
            fmt.Println("Invalid marital status format. Marital status should be true or false.")
            return
        }
    default:
        fmt.Println("Invalid search option.")
        return
    }

    req := &pb.SearchUserRequest{Query: query}
    resp, err := client.SearchUsers(context.Background(), req)
    if err != nil {
        handleRPCError(err)
        return
    }

    fmt.Println("Search results:")
    if len(resp.GetUsers()) == 0 {
        fmt.Println("No users found matching the criteria.")
    } else {
        for _, user := range resp.GetUsers() {
            fmt.Printf("id:%d fname:%q city:%q phone:%d height:%.1f married:%t\n",
                user.Id, user.Fname, user.City, user.Phone, user.Height, user.Married)
        }
    }
}

func readIntInput(prompt string, min, max int) int {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number.")
			continue
		}
		if num < min || num > max {
			fmt.Printf("Input out of range (%d-%d). Please enter a number between %d and %d.\n", min, max, min, max)
			continue
		}
		return num
	}
}

func readIDsInput() []int32 {
	var ids []int32
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter user IDs (comma-separated): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	input = strings.TrimSpace(input)
	if input == "" {
		fmt.Println("No input provided.")
		return ids
	}
	idStrs := strings.Split(input, ",")
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			fmt.Printf("Invalid input: '%s'. Please enter comma-separated integers.\n", idStr)
			continue
		}
		if id <= 0 {
			fmt.Printf("Invalid user ID: %d. Please enter positive integers only.\n", id)
			continue
		}
		ids = append(ids, int32(id))
	}
	return ids
}

func readStringInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	return strings.TrimSpace(input)
}

func handleRPCError(err error) {
	statusErr, ok := status.FromError(err)
	if ok && statusErr.Code() == codes.NotFound {
		fmt.Println("User not found.")
	} else {
		log.Fatalf("RPC failed: %v", err)
	}
}
