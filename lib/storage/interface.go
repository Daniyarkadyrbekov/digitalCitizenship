package storage

import "github.com/volatiletech/authboss"

type Storage interface {
	authboss.CreatingServerStorer
	AddInteraction(firstUserID, secondUserID string, at int64) error
	InteractedWithInfected(userID string) (bool, error)
	GetInfectedList() ([]string, error)
	AddInfected(userID string) error
}
