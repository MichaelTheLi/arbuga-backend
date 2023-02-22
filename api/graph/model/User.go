package model

type User struct {
	ID         string       `json:"id"`
	Login      *string      `json:"login"`
	Password   *string      `json:"-"`
	Name       string       `json:"name"`
	Ecosystems []*Ecosystem `json:"ecosystems"`
	Ecosystem  *Ecosystem   `json:"ecosystem"`
}
