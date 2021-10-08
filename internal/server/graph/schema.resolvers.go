package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"go_server/internal/server/graph/generated"
	"go_server/internal/server/graph/model"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	query := map[string]interface{}{}

	dbUsers, err := r.store.ListUsers(query)

	if err != nil {
		return nil, err
	}

	users := []*model.User{}
	for _, user := range dbUsers {
		users = append(
			users,
			&model.User{
				FirstName: *user.FirstName,
				LastName:  *user.LastName,
				Auth0ID:  *user.Auth0ID,
				UserID:    user.UserID.String(),
				CreatedAt:    user.CreatedAt.String(),
				UpdatedAt:    user.UpdatedAt.String(),
			})
	}

	return users, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
