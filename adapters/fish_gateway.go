package adapters

import (
	"arbuga/backend/app"
	"errors"
	"sort"
	"strings"
)

type LocalFishGateway struct {
	Fish map[string]*app.Fish
}

func (gateway *LocalFishGateway) GetFishById(id string) (*app.Fish, error) {
	fish, success := gateway.Fish[id]
	if success != true {
		return nil, errors.New("fish not found")
	}
	return fish, nil
}

func (gateway *LocalFishGateway) SearchFishBySubstring(substring string) ([]*app.Fish, error) {
	var fishListFound []*app.Fish

	for _, fish := range gateway.Fish {
		if strings.Contains(strings.ToLower(fish.Fish.Name), strings.ToLower(substring)) {
			fishListFound = append(fishListFound, fish)
		}
	}

	sort.Slice(fishListFound, func(i, j int) bool {
		return fishListFound[i].Id < fishListFound[j].Id
	})

	return fishListFound, nil
}

func (gateway *LocalFishGateway) SaveFish(fish *app.Fish) (*app.Fish, error) {
	gateway.Fish[fish.Id] = fish
	return fish, nil
}
