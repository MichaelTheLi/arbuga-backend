package adapters

import (
	"arbuga/backend/app"
	"strings"
)

type LocalFishGateway struct {
	Fish map[string]*app.Fish
}

func (gateway *LocalFishGateway) GetFishById(id string) (*app.Fish, error) {
	fish, _ := gateway.Fish[id]
	return fish, nil
}

func (gateway *LocalFishGateway) SearchFishBySubstring(substring string) ([]*app.Fish, error) {
	var fishListFound []*app.Fish

	for _, fish := range gateway.Fish {
		if strings.Contains(strings.ToLower(fish.Fish.Name), strings.ToLower(substring)) {
			fishListFound = append(fishListFound, fish)
		}
	}

	return fishListFound, nil
}

func (gateway *LocalFishGateway) SaveFish(fish *app.Fish) (*app.Fish, error) {
	gateway.Fish[fish.Id] = fish
	return fish, nil
}
