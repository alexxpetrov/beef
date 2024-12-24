package user

import (
	erdtreev1 "app-websocket/gen/erdtree/v1"
	"app-websocket/gen/erdtree/v1/erdtreev1connect"
	userInfov1 "app-websocket/gen/user/v1"
	"context"
	"fmt"

	"connectrpc.com/connect"
)

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserInfoServer struct {
	cacheClient erdtreev1connect.ErdtreeStoreClient
}

func NewExternalClient(erdtreeClient erdtreev1connect.ErdtreeStoreClient) *UserInfoServer {
	return &UserInfoServer{
		cacheClient: erdtreeClient,
	}
}

func (s *UserInfoServer) GetUserInfo(
	ctx context.Context,
	req *connect.Request[userInfov1.GetInfoRequest],
) (*connect.Response[userInfov1.GetInfoResponse], error) {
	var err error

	loginCacheRequest := connect.NewRequest(&erdtreev1.GetRequest{
		Key: req.Msg.UserId + "-login",
	})

	loginTimestamp, err := s.cacheClient.Get(ctx, loginCacheRequest)
	if err != nil {
		loginTimestamp = connect.NewResponse(&erdtreev1.GetResponse{
			Value: []byte(""),
		})
		fmt.Println("Error retrieving login record in Erdtree:", err)
	}

	registerCacheRequest := connect.NewRequest(&erdtreev1.GetRequest{
		Key: req.Msg.UserId + "-register",
	})

	registerTimestamp, err := s.cacheClient.Get(ctx, registerCacheRequest)
	if err != nil {
		registerTimestamp = connect.NewResponse(&erdtreev1.GetResponse{
			Value: []byte(""),
		})
		fmt.Println("Error retrieving register record in Erdtree:", err)
	}

	joinRoomCacheRequest := connect.NewRequest(&erdtreev1.GetRequest{
		Key: req.Msg.UserId + "-joinedRoom",
	})

	joinRoomId, err := s.cacheClient.Get(ctx, joinRoomCacheRequest)
	if err != nil {
		joinRoomId = connect.NewResponse(&erdtreev1.GetResponse{
			Value: []byte(""),
		})
		fmt.Println("Error retrieving joinedRoom record in Erdtree:", err)
	}

	leftRoomCacheRequest := connect.NewRequest(&erdtreev1.GetRequest{
		Key: req.Msg.UserId + "-leftRoom",
	})

	leftRoomId, err := s.cacheClient.Get(ctx, leftRoomCacheRequest)
	if err != nil {
		leftRoomId = connect.NewResponse(&erdtreev1.GetResponse{
			Value: []byte(""),
		})
		fmt.Println("Error retrieving leftRoom record in Erdtree:", err)
	}

	lastMessageRequest := connect.NewRequest(&erdtreev1.GetRequest{
		Key: req.Msg.UserId + "-lastMessage",
	})

	lastMessage, err := s.cacheClient.Get(ctx, lastMessageRequest)
	if err != nil {
		lastMessage = connect.NewResponse(&erdtreev1.GetResponse{
			Value: []byte(""),
		})
		fmt.Println("Error retrieving leftRoom record in Erdtree:", err)
	}

	lastMessageRoomRequest := connect.NewRequest(&erdtreev1.GetRequest{
		Key: req.Msg.UserId + "-lastMessageRoom",
	})

	lastMessageRoomId, err := s.cacheClient.Get(ctx, lastMessageRoomRequest)
	if err != nil {
		lastMessageRoomId = connect.NewResponse(&erdtreev1.GetResponse{
			Value: []byte(""),
		})
		fmt.Println("Error retrieving leftRoom record in Erdtree:", err)
	}

	userInfoResponse := connect.NewResponse(&userInfov1.GetInfoResponse{
		LoginTimestamp:    string(loginTimestamp.Msg.Value),
		RegisterTimestamp: string(registerTimestamp.Msg.Value),
		JoinedRoomId:      string(joinRoomId.Msg.Value),
		LeftRoomId:        string(leftRoomId.Msg.Value),
		LastMessage:       string(lastMessage.Msg.Value),
		LastMessageRoomId: string(lastMessageRoomId.Msg.Value),
	})

	return userInfoResponse, nil
}
