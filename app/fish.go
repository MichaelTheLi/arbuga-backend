package app

import (
	"arbuga/backend/domain"
)

type FishService struct {
	Gateway FishGateway
}

type Fish struct {
	Id   string
	Fish *domain.Fish
}

func (service *FishService) GetFishById(id string) (*Fish, error) {
	return service.Gateway.GetFishById(id)
}

func (service *FishService) SearchFishBySubstring(substring string) ([]*Fish, error) {
	return service.Gateway.SearchFishBySubstring(substring)
}
