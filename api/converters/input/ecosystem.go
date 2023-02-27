package input

import (
	"arbuga/backend/api/graph/model"
	"arbuga/backend/app"
)

func ConvertEcosystemInput(input model.EcosystemInput) app.EcosystemInput {
	var aquarium *app.AquariumGlassInput
	if inputAquarium := input.Aquarium; inputAquarium != nil {
		var dimensions *app.DimensionsInput

		if inputDimensions := inputAquarium.Dimensions; inputDimensions != nil {
			dimensions = &app.DimensionsInput{
				Width:  inputDimensions.Width,
				Height: inputDimensions.Height,
				Length: inputDimensions.Length,
			}
		}

		aquarium = &app.AquariumGlassInput{
			Dimensions:         dimensions,
			GlassThickness:     inputAquarium.GlassThickness,
			SubstrateThickness: inputAquarium.SubstrateThickness,
			DecorationsVolume:  inputAquarium.DecorationsVolume,
		}
	}
	return app.EcosystemInput{
		Name:     input.Name,
		Aquarium: aquarium,
	}
}
