package resolver

import (
	"blog-fanchiikawa-service/graph/model"
	"context"
)

// Users handles the users query
func (r *Resolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.UserService.GetUsers(10) // Default limit of 10
}

// FetchLastData handles the fetchLastData query
func (r *Resolver) FetchLastData(ctx context.Context) (string, error) {
	return r.StorageService.GetLastData()
}