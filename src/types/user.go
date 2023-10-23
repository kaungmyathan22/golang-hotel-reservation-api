package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minFirstNameLength = 2
	minLastNameLength  = 2
	minPasswrodLength  = 7
)

type CreateUserPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (payload CreateUserPayload) Validate() map[string]string {
	errors := map[string]string{}
	if len(payload.FirstName) < minFirstNameLength {
		errors["firstName"] = fmt.Sprintf("firstName length should be at least %d", minFirstNameLength)
	}
	if len(payload.LastName) < minLastNameLength {
		errors["lastName"] = fmt.Sprintf("lastName length should be at least %d", minLastNameLength)
	}
	if len(payload.Password) < minPasswrodLength {
		errors["password"] = fmt.Sprintf("password length should be at least %d", minPasswrodLength)
	}
	if !isEmailValid(payload.Email) {
		errors["email"] = fmt.Sprintf("email is invalid.")
	}
	return errors
}

func isEmailValid(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isMatch, _ := regexp.MatchString(emailRegex, email)
	return isMatch
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstName" json:"firstName"`
	LastName  string             `bson:"lastName" json:"lastName"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
}

func NewUserFromParams(params CreateUserPayload) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  string(hashedPassword),
	}, nil
}
