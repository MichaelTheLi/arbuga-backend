package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	model2 "arbuga/backend/api/graph/model"
	"arbuga/backend/auth"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(_ context.Context, login string, password string) (*model2.LoginResult, error) {
	randValue, _ := rand.Int(rand.Reader, big.NewInt(100))

	user, _ := r.UsersState.GetUserByLogin(login)

	if user == nil {
		// hash and salt the password
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		hashedPassString := string(hashedPass)

		if err != nil {
			return nil, err
		}
		user = &model2.User{
			ID:         fmt.Sprintf("T%d", randValue),
			Name:       "Michael" + fmt.Sprintf("T%d", randValue), // TODO Username
			Login:      &login,
			Password:   &hashedPassString,
			Ecosystems: []*model2.Ecosystem{},
		}
		r.UsersState.Users[user.ID] = user
		log.Println(r.UsersState.Users)
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password))
		if err != nil {
			return nil, errors.New("error")
		}
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &model2.LoginResult{User: user, Token: &token}, nil
}

// SaveEcosystem is the resolver for the saveEcosystem field.
func (r *mutationResolver) SaveEcosystem(ctx context.Context, id *string, ecosystem model2.EcosystemInput) (*model2.EcosystemUpdateResult, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, errors.New("not authenticated")
	}

	user, err := r.UsersState.GetUserByID(user.ID)

	if user == nil || err != nil {
		return nil, errors.New("not authenticated")
	}

	var ecosystemFound *model2.Ecosystem = nil

	result := model2.EcosystemUpdateResult{
		Success:   true,
		Error:     nil,
		Ecosystem: nil,
	}

	if id != nil {
		for _, v := range user.Ecosystems {
			if v.ID == *id {
				ecosystemFound = v
				break
			}
		}

		if ecosystemFound == nil {
			result.Success = false
			s := "Ecosystem not found"
			result.Error = &s
		}
	} else {
		randValue, _ := rand.Int(rand.Reader, big.NewInt(100))
		newId := fmt.Sprintf("T%d", randValue)
		ecosystemFound = &model2.Ecosystem{
			ID: newId,
			Aquarium: &model2.AquariumGlass{
				Dimensions: &model2.Dimensions{
					Width:  0,
					Height: 0,
					Length: 0,
				},
				GlassThickness: 0,
			},
		}
		user.Ecosystems = append(user.Ecosystems, ecosystemFound)
	}

	if ecosystemFound != nil {
		// TODO More sophisticated propagation of the data required. Also, more fields
		if ecosystem.Name != nil {
			ecosystemFound.Name = *ecosystem.Name
		}
		if aquarium := ecosystem.Aquarium; aquarium != nil {
			if dimensions := aquarium.Dimensions; dimensions != nil {
				if dimensions.Width != nil {
					ecosystemFound.Aquarium.Dimensions.Width = *dimensions.Width
				}
				if dimensions.Height != nil {
					ecosystemFound.Aquarium.Dimensions.Height = *dimensions.Height
				}
				if dimensions.Length != nil {
					ecosystemFound.Aquarium.Dimensions.Length = *dimensions.Length
				}
			}
		}

		ecosystemFound.Analysis = calculateAnalysis()
		result.Ecosystem = ecosystemFound
	}

	return &result, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func calculateAnalysis() []*model2.EcosystemAnalysisCategory {
	return []*model2.EcosystemAnalysisCategory{
		{
			ID:          "1",
			Name:        "Filtration",
			Description: "How good filtration is",
			Status:      "ok",
			Messages:    nil,
		},
		{
			ID:          "2",
			Name:        "Temperature",
			Description: "How good filtration is",
			Status:      "moderate",
			Messages: []*model2.EcosystemAnalysisMessage{
				{
					ID:          "1",
					Name:        "Low temperature: fish",
					Description: "Temperature too low for some of the fish",
					Status:      "bad",
				},
				{
					ID:          "2",
					Name:        "Not optimal temperature: fish",
					Description: "Temperature on the low edge for some of the fish",
					Status:      "moderate",
				},
			},
		},
		{
			ID:          "3",
			Name:        "Fish compatibility",
			Description: "Are your fish compatible with each other",
			Status:      "bad",
			Messages: []*model2.EcosystemAnalysisMessage{
				{
					ID:          "1",
					Name:        "Barbus",
					Description: "These guys are likely to behave aggressive towards most of your fish",
					Status:      "bad",
				},
				{
					ID:          "2",
					Name:        "Bolivian Ram",
					Description: "Bolivian Ram are too shy for the neighbourhood with barbuses, you might notice some issues with the food",
					Status:      "moderate",
				},
			},
		},
		{
			ID:          "3",
			Name:        "Plants compatibility",
			Description: "Are your plants compatible with the water in your tank",
			Status:      "moderate",
			Messages: []*model2.EcosystemAnalysisMessage{
				{
					ID:          "1",
					Name:        "Cryptocoryne might melt",
					Description: "Cryptocoryne might melt, water too soft",
					Status:      "moderate",
				},
			},
		},
	}
}