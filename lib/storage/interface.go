package storage

type Storage interface {
	AuthSave(userID int64) (AuthInfo, error)
	AuthRemove(userID int64) (AuthInfo, error)
	AuthCheck(token string) (int64, error)
}
