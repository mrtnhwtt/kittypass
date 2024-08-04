package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(databasePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("error while openening sqlite database: %s", err)
	}

	storage := &Storage{db: db}

	// Ensure the database has the required table
	if err := storage.initDB(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error while running initialization query: %s", err)
	}

	return storage, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) initDB() error {
	query := `CREATE TABLE IF NOT EXISTS passwords (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL UNIQUE,
        username TEXT NOT NULL,
        encrypted_password TEXT NOT NULL,
        salt TEXT NOT NULL
    );`
	_, err := s.db.Exec(query)
	return err
}

func (s *Storage) Save(name, username, encryptedPassword, salt string) (int64, error) {
	query := `
    INSERT INTO passwords (name, username, encrypted_password, salt)
    VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, name, username, encryptedPassword, salt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Storage) Read(name string) (map[string]string, error) {
	query := `SELECT username, encrypted_password, salt FROM passwords WHERE name = ?`
	row := s.db.QueryRow(query, name)

	var username, encryptedPassword, salt string
	err := row.Scan(&username, &encryptedPassword, &salt)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"name":               name,
		"username":           username,
		"encrypted_password": encryptedPassword,
		"salt":               salt,
	}, nil
}

func (s *Storage) List(name, username string) ([]map[string]string, error) {
	query := `SELECT username, name FROM passwords WHERE name LIKE ? OR username LIKE ?`
	name = "%" + name + "%"
	username = "%" + username + "%"
	rows, err := s.db.Query(query, name, username)
	if err != nil {
		return nil, err
	}
	var loginList []map[string]string
	for rows.Next() {
		var loginName, loginUsername string
		err := rows.Scan(&loginUsername, &loginName)
		if err != nil {
			return nil, err
		}
		loginList = append(loginList, map[string]string{"name": loginName, "username": loginUsername})
	}
	return loginList, nil
}

func (s *Storage) ListAll() ([]map[string]string, error) {
	query := `SELECT username, name FROM passwords`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	var loginList []map[string]string
	for rows.Next() {
		var loginName, loginUsername string
		err := rows.Scan(&loginUsername, &loginName)
		if err != nil {
			return nil, err
		}
		loginList = append(loginList, map[string]string{"name": loginName, "username": loginUsername})
	}
	return loginList, nil
}

func (s *Storage) Update(id int64, name, username, encryptedPassword, salt, nonce string) error {
	query := `
    UPDATE passwords
    SET encrypted_password = ?
    WHERE name = ?`
	_, err := s.db.Exec(query, encryptedPassword, name)
	return err
}

func (s *Storage) Delete(name string) error {
	query := `DELETE FROM passwords WHERE name = ?`
	_, err := s.db.Exec(query, name)
	return err
}
