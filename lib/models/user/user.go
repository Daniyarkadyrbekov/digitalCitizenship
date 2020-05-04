package user

import "errors"

//type User struct {
//	id          int64
//	hashedIIN   string
//	phoneNumber string
//}

type User struct {
	ID       int64
	IINHash  string
	Password string
}

func (u *User) GetPID() string    { return u.IINHash }
func (u *User) PutPID(pid string) { u.IINHash = pid }
func (u *User) GetPassword() (password string) {
	return u.Password
}
func (u *User) PutPassword(password string) {
	u.Password = password
}

func (u *User) Validate() []error {
	return []error{errors.New("validate Err custom")}
}
