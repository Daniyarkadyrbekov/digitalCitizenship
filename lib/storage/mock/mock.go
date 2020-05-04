package mock

import (
	"context"

	"github.com/digitalCitizenship/lib/models"

	"github.com/digitalCitizenship/lib/models/user"

	"github.com/volatiletech/authboss"
)

// Mock stores users in memory
type Mock struct {
	users        map[string]user.User
	interactions map[int64]models.Interactions
	infectedList []int64
}

// NewMemStorer constructor
func New() *Mock {
	return &Mock{
		users: map[string]user.User{
			"123456781234": {
				ID:       1,
				IINHash:  "123456781234",
				Password: "password",
			},
		},
		interactions: map[int64]models.Interactions{},
		infectedList: []int64{},
	}
}

// Save the user
func (m *Mock) Save(_ context.Context, bossUser authboss.User) error {
	u := bossUser.(*user.User)
	m.users[u.IINHash] = *u

	return nil
}

// Load the user
func (m *Mock) Load(_ context.Context, key string) (user authboss.User, err error) {

	u, ok := m.users[key]
	if !ok {
		return nil, authboss.ErrUserNotFound
	}

	return &u, nil
}

// New user creation
func (m *Mock) New(_ context.Context) authboss.User {
	return &user.User{}
}

// Create the user
func (m *Mock) Create(_ context.Context, bossUser authboss.User) error {
	u := bossUser.(*user.User)

	if _, ok := m.users[u.IINHash]; ok {
		return authboss.ErrUserFound
	}

	m.users[u.IINHash] = *u
	return nil
}

func (m *Mock) AddInteraction(firstUserID, secondUserID int64, at int64) error {

	if _, ok := m.interactions[firstUserID]; !ok {
		m.interactions[firstUserID] = models.NewInteractions()
	}
	m.interactions[firstUserID].Add(secondUserID, at)

	return nil
}
func (m *Mock) InteractedWithInfected(firstUserID int64) (bool, error) {

	if _, ok := m.interactions[firstUserID]; !ok {
		return false, nil
	}

	for _, infectedID := range m.infectedList {
		if m.interactions[firstUserID].Search(infectedID) {
			return true, nil
		}
	}

	return false, nil
}
func (m *Mock) GetInfectedList() ([]int64, error) {
	return m.infectedList, nil
}
func (m *Mock) AddInfected(userID int64) error {
	m.infectedList = append(m.infectedList, userID)
	return nil
}
