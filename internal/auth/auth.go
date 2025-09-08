package auth

import (
	"github.com/Unfield/FileHopper/data"
	"github.com/Unfield/FileHopper/internal/db"
	"github.com/Unfield/FileHopper/utils"
)

type Authenticator struct {
	db db.DBDriver
}

func NewAuthenticator(dbDriver db.DBDriver) (*Authenticator, error) {
	return &Authenticator{
		db: dbDriver,
	}, nil
}

func (a *Authenticator) CreateUser(username, password, homedir string, roles []string) (*data.User, error) {
	userId, err := utils.GenerateNanoid()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	var user = data.User{
		ID:             userId,
		Username:       username,
		HashedPassword: hashedPassword,
		HomeDir:        homedir,
		Roles:          roles,
	}

	err = a.db.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *Authenticator) Authenticate(username, password string) bool {
	user, err := a.db.GetUser(username)
	if err != nil {
		return false
	}

	return utils.ComparePassword(password, user.HashedPassword)
}

func (a *Authenticator) ListUsers() []data.User {
	users, err := a.db.ListUsers()
	if err != nil || len(users) <= 0 {
		return []data.User{}
	}
	return users
}
