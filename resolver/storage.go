package resolver

import (
	"context"
)

// FetchLastData handles the fetchLastData query
func (r *Resolver) FetchLastData(ctx context.Context) (string, error) {
	return r.StorageService.GetLastData()
}