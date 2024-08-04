package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
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
	vaultsQuery := `CREATE TABLE IF NOT EXISTS vaults (
        uuid TEXT PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
        hashed_master_password TEXT NOT NULL,
		salt TEXT NOT NULL,
        date_created DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := s.db.Exec(vaultsQuery)
	if err != nil {
		return err
	}

	passwordsQuery := `CREATE TABLE IF NOT EXISTS passwords (
        uuid TEXT PRIMARY KEY,
        vault_uuid TEXT NOT NULL,
        name TEXT NOT NULL UNIQUE,
        username TEXT NOT NULL,
        encrypted_password TEXT NOT NULL,
        date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(vault_uuid) REFERENCES vaults(uuid)
    );`
	_, err = s.db.Exec(passwordsQuery)
	return err
}

func (s *Storage) NewVault(name, description, hashedMaster, salt string) (int64, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return 0, err
	}
	query := `
    INSERT INTO vaults (uuid, name, description, hashed_master_password, salt)
    VALUES (?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, uuid, name, description, hashedMaster, salt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Storage) GetVault(name string) (map[string]string, error) {
	query := `SELECT uuid, description, hashed_master_password, salt, date_created FROM vaults WHERE name = ?`
	row := s.db.QueryRow(query, name)
	var uuid, description, hashed_master_password, salt, date_created string
	err := row.Scan(&description, &hashed_master_password, &date_created)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"uuid":                   uuid,
		"description":            description,
		"hashed_master_password": hashed_master_password,
		"salt":                   salt,
		"date_created":           date_created,
	}, nil
}

func (s *Storage) DeleteVault(name string) error {

	return nil
}

func (s *Storage) SaveLogin(vault_uuid, name, username, encryptedPassword string) (int64, error) {
	query := `
    INSERT INTO passwords (vault_uuid, name, username, encrypted_password, salt)
    VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, vault_uuid, name, username, encryptedPassword)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (s *Storage) ReadLogin(vault_uuid, name string) (map[string]string, error) {
	query := `SELECT username, encrypted_password, salt FROM passwords WHERE name = ? AND vault_uuid = ?`
	row := s.db.QueryRow(query, name, vault_uuid)

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

func (s *Storage) ListLogin(vault_uuid, name, username string) ([]map[string]string, error) {
	query := `SELECT username, name FROM passwords`
	var conditions []string
	var args []interface{}

	if vault_uuid != "" {
		conditions = append(conditions, "vault_uuid = ?")
		args = append(args, vault_uuid)
	}
	if name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}
	if username != "" {
		conditions = append(conditions, "username LIKE ?")
		args = append(args, "%"+username+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (s *Storage) Update(vault_uuid, target, name, username, encryptedPassword, salt string) error {
	var args []interface{}
	query := `UPDATE passwords SET`
	whereClause := " WHERE name = ?, vault_uuid = ?"
	if name != "" {
		query += " name = ?"
		args = append(args, name)
	}
	if username != "" {
		query += ", username = ?"
		args = append(args, username)
	}
	if encryptedPassword != "" && salt != "" {
		query += ", encrypted_password = ?, salt = ?"
		args = append(args, encryptedPassword)
		args = append(args, salt)
	}

	query += whereClause
	args = append(args, target)
	args = append(args, vault_uuid)

	_, err := s.db.Exec(query, args...)
	return err
}

func (s *Storage) DeleteLogin(vault_uuid, name string) error {
	var args []interface{}
	query := `DELETE FROM passwords WHERE`

	if name != "" {
		query += " name = ? AND"
		args = append(args, name)
	}
	query += " vault_uuid = ?"
	args = append(args, vault_uuid)
	_, err := s.db.Exec(query, args...)
	return err
}
