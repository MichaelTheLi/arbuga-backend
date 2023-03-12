package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"arbuga/backend/api/converters/output"
	"arbuga/backend/api/graph/model"
	"context"
	"encoding/base64"
	"errors"
)

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	user := ForContext(ctx)
	if user == nil {
		return nil, errors.New("not authenticated")
	}

	userFromService, err := r.UserService.GetUserById(user.ID)

	if err != nil {
		return nil, err
	}

	return output.ConvertUser(userFromService), nil
}

// FishList is the resolver for the fishList field.
//
//goland:noinspection GoUnusedParameter
func (r *queryResolver) FishList(ctx context.Context, substring *string, first *int, after *string) (*model.FishListConnection, error) {
	if substring == nil {
		raw := ""
		substring = &raw
	}

	// The cursor is base64 encoded by convention, so we need to decode it first
	var decodedCursor string
	if after != nil {
		b, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, err
		}
		decodedCursor = string(b)
	}

	// TODO Cursor-aware app service?
	fishList, err := r.FishService.SearchFishBySubstring(*substring)

	if err != nil {
		return nil, err
	}

	count := 0
	currentPage := false

	// If no cursor present start from the top
	if decodedCursor == "" {
		currentPage = true
	}

	hasNextPage := false
	pageInfo := &model.PageInfo{
		StartCursor: "",
		EndCursor:   "",
		HasNextPage: &hasNextPage,
	}

	outputFishList := output.ConvertFishList(fishList)
	edges := make([]*model.FishListEdge, *first)

	for i, fish := range outputFishList {
		if currentPage && count < *first {
			edges[count] = &model.FishListEdge{
				Cursor: base64.StdEncoding.EncodeToString([]byte(fish.ID)),
				Node:   fish,
			}
			count++
		}

		if fish.ID == decodedCursor {
			currentPage = true
		}

		// If there are any elements left after the current page
		// we indicate that in the response
		if count == *first && i < len(outputFishList)-1 {
			*pageInfo.HasNextPage = true
		}
	}

	if count > 0 {
		pageInfo.StartCursor = base64.StdEncoding.EncodeToString([]byte(edges[0].Node.ID))
		pageInfo.EndCursor = base64.StdEncoding.EncodeToString([]byte(edges[count-1].Node.ID))
	}

	return &model.FishListConnection{
		Edges:    edges[:count],
		PageInfo: pageInfo,
	}, nil
}

// Fish is the resolver for the fish field.
//
//goland:noinspection GoUnusedParameter
func (r *queryResolver) Fish(ctx context.Context, id string) (*model.Fish, error) {
	fish, err := r.FishService.GetFishById(id)

	if err != nil {
		return nil, err
	}

	return output.ConvertFish(fish), nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
