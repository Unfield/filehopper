package drivers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Unfield/FileHopper/internal/auth"
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
			username TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			roles TEXT NOT NULL,
			home_dir TEXT
		);
	`)

	return err
}

func (s *SqliteDriver) Close() error {
	return s.db.Close()
}

func (s *SqliteDriver) CreateUser(u auth.User) error {
	_, err := s.db.Exec("INSERT INTO users (id, username, password_hash, roles, home_dir VALUES (?,?,?,?,?)",
		u.ID, u.Username, u.HashedPassword, strings.Join(u.Roles, ","), u.HomeDir,
	)
	return err
}

func (s *SqliteDriver) GetUser(username string) (*auth.User, error) {
	row := s.db.QueryRow("SELECT * FROM users WHERE username=?", username)
	var uid, uname, password_hash, roles, home_dir string
	if err := row.Scan(&uid, &uname, &password_hash, &roles, &home_dir); err != nil {
		return nil, err
	}
	return &auth.User{
		ID:             uid,
		Username:       uname,
		HashedPassword: password_hash,
		Roles:          strings.Split(roles, ","),
		HomeDir:        home_dir,
	}, nil
}

func (s *SqliteDriver) ListUsers() ([]auth.User, error) {
	rows, err := s.db.Query("SELECT username, password_hash, roles, home_dir FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []auth.User
	for rows.Next() {
		var uid, uname, password_hash, roles, home_dir string
		if err := rows.Scan(&uid, &uname, &password_hash, &roles, &home_dir); err != nil {
			return nil, err
		}
		users = append(users, auth.User{
			ID:             uid,
			Username:       uname,
			HashedPassword: password_hash,
			Roles:          strings.Split(roles, ","),
			HomeDir:        home_dir,
		})
	}
	return users, rows.Err()
}

func (s *SqliteDriver) UpdateUser(u auth.User) error {
	_, err := s.db.Exec(
		"UPDATE users SET password_hash=?, roles=?, home_dir=? WHERE username=?",
		u.HashedPassword, strings.Join(u.Roles, ","), u.HomeDir, u.Username,
	)
	return err
}

func (s *SqliteDriver) DeleteUser(username string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE username=?", username)
	return err
}
