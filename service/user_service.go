package service

import (
	"blog-fanchiikawa-service/db"
	"blog-fanchiikawa-service/graph/model"
	"blog-fanchiikawa-service/greetings"
	"blog-fanchiikawa-service/repository"
	"log"
)

// UserService defines the interface for user business logic
type UserService interface {
	// Login handles user login/registration logic
	Login(nickname, email, deviceID string) (*model.User, error)
	
	// GetUsers retrieves a list of users
	GetUsers(limit int) ([]*model.User, error)
}

// userService implements UserService interface
type userService struct {
	userRepo       repository.UserRepository
	deviceRepo     repository.UserDeviceRepository
	transactionMgr repository.TransactionManager
}

// NewUserService creates a new UserService instance
func NewUserService(
	userRepo repository.UserRepository,
	deviceRepo repository.UserDeviceRepository,
	transactionMgr repository.TransactionManager,
) UserService {
	return &userService{
		userRepo:       userRepo,
		deviceRepo:     deviceRepo,
		transactionMgr: transactionMgr,
	}
}

// Login handles user login/registration logic
func (s *userService) Login(nickname, email, deviceID string) (*model.User, error) {
	// Check if user exists by email
	existingUser, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		// User exists, return the user
		return s.convertToGraphQLUser(existingUser), nil
	}

	// User doesn't exist, create new user and device in transaction
	var newUser *db.User
	err = s.transactionMgr.WithTransaction(func() error {
		// Create new user
		newUser = &db.User{
			Nickname: nickname,
			Email:    email,
		}

		if err := s.userRepo.Create(newUser); err != nil {
			return err
		}

		// Create user device
		userDevice := &db.UserDevice{
			UserID:   newUser.ID,
			DeviceID: deviceID,
		}

		return s.deviceRepo.Create(userDevice)
	})

	if err != nil {
		return nil, err
	}

	// Log greeting message
	message, _ := greetings.Hello(newUser.Nickname)
	log.Println(message)

	return s.convertToGraphQLUser(newUser), nil
}

// GetUsers retrieves a list of users
func (s *userService) GetUsers(limit int) ([]*model.User, error) {
	dbUsers, err := s.userRepo.List(limit, 0)
	if err != nil {
		return nil, err
	}

	// Convert database models to GraphQL models
	var users []*model.User
	for _, dbUser := range dbUsers {
		users = append(users, s.convertToGraphQLUser(dbUser))
	}

	return users, nil
}

// convertToGraphQLUser converts database User model to GraphQL User model
func (s *userService) convertToGraphQLUser(dbUser *db.User) *model.User {
	return &model.User{
		ID:        dbUser.ID,
		Nickname:  dbUser.Nickname,
		Email:     dbUser.Email,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}