package user

import "errors"

type User struct {
	IIN      string
	Password string
}

func (u *User) GetPID() string    { return u.IIN }
func (u *User) PutPID(pid string) { u.IIN = pid }
func (u *User) GetPassword() (password string) {
	return u.Password
}
func (u *User) PutPassword(password string) {
	u.Password = password
}

func (u *User) Validate() []error {
	return []error{errors.New("validate Err custom")}
}
