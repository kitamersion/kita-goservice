package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/kitamersion/go-goservice/internal/domain/entities"
	"github.com/kitamersion/go-goservice/internal/domain/repositories"
	"github.com/kitamersion/go-goservice/internal/events/producer"
	"github.com/kitamersion/go-goservice/internal/events/types"
)

type UserService struct {
	userRepo repositories.UserRepository
	producer *producer.Producer
}

func NewUserService(userRepo repositories.UserRepository, producer *producer.Producer) *UserService {
	return &UserService{
		userRepo: userRepo,
		producer: producer,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *entities.User) error {
	if err := s.userRepo.Create(user); err != nil {
		return err
	}

	// Publish user created event
	event := types.NewEvent(types.UserCreated, types.UserCreatedEvent{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	})

	return s.producer.PublishEvent(ctx, event)
}

func (s *UserService) GetUserByID(id uuid.UUID) (*entities.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateUser(ctx context.Context, user *entities.User) error {
	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Publish user updated event
	event := types.NewEvent(types.UserUpdated, types.UserUpdatedEvent{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
	})

	return s.producer.PublishEvent(ctx, event)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(id); err != nil {
		return err
	}

	// Publish user deleted event
	event := types.NewEvent(types.UserDeleted, types.UserDeletedEvent{
		UserID: id,
	})

	return s.producer.PublishEvent(ctx, event)
}
