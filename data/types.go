package data

type User struct {
	ID             string
	Username       string
	HashedPassword string
	HomeDir        string
	Roles          []string
}
