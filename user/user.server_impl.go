package user

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserServerImpl struct {
	Driver *neo4j.DriverWithContext
	UnimplementedUserServer
}

func (userServer *UserServerImpl) GetUserInfo(ctx context.Context, req *UserInfoRequest) (*UserInfoResponse, error) {
	userModel := UserModel{
		Driver: userServer.Driver,
	}

	users, err := userModel.getUserDetails(req.Username, req.UsersToRequest)
	if err != nil {
		return nil, err
	}

	return &UserInfoResponse{
		Users: users,
	}, nil
}
