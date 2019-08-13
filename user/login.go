package user

type LoginService interface {
	Login(username string, pwd string) bool
}

type BasicLoginService struct {
	db map[string]string
}

func (b *BasicLoginService) Login(username string, pwd string) bool {
	// encode the password
	epwd := encodePassword(pwd)
	val := b.db[username]

	if val == "" {
		return false
	} else if val == epwd {
		return true
	} else {
		return false
	}
}
