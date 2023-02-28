package app

type UserGateway interface {
	GetUserByLogin(login string) (*User, error)
	GetUserByID(id string) (*User, error)
	SaveUser(*User) (*User, error)
}

type AuthService interface {
	HashFromPassword(password string) (string, error)
	IsHashValid(expectedHash string, password string) (bool, error)
}

type TokenService interface {
	GenerateToken(user *User) (string, error)
	GetUserIdFromToken(token string) (string, error)
}
