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
	Ecosystem *Ecosystem
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

type User struct {
	ID           string
	Owner        *domain.Owner
	Login        *string
	PasswordHash *string
	Ecosystems   []*Ecosystem
}

type Ecosystem struct {
	ID        string
	Ecosystem *domain.Ecosystem
	Fish      []*Fish
}

func (service *UserService) GetUserById(id string) (*User, error) {
	return service.Gateway.GetUserByID(id)
}

func (service *UserService) SaveEcosystem(user *User, input *EcosystemInput) *EcosystemUpdateResult {
	randValue, _ := rand.Int(rand.Reader, big.NewInt(100))
	newId := fmt.Sprintf("T%d", randValue)
	newEcosystem := &Ecosystem{
		ID: newId,
		Ecosystem: &domain.Ecosystem{
			Aquarium: &domain.AquariumGlass{
				Dimensions: &domain.Dimensions{
					Width:  0,
					Height: 0,
					Length: 0,
				},
				GlassThickness: 0,
			},
		},
	}
	// TODO Hide in app.User? Kind of aggregate?
	user.Ecosystems = append(user.Ecosystems, newEcosystem)
	user.Owner.Ecosystems = append(user.Owner.Ecosystems, newEcosystem.Ecosystem)

	return service.UpdateEcosystem(user, newId, input)
}

func (service *UserService) UpdateEcosystem(user *User, id string, input *EcosystemInput) *EcosystemUpdateResult {
	result := EcosystemUpdateResult{
		Success:   true,
		Error:     nil,
		Ecosystem: nil,
	}

	var ecosystemFound *Ecosystem = nil

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
			ecosystemFound.Ecosystem.Name = *input.Name
		}
		if aquarium := input.Aquarium; aquarium != nil {
			if dimensions := aquarium.Dimensions; dimensions != nil {
				if dimensions.Width != nil {
					ecosystemFound.Ecosystem.Aquarium.Dimensions.Width = *dimensions.Width
				}
				if dimensions.Height != nil {
					ecosystemFound.Ecosystem.Aquarium.Dimensions.Height = *dimensions.Height
				}
				if dimensions.Length != nil {
					ecosystemFound.Ecosystem.Aquarium.Dimensions.Length = *dimensions.Length
				}
			}
		}

		ecosystemFound.Ecosystem.Analysis = ecosystemFound.Ecosystem.CalculateAnalysis()

		_, err := service.Gateway.SaveUser(user)

		if err != nil {
			result.Success = false
			s := "Ecosystem not updated"
			result.Error = &s
		}

		result.Ecosystem = ecosystemFound
	}

	return &result
}
