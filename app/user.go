package app

import (
	"arbuga/backend/domain"
	"crypto/rand"
	"fmt"
	"math/big"
)

type UserService struct {
	Gateway UserGateway
}

type EcosystemUpdateResult struct {
	Ecosystem *domain.Ecosystem
	Success   bool
	Error     *string
}

type EcosystemInput struct {
	Name     *string
	Aquarium *AquariumGlassInput
}

type AquariumGlassInput struct {
	Dimensions         *DimensionsInput
	GlassThickness     *int
	SubstrateThickness *int
	DecorationsVolume  *int
}
type DimensionsInput struct {
	Width  *int
	Height *int
	Length *int
}

func (service *UserService) GetUserById(id string) (*domain.User, error) {
	return service.Gateway.GetUserByID(id)
}

func (service *UserService) SaveEcosystem(user *domain.User, input *EcosystemInput) *EcosystemUpdateResult {
	randValue, _ := rand.Int(rand.Reader, big.NewInt(100))
	newId := fmt.Sprintf("T%d", randValue)
	newEcosystem := &domain.Ecosystem{
		ID: newId,
		Aquarium: &domain.AquariumGlass{
			Dimensions: &domain.Dimensions{
				Width:  0,
				Height: 0,
				Length: 0,
			},
			GlassThickness: 0,
		},
	}
	user.Ecosystems = append(user.Ecosystems, newEcosystem)

	return service.UpdateEcosystem(user, newId, input)
}

func (service *UserService) UpdateEcosystem(user *domain.User, id string, input *EcosystemInput) *EcosystemUpdateResult {
	result := EcosystemUpdateResult{
		Success:   true,
		Error:     nil,
		Ecosystem: nil,
	}

	var ecosystemFound *domain.Ecosystem = nil

	for _, v := range user.Ecosystems {
		if v.ID == id {
			ecosystemFound = v
			break
		}
	}

	if ecosystemFound == nil {
		result.Success = false
		s := "Ecosystem not found"
		result.Error = &s
	} else {
		// TODO More sophisticated propagation of the data required. Also, more fields
		if input.Name != nil {
			ecosystemFound.Name = *input.Name
		}
		if aquarium := input.Aquarium; aquarium != nil {
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

		ecosystemFound.Analysis = ecosystemFound.CalculateAnalysis()
		result.Ecosystem = ecosystemFound
	}

	return &result
}
