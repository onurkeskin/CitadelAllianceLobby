package domain

type IUser interface {
	GetID() string
	IsValid() bool
	IsCodeVerified(code string) bool
	IsCredentialsVerified(password string) bool
	SetPassword(password string) error
}

type IUsers interface{}
