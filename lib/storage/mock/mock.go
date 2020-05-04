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
	interactions map[string]*models.Interactions
	infectedList []string
}

// NewMemStorer constructor
func New() *Mock {
	return &Mock{
		users: map[string]user.User{
			"123456781234": {
				IIN:      "123456781234",
				Password: "password",
			},
		},
		interactions: map[string]*models.Interactions{},
		infectedList: []string{},
	}
}

// Save the user
func (m *Mock) Save(_ context.Context, bossUser authboss.User) error {
	u := bossUser.(*user.User)
	m.users[u.IIN] = *u

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

	if _, ok := m.users[u.IIN]; ok {
		return authboss.ErrUserFound
	}

	m.users[u.IIN] = *u
	return nil
}

func (m *Mock) AddInteraction(firstUserIIN, secondUserIIN string, at int64) error {

	if _, ok := m.interactions[firstUserIIN]; !ok {
		m.interactions[firstUserIIN] = models.NewInteractions()
	}
	m.interactions[firstUserIIN].Add(secondUserIIN, at)

	return nil
}
func (m *Mock) InteractedWithInfected(firstUserIIN string) (bool, error) {

	if _, ok := m.interactions[firstUserIIN]; !ok {
		return false, nil
	}

	for _, infectedIIN := range m.infectedList {
		if m.interactions[firstUserIIN].Search(infectedIIN) {
			return true, nil
		}
	}

	return false, nil
}
func (m *Mock) GetInfectedList() ([]string, error) {
	return m.infectedList, nil
}
func (m *Mock) AddInfected(userIIN string) error {
	m.infectedList = append(m.infectedList, userIIN)
	return nil
}
