package domain

type User struct {
	ID           string
	Login        *string
	PasswordHash *string
	Name         string
	Ecosystems   []*Ecosystem
}
