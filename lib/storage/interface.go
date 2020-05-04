package storage

import "github.com/volatiletech/authboss"

type Storage interface {
	authboss.CreatingServerStorer
	AddInteraction(firstUserID, secondUserID int64, at int64) error
	InteractedWithInfected(userID int64) (bool, error)
	GetInfectedList() ([]int64, error)
	AddInfected(userID int64) error
}
