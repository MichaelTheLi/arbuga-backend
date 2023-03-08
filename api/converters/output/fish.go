package output

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/app"
)

func ConvertFish(domainModel *app.Fish) *model.Fish {
	return &model.Fish{
		ID:          domainModel.Id,
		Name:        domainModel.Fish.Name,
		Description: domainModel.Fish.Description,
	}
}

func ConvertFishList(domainList []*app.Fish) []*model.Fish {
	var fishList []*model.Fish
	fishList = []*model.Fish{}

	for _, fish := range domainList {
		modelFish := ConvertFish(fish)
		fishList = append(fishList, modelFish)
	}

	return fishList
}
