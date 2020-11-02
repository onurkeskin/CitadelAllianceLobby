package service

import (
	"github.com/dgrijalva/jwt-go"

	domain "keon.com/CitadelAllianceLobbyServer/user-service/dbclient/domain"
	userModel "keon.com/CitadelAllianceLobbyServer/user-service/model"
)

var (
	// Define a secure key string used
	// as a salt when hashing our tokens.
	// Please make your own way more secure than this,
	// use a randomly generated md5 hash or something.
	key = []byte("85429ec091e53154f7bfaab36c30d7fc")
)

// CustomClaims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type CustomClaims struct {
	User *userModel.User
	jwt.StandardClaims
}

type Authable interface {
	Decode(token string) (*CustomClaims, error)
	Encode(user *userModel.User) (string, error)
}

type TokenService struct {
	Repo domain.IUserRepository
}

// Decode a token string into a token object
func (srv *TokenService) Decode(token string) (*CustomClaims, error) {

	// Parse the token
	tokenType, err := jwt.ParseWithClaims(string(key), &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// Validate the token and return the custom claims
	if claims, ok := tokenType.Claims.(*CustomClaims); ok && tokenType.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

// Encode a claim into a JWT
func (srv *TokenService) Encode(user *userModel.User) (string, error) {
	// Create the Claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "keon.citadelalliance.user",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(key)
}
