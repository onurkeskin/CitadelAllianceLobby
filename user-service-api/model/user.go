package model

import (
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Roles []string

type Users []*User
type User struct {
	// The id of the user.
	//
	// required: true
	ID string `json:"id,omitempty" db:"id"`
	// The username of the user.
	//
	// required: true
	Username string `json:"username,omitempty" db:"username"`
	// The email of the user.
	//
	// required: false
	Email string `json:"email,omitempty" db:"email"`
	Roles Roles  `json:"roles,omitempty" db:"roles"`
	// The status of the user.
	//
	// required: true
	Status           string    `json:"status,omitempty" db:"status"`
	LastModifiedDate time.Time `json:"lastModifiedDate" db:"lastModifiedDate"`
	CreatedDate      time.Time `json:"createdDate,omitempty" db:"createdDate"`

	// fields are not exported to JSON
	ConfirmationCode string `json:"-" db:"conformationCode"`
	HashedPassword   string `json:"-" db:"hashedPassword"`
}

// GetID returns user id
func (user *User) GetID() string {
	return user.ID
}

// IsValid Ensures that the customer object is valid
func (user *User) IsValid() bool {
	// Note regarding email address validation,
	// as long as it `*looks* like an address, we'll allow it.
	// User need to confirm account by click on link sent to the email anyways
	return len(user.Username) > 0 &&
		len(user.Email) > 0 &&
		strings.Contains(user.Email, "@")
}

// IsCodeVerified verify the given code
func (user *User) IsCodeVerified(code string) bool {
	return (user.ConfirmationCode == code)
}

// IsCredentialsVerified verify the given credentials
func (user *User) IsCredentialsVerified(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return (err == nil)
}

// SetPassword encrypts the given plain text password
func (user *User) SetPassword(password string) error {
	passwordBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	return nil
}

// SetPassword encrypts the given plain text password
func (user *User) GenerateConfirmationCode() {
	user.ConfirmationCode = generateNewUniqueCode()
}

// generateNewUniqueCode generates a new confirmation code
func generateNewUniqueCode() string {
	// set code format
	//uuid.SwitchFormat(uuid.Clean)
	id, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return id.String()
}
