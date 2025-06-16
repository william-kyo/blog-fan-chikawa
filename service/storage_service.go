package service

import (
	"blog-fanchiikawa-service/sdk"
	"fmt"
)

// StorageService defines the interface for storage operations
type StorageService interface {
	// GetLastData retrieves the last data from storage
	GetLastData() (string, error)
}

// storageService implements StorageService interface
type storageService struct{}

// NewStorageService creates a new StorageService instance
func NewStorageService() StorageService {
	return &storageService{}
}

// GetLastData retrieves the last data from storage
func (s *storageService) GetLastData() (string, error) {
	data, err := sdk.GetLastData()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve last data: %w", err)
	}
	return data, nil
}