package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kitamersion/go-goservice/graph/model"
	"github.com/kitamersion/go-goservice/internal/domain/entities"
)

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (string, error) {
	entity := &entities.UserEntity{
		Name:  input.Name,
		Email: input.Email,
	}

	id, err := r.UserService.CreateUser(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return id.String(), nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID")
	}
	user, err := r.UserService.GetUserByID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID")
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return &model.User{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
