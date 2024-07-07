package main

import (
	"context"
	"testing"

	pb "grpc-user-service/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetUserByID(t *testing.T) {
	s := newServer()

	tests := []struct {
		id       int32
		expected *pb.User
		code     codes.Code
	}{
		{1, s.users[1], codes.OK},
		{2, s.users[2], codes.OK},
		{3, s.users[3], codes.OK},
		{0, nil, codes.InvalidArgument},
		{99, nil, codes.NotFound},
	}

	for _, test := range tests {
		req := &pb.GetUserRequest{Id: test.id}
		resp, err := s.GetUserByID(context.Background(), req)

		if test.code == codes.OK {
			assert.Nil(t, err)
			assert.Equal(t, test.expected, resp.User)
		} else {
			assert.NotNil(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, test.code, status.Code(err))
		}
	}
}

func TestGetUsersByIDs(t *testing.T) {
	s := newServer()

	tests := []struct {
		ids       []int32
		expected  []*pb.User
		notFound  []int32
		code      codes.Code
	}{
		{[]int32{1, 2}, []*pb.User{s.users[1], s.users[2]}, nil, codes.OK},
		{[]int32{1, 99}, []*pb.User{s.users[1]}, []int32{99}, codes.OK},
		{[]int32{}, nil, nil, codes.InvalidArgument},
	}

	for _, test := range tests {
		req := &pb.GetUsersRequest{Ids: test.ids}
		resp, err := s.GetUsersByIDs(context.Background(), req)

		if test.code == codes.OK {
			assert.Nil(t, err)
			assert.Equal(t, test.expected, resp.Users)
			assert.Equal(t, test.notFound, resp.NotFoundIds)
		} else {
			assert.NotNil(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, test.code, status.Code(err))
		}
	}
}

func TestGetAllUsers(t *testing.T) {
	s := newServer()

	req := &pb.GetAllUsersRequest{}
	resp, err := s.GetAllUsers(context.Background(), req)

	assert.Nil(t, err)
	assert.Equal(t, len(s.users), len(resp.Users))

	for _, user := range resp.Users {
		expectedUser, exists := s.users[user.Id]
		assert.True(t, exists)
		assert.Equal(t, expectedUser, user)
	}
}

func TestSearchUsers(t *testing.T) {
	s := newServer()

	tests := []struct {
		query    string
		expected []*pb.User
		code     codes.Code
	}{
		{"Steve", []*pb.User{s.users[1]}, codes.OK},
		{"LA", []*pb.User{s.users[1], s.users[3]}, codes.OK},
		{"true", []*pb.User{s.users[1], s.users[3]}, codes.OK},
		{"false", []*pb.User{s.users[2]}, codes.OK},
		{"unknown", nil, codes.NotFound},
		{"", nil, codes.InvalidArgument},
	}

	for _, test := range tests {
		req := &pb.SearchUserRequest{Query: test.query}
		resp, err := s.SearchUsers(context.Background(), req)

		if test.code == codes.OK {
			assert.Nil(t, err)
			assert.Equal(t, test.expected, resp.Users)
		} else {
			assert.NotNil(t, err)
			assert.Nil(t, resp)
			assert.Equal(t, test.code, status.Code(err))
		}
	}
}
