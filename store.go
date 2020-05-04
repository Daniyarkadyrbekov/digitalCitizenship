package main

import (
	"context"
	"errors"

	"github.com/volatiletech/authboss"
)

type User struct {
	ID       int64
	IINHash  string
	password string
}

func (u *User) GetPID() string    { return u.IINHash }
func (u *User) PutPID(pid string) { u.IINHash = pid }
func (u *User) GetPassword() (password string) {
	return u.password
}
func (u *User) PutPassword(password string) {
	u.password = password
}

func (u *User) Validate() []error {
	return []error{errors.New("validate Err custom")}
}

// MemStorer stores users in memory
type MemStorer struct {
	Users  map[string]User
	Tokens map[string][]string
}

// NewMemStorer constructor
func NewMemStorer() *MemStorer {
	return &MemStorer{
		Users: map[string]User{
			"123456781234": {
				ID:       1,
				IINHash:  "123456781234",
				password: "password",
			},
		},
		Tokens: make(map[string][]string),
	}
}

// Save the user

func (m MemStorer) Save(_ context.Context, user authboss.User) error {
	u := user.(*User)
	m.Users[u.IINHash] = *u

	return nil
}

// Load the user
func (m MemStorer) Load(_ context.Context, key string) (user authboss.User, err error) {

	u, ok := m.Users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &u, nil
}

// New user creation
func (m MemStorer) New(_ context.Context) authboss.User {
	return &User{}
}

// Create the user
func (m MemStorer) Create(_ context.Context, user authboss.User) error {
	u := user.(*User)

	if _, ok := m.Users[u.IINHash]; ok {
		return authboss.ErrUserFound
	}

	m.Users[u.IINHash] = *u
	return nil
}
