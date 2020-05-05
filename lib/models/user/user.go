package user

type User struct {
	IIN      string
	Password string
	Mac      string
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
	return nil
}

func (u *User) GetArbitrary() (arbitrary map[string]string) {
	return map[string]string{
		"mac": u.Mac,
	}
}

func (u *User) PutArbitrary(arbitrary map[string]string) {
	if n, ok := arbitrary["mac"]; ok {
		u.Mac = n
	}
}
