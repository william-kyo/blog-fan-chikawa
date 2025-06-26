package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"context"
)

// Login handles the login mutation
func (r *Resolver) Login(ctx context.Context, input model.LoginUser) (*model.User, error) {
	return r.UserService.Login(input.Nickname, input.Email, input.DeviceID)
}

// Users handles the users query
func (r *Resolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.UserService.GetUsers(10) // Default limit of 10
}
