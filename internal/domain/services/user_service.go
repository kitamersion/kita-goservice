package services

import (
	"context"
	"errors"
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

func (s *UserService) CreateUser(ctx context.Context, user *entities.UserEntity) (uuid.UUID, error) {
	entity := &entities.UserEntity{
		ID:        uuid.New(),
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.userRepo.Create(entity); err != nil {
		return uuid.UUID{}, err
	}
	event := &userpb.UserCreated{
		Id:        entity.ID.String(),
		Email:     entity.Email,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt.Unix(),
	}

	if err := s.producer.PublishEvent(ctx, event); err != nil {
		return uuid.UUID{}, err
	}

	return entity.ID, nil
}

func (s *UserService) GetUserByID(id uuid.UUID) (*entities.UserEntity, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
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
