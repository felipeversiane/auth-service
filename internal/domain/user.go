package domain

import (
	"time"

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
}

func New(id uuid.UUID, email, password, phone, firstName, lastName string, createdAt time.Time, updatedAt time.Time) UserInterface {
	return &user{
		id:        id,
		email:     email,
		password:  password,
		phone:     phone,
		firstName: firstName,
		lastName:  lastName,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
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
