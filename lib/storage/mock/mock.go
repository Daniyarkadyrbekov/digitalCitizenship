package mock

import (
	"github.com/digitalCitizenship/lib/storage"
)

type Mock struct {
	users []storage.AuthInfo
}

func New() Mock {
	return Mock{users: []storage.AuthInfo{}}
}

func (m Mock) AuthSave(userID int64) (info storage.AuthInfo, err error) {
	return
}

func (m Mock) AuthRemove(userID int64) (info storage.AuthInfo, err error) {
	return
}

func (m Mock) AuthCheck(token string) (userID int64, err error) {
	return
}
