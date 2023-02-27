package domain

type AquariumGlass struct {
	Dimensions         *Dimensions
	GlassThickness     int
	SubstrateThickness *int
	DecorationsVolume  *int
}

type Dimensions struct {
	Width  int
	Height int
	Length int
}
