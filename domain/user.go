package domain

type Owner struct {
	Name       string
	Ecosystems []*Ecosystem
}

func NewOwner(name string) *Owner {
	return &Owner{
		Name:       name,
		Ecosystems: nil,
	}
}
