package permissions

import (
	"context"
	"doki.co.in/doki_user_service/user"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type PermissionServerImpl struct {
	Driver *neo4j.DriverWithContext
	UnimplementedPermissionServer
}

func (permissionServer *PermissionServerImpl) CanSendMessage(_ context.Context, req *CanSendMessageRequest) (*PermissionResponse, error) {
	userModel := user.UserModel{
		Driver: permissionServer.Driver,
	}

	// when block feature is implemented
	// use a single db query to get the result
	isBlock, err := userModel.IsUserBlocked(req.Username, req.RecipientUsername)
	if err != nil {
		return nil, err
	}

	isFriend, err := userModel.IsUserFriend(req.Username, req.RecipientUsername)
	if err != nil {
		return nil, err
	}

	return &PermissionResponse{CanPerformAction: !isBlock && isFriend}, nil
}