package drivers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Unfield/FileHopper/data"
	_ "modernc.org/sqlite"
)

type SqliteDriver struct {
	db *sql.DB
}

func NewSqliteDriver() (*SqliteDriver, error) {
	return &SqliteDriver{}, nil
}

func (s *SqliteDriver) Init(dsn string) error {
	var err error
	s.db, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed to open sqlite db: %w", err)
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			roles TEXT NOT NULL,
			home_dir TEXT
		);

		CREATE TABLE IF NOT EXISTS quota (
		    id TEXT PRIMARY KEY,
		    user_id TEXT NOT NULL,
		    type TEXT NOT NULL,
		    max INTEGER NOT NULL,
		    used INTEGER NOT NULL
		);
	`)

	return err
}

func (s *SqliteDriver) Close() error {
	return s.db.Close()
}

func (s *SqliteDriver) CreateUser(u data.User) error {
	_, err := s.db.Exec("INSERT INTO users (id, username, password_hash, roles, home_dir) VALUES (?,?,?,?,?)",
		u.ID, u.Username, u.HashedPassword, strings.Join(u.Roles, ","), u.HomeDir,
	)
	return err
}

func (s *SqliteDriver) GetUser(username string) (*data.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE username = ?", username)
	var uid, uname, password_hash, roles, home_dir string
	if err := row.Scan(&uid, &uname, &password_hash, &roles, &home_dir); err != nil {
		return nil, err
	}
	return &data.User{
		ID:             uid,
		Username:       uname,
		HashedPassword: password_hash,
		Roles:          strings.Split(roles, ","),
		HomeDir:        home_dir,
	}, nil
}

func (s *SqliteDriver) ListUsers() ([]data.User, error) {
	rows, err := s.db.Query("SELECT id, username, password_hash, roles, home_dir FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []data.User
	for rows.Next() {
		var uid, uname, password_hash, roles, home_dir string
		if err := rows.Scan(&uid, &uname, &password_hash, &roles, &home_dir); err != nil {
			return nil, err
		}
		users = append(users, data.User{
			ID:             uid,
			Username:       uname,
			HashedPassword: password_hash,
			Roles:          strings.Split(roles, ","),
			HomeDir:        home_dir,
		})
	}

	return users, rows.Err()
}

func (s *SqliteDriver) UpdateUser(u data.User) error {
	_, err := s.db.Exec(
		"UPDATE users SET password_hash = ?, roles = ?, home_dir = ? WHERE username = ?",
		u.HashedPassword, strings.Join(u.Roles, ","), u.HomeDir, u.Username,
	)
	return err
}

func (s *SqliteDriver) DeleteUser(username string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE username = ?", username)
	return err
}

func (s *SqliteDriver) CreateQuota(q data.Quota) error {
	_, err := s.db.Exec("INSERT INTO quota (id, user_id, type, max, used) VALUES (?,?,?,?,?)",
		q.ID, q.UserID, q.Type, q.Max, q.Used,
	)
	return err
}

func (s *SqliteDriver) GetQuota(userId string) (*data.Quota, error) {
	row := s.db.QueryRow("SELECT id, user_id, type, max, used FROM quota WHERE user_id = ?", userId)
	var id, _userId, _type string
	var max, used int64
	if err := row.Scan(&id, &_userId, &_type, &max, &used); err != nil {
		return nil, err
	}
	return &data.Quota{
		ID:     id,
		UserID: _userId,
		Type:   data.LimitType(_type),
		Max:    max,
		Used:   max,
	}, nil
}

func (s *SqliteDriver) UpdateQuota(q data.Quota) error {
	_, err := s.db.Exec(
		"UPDATE quota SET type = ?, max = ?, used = ? WHERE id = ?",
		q.Type, q.Max, q.Used, q.ID,
	)
	return err
}

func (s *SqliteDriver) DeleteQuota(userId string) error {
	_, err := s.db.Exec("DELETE * FROM quota WHERE user_id = ?", userId)
	return err
}
