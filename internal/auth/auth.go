package auth

type User struct {
	ID             string
	Username       string
	HashedPassword string
	HomeDir        string
	Roles          []string
}

func CreateUser(username, password, homedir, role string) *User {
	return &User{}
}

func Authenticate(username, password string) bool {
	return true
}

func ListUsers() []User {
	return []User{}
}
