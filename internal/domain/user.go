package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
)

type user struct {
	id        uuid.UUID
	email     string
	password  string
	phone     string
	firstName string
	lastName  string
	createdAt time.Time
	updatedAt time.Time
}

type UserInterface interface {
	GetID() uuid.UUID
	GetEmail() string
	GetPassword() string
	GetPhone() string
	GetFirstName() string
	GetLastName() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	ComparePassword(password string) bool
}

func New(email, password, phone, firstName, lastName string) UserInterface {
	user := &user{
		id:        uuid.Must(uuid.NewRandom()),
		email:     email,
		password:  hashPassword(password),
		phone:     phone,
		firstName: firstName,
		lastName:  lastName,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}
	return user
}

func (u *user) GetID() uuid.UUID {
	return u.id
}

func (u *user) GetEmail() string {
	return u.email
}

func (u *user) GetPassword() string {
	return u.password
}

func (u *user) GetPhone() string {
	return u.phone
}

func (u *user) GetFirstName() string {
	return u.firstName
}

func (u *user) GetLastName() string {
	return u.lastName
}

func (u *user) GetCreatedAt() time.Time {
	return u.createdAt
}

func (u *user) GetUpdatedAt() time.Time {
	return u.updatedAt
}

func (u *user) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	return err == nil
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}
