package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kitamersion/go-goservice/internal/domain/entities"
	"github.com/kitamersion/go-goservice/internal/domain/repositories"
	"github.com/kitamersion/go-goservice/internal/events/producer"
	"github.com/kitamersion/go-goservice/internal/events/proto/events/userpb"
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

func (s *UserService) CreateUser(ctx context.Context, user *entities.UserEntity) error {
	if err := s.userRepo.Create(user); err != nil {
		return err
	}
	event := &userpb.UserCreated{
		Id:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Unix(),
	}

	return s.producer.PublishEvent(ctx, event)
}

func (s *UserService) GetUserByID(id uuid.UUID) (*entities.UserEntity, error) {
	return s.userRepo.GetByID(id)
}

func (s *UserService) UpdateUser(ctx context.Context, user *entities.UserEntity) error {
	if err := s.userRepo.Update(user); err != nil {
		return err
	}
	event := &userpb.UserUpdated{
		Id:        user.ID.String(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}

	return s.producer.PublishEvent(ctx, event)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := s.userRepo.Delete(id); err != nil {
		return err
	}
	event := &userpb.UserDeleted{
		Id:        id.String(),
		DeletedAt: time.Now().Unix(),
	}

	return s.producer.PublishEvent(ctx, event)
}
