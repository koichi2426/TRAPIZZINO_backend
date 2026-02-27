package entities

import (
	"src/domain/value_objects"
)

type User struct {
	ID             value_objects.ID
	Username       string
	Email          value_objects.Email
	HashedPassword string
}

func NewUser(id int, username, email, hashedPassword string) (*User, error) {
	uid, err := value_objects.NewID(id)
	if err != nil {
		return nil, err
	}
	emailVO, err := value_objects.NewEmail(email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:             uid,
		Username:       username,
		Email:          emailVO,
		HashedPassword: hashedPassword,
	}, nil
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindByID(id value_objects.ID) (*User, error)
	FindByEmail(email value_objects.Email) (*User, error)
	Update(user *User) error
	Delete(id value_objects.ID) error
}
