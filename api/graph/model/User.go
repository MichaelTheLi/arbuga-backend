package model

type User struct {
	ID         string       `json:"id"`
	Login      *string      `json:"login"`
	Name       string       `json:"name"`
	Ecosystems []*Ecosystem `json:"ecosystems"`
	Ecosystem  *Ecosystem   `json:"ecosystem"`
}
