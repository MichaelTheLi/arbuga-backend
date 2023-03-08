package domain

type Fish struct {
	Name        string
	Description string
}

func NewFish(name string, description string) *Fish {
	return &Fish{
		Name:        name,
		Description: description,
	}
}
