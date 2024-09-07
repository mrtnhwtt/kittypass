package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(databasePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Printf("failed to open sqlite database at path %s. err: %s", databasePath, err)
		return nil, StorageAccessError{}
	}

	storage := &Storage{db: db}

	if err := storage.initDB(); err != nil {
		db.Close()
		log.Println("failed to initialize sqlite database with required tables")
		return nil, StorageInitError{}
	}

	return storage, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) initDB() error {
	vaultsQuery := `CREATE TABLE IF NOT EXISTS vaults (
        uuid TEXT PRIMARY KEY,
        name TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
        hex_hashed_master_password TEXT NOT NULL,
        date_created DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := s.db.Exec(vaultsQuery)
	if err != nil {
		log.Printf("error when running create query for vaults table: %s", err)
		return err
	}

	passwordsQuery := `CREATE TABLE IF NOT EXISTS passwords ( 
		identifier TEXT NOT NULL UNIQUE,
        vault_uuid TEXT NOT NULL,
        name TEXT NOT NULL,
        username TEXT NOT NULL,
        hex_encrypted_password TEXT NOT NULL,
		hex_salt TEXT NOT NULL,
        date_created DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY(vault_uuid) REFERENCES vaults(uuid)
    );`
	_, err = s.db.Exec(passwordsQuery)
	if err != nil {
		log.Printf("error when running create query for passwords table: %s", err)
		return err
	}
	return nil
}

func (s *Storage) SaveVault(name, description, hexHashedMaster string) (int64, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return 0, fmt.Errorf("error while generating an uuid for the vault: %s", err)
	}
	query := `INSERT INTO vaults (uuid, name, description, hex_hashed_master_password) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, uuid, name, description, hexHashedMaster)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, StorageConstraintError{Field: "name", Type: "Vault"}
			}
		}
		log.Printf("failed to save vault. err: %s", err)
		return 0, StorageUpdateError{}
	}
	return result.LastInsertId()
}

func (s *Storage) GetVault(name string) (map[string]string, error) {
	query := `SELECT uuid, description, hex_hashed_master_password FROM vaults WHERE name = ?`
	row := s.db.QueryRow(query, name)
	var uuid, description, hexHashedMasterPassword string
	err := row.Scan(&uuid, &description, &hexHashedMasterPassword)
	if err != nil {
		log.Printf("failed to read entry from database: %s", err)
		if errors.Is(sql.ErrNoRows, err) {
			log.Printf("vault %s does not exist", name)
		}
		return nil, VaultNotFound{}
	}
	return map[string]string{
		"uuid":                       uuid,
		"description":                description,
		"hex_hashed_master_password": hexHashedMasterPassword,
	}, nil
}

func (s *Storage) ListVault(name string) ([]map[string]string, error) {
	var rows *sql.Rows
	var err error
	query := `SELECT name, description, date_created FROM vaults`
	if name != "" {
		query = query + " WHERE name LIKE ?"
		name = "%" + name + "%"
		rows, err = s.db.Query(query, name)
	} else {
		rows, err = s.db.Query(query)
	}
	if err != nil {
		log.Printf("failed to get result from database: %s", err)
		return nil, VaultNotFound{}
	}
	defer rows.Close()

	var vaultList []map[string]string
	for rows.Next() {
		var vaultName, vaultDesc, vaultCreationDate string
		err := rows.Scan(&vaultName, &vaultDesc, &vaultCreationDate)
		if err != nil {
			log.Printf("error while scanning results of query. err: %s", err)
			return nil, StorageReadError{}
		}
		vaultList = append(vaultList, map[string]string{"name": vaultName, "description": vaultDesc, "date_created": vaultCreationDate})
	}
	return vaultList, nil
}

func (s *Storage) UpdateVault(vaultUuid, newName, newDescription, newMasterPassHashedHex, newSalt string, loginList []map[string]string) (map[string]int, error) { //TODO: update this for login unique salt
	affectedLogin := 0
	affectedVault := 0
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %s", err)
		return nil, StorageUpdateError{}
	}
	defer func() {
		if err != nil {
			log.Printf("rolling back update because an error happened. err: %s", err)
			tx.Rollback()
		}
	}()
	if len(loginList) > 0 {
		loginQuery := `UPDATE passwords SET hex_encrypted_password = ? WHERE name = ? AND vault_uuid = ?`
		for _, login := range loginList {
			res, err := tx.Exec(loginQuery, login["newHexEncrypted"], login["name"], vaultUuid)
			if err != nil {
				log.Printf("failed to update passwords associated with the vault: %s", err)
				return nil, StorageUpdateError{}
			}

			aff, err := res.RowsAffected()
			if err != nil {
				log.Printf("could not get the number of updated password entries: %s", err)
				return nil, StorageUpdateError{}
			}
			affectedLogin += int(aff)
		}

		if affectedLogin != len(loginList) {
			log.Printf("failed to update all logins password. %d updated logins for %d login. Aborting update", affectedLogin, len(loginList))
			return nil, StorageUpdateError{}
		}
	}

	var args []interface{}
	vaultQuery := `UPDATE vaults SET`
	whereClause := " WHERE uuid = ?"
	var setClause []string
	if newName != "" {
		setClause = append(setClause, " name = ?")
		args = append(args, newName)
	}
	if newDescription != "" {
		setClause = append(setClause, " description = ?")
		args = append(args, newDescription)
	}
	if len(loginList) > 0 {
		setClause = append(setClause, " hex_salt = ?")
		args = append(args, newSalt)
		setClause = append(setClause, " hex_hashed_master_password = ?")
		args = append(args, newMasterPassHashedHex)
	}
	vaultQuery += strings.Join(setClause, ",")
	vaultQuery += whereClause
	args = append(args, vaultUuid)
	res, err := tx.Exec(vaultQuery, args...)
	if err != nil {
		log.Printf("failed to delete the vault: %s", err)
		return nil, StorageUpdateError{}
	}

	aff, err := res.RowsAffected()
	if err != nil {
		log.Printf("could not get the number of deleted vault entries: %s", err)
		return nil, StorageUpdateError{}
	}
	affectedVault += int(aff)

	if affectedVault != 1 {
		log.Println("failed to update vault.")
		return nil, StorageUpdateError{}
	}
	if err = tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %s", err)
		return nil, StorageUpdateError{}
	}
	return map[string]int{"updated_login": affectedLogin, "updated_vault": affectedVault}, nil
}

func (s *Storage) ReadLogins(vault_uuid string) ([]map[string]string, error) {
	query := `SELECT identifier, name, username, hex_encrypted_password, hex_salt FROM passwords WHERE vault_uuid = ?`
	rows, err := s.db.Query(query, vault_uuid)
	if err != nil {
		log.Printf("failed to query database for logins associated with vauld uuid %s. err: %s", vault_uuid, err)
		return nil, StorageReadError{}
	}
	defer rows.Close()

	var loginList []map[string]string
	for rows.Next() {
		var identifier, name, username, hexEncryptedPassword, hexSalt string
		err := rows.Scan(&identifier, &name, &username, &hexEncryptedPassword, &hexSalt)
		if err != nil {
			log.Printf("error while scanning results of query. err: %s", err)
			return nil, StorageReadError{}
		}
		loginList = append(loginList, map[string]string{"identifier": identifier, "name": name, "username": username, "hex_enc_pass": hexEncryptedPassword, "hex_salt": hexSalt})
	}
	return loginList, nil
}

func (s *Storage) DeleteVault(name, vault_uuid string) (map[string]int64, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("failed to begin transaction: %s", err)
		return nil, StorageReadError{}
	}

	defer func() {
		if err != nil {
			log.Printf("rolling back update because an error happened. err: %s", err)
			tx.Rollback()
		}
	}()

	loginQuery := `
		DELETE FROM passwords 
		WHERE vault_uuid = (
			SELECT uuid FROM vaults 
			WHERE name = ? AND uuid = ?
		)`
	res, err := tx.Exec(loginQuery, name, vault_uuid)
	if err != nil {
		log.Printf("failed to delete passwords associated with the vault: %s", err)
		return nil, StorageUpdateError{}
	}

	affectedLogin, err := res.RowsAffected()
	if err != nil {
		log.Printf("could not get the number of deleted password entries: %s", err)
		return nil, StorageUpdateError{}
	}

	vaultQuery := `DELETE FROM vaults WHERE name = ? AND uuid = ?`
	res, err = tx.Exec(vaultQuery, name, vault_uuid)
	if err != nil {
		log.Printf("failed to delete the vault: %s", err)
		return nil, StorageUpdateError{}
	}

	affectedVault, err := res.RowsAffected()
	if err != nil {
		log.Printf("could not get the number of deleted vault entries: %s", err)
		return nil, StorageUpdateError{}

	}

	if err = tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %s", err)
		return nil, StorageUpdateError{}

	}

	return map[string]int64{"delete_login": affectedLogin, "delete_vault": affectedVault}, nil
}

func (s *Storage) SaveLogin(vaultUuid, name, username, hexEncryptedPassword, hex_salt string) (int64, error) {
	identifier := vaultUuid + "_" + name
	query := `INSERT INTO passwords (vault_uuid, identifier, name, username, hex_encrypted_password, hex_salt) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := s.db.Exec(query, vaultUuid, identifier, name, username, hexEncryptedPassword, hex_salt)
	if err != nil {
		log.Printf("failed to add new login %s to vault %s. err: %s", name, vaultUuid, err)
		return 0, StorageUpdateError{}
	}
	return result.LastInsertId()
}

func (s *Storage) ReadLogin(vault_uuid, name string) (map[string]string, error) {
	query := `SELECT username, hex_encrypted_password, hex_salt FROM passwords WHERE name = ? AND vault_uuid = ?`
	row := s.db.QueryRow(query, name, vault_uuid)

	var username, hexEncryptedPassword, hexSalt string
	err := row.Scan(&username, &hexEncryptedPassword, &hexSalt)
	if err != nil {
		log.Printf("error while scanning results of query. err: %s", err)
		return nil, StorageReadError{}
	}

	return map[string]string{
		"name":                   name,
		"username":               username,
		"hex_encrypted_password": hexEncryptedPassword,
		"hex_salt":               hexSalt,
	}, nil
}

func (s *Storage) ListLogin(vault_uuid, name, username string) ([]map[string]string, error) {
	query := `SELECT 
		p.username,
		p.name,
		p.date_created,
		v.name
	FROM 
		passwords p
	INNER JOIN 
		vaults v 
	ON 
		p.vault_uuid = v.uuid`
	var conditions []string
	var args []interface{}

	if vault_uuid != "" {
		conditions = append(conditions, "p.vault_uuid = ?")
		args = append(args, vault_uuid)
	}
	if name != "" {
		conditions = append(conditions, "p.name LIKE ?")
		args = append(args, "%"+name+"%")
	}
	if username != "" {
		conditions = append(conditions, "p.username LIKE ?")
		args = append(args, "%"+username+"%")
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Printf("failed to query database for logins with name %s and username %s associated with vauld uuid %s. err: %s", name, username, vault_uuid, err)
		return nil, StorageReadError{}
	}
	defer rows.Close()

	var loginList []map[string]string
	for rows.Next() {
		var loginName, loginUsername, dateCreated, vaultName string
		err := rows.Scan(&loginUsername, &loginName, &dateCreated, &vaultName)
		if err != nil {
			log.Printf("error while scanning results of query. err: %s", err)
			return nil, StorageReadError{}
		}
		loginList = append(loginList, map[string]string{"name": loginName, "username": loginUsername, "timestamp": dateCreated, "vault_name": vaultName})
	}
	return loginList, nil
}

func (s *Storage) UpdateLogin(vaultUuid, target, name, username, hexEncryptedPassword, hexSalt string) (int64, error) {
	var args []interface{}
	query := `UPDATE passwords SET`
	whereClause := " WHERE name = ? AND vault_uuid = ?"
	var setClause []string
	if name != "" {
		setClause = append(setClause, " name = ?")
		args = append(args, name)
		setClause = append(setClause, " identifier = ?")
		args = append(args, vaultUuid+"_"+name)
	}
	if username != "" {
		setClause = append(setClause, " username = ?")
		args = append(args, username)
	}
	if hexEncryptedPassword != "" {
		setClause = append(setClause, " hex_encrypted_password = ?")
		args = append(args, hexEncryptedPassword)
		setClause = append(setClause, " hex_salt = ?")
		args = append(args, hexSalt)
	}
	query += strings.Join(setClause, ",")
	query += whereClause
	args = append(args, target)
	args = append(args, vaultUuid)
	res, err := s.db.Exec(query, args...)
	if err != nil {
		log.Printf("failed to update login name %s and username %s associated with vault uuid %s. err: %s", name, username, vaultUuid, err)
		return 0, StorageUpdateError{}
	}
	aff, err := res.RowsAffected()
	if err != nil {
		log.Printf("could not get the number of updated login entries: %s", err)
		return 0, StorageUpdateError{}
	}

	return aff, nil
}

func (s *Storage) DeleteLogin(vault_uuid, name string) error {
	query := `DELETE FROM passwords WHERE name = ? AND vault_uuid = ?`

	res, err := s.db.Exec(query, name, vault_uuid)
	if err != nil {
		log.Printf("failed to delete login name %s associated with vault uuid %s. err: %s", name, vault_uuid, err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("could not get the number of deleted login entries: %s", err)
		return StorageUpdateError{}
	}
	if affected < 1 {
		log.Printf("no log were deleted when attempting to delete login name %s associated with vault uuid %s. err: %s", name, vault_uuid, err)
		return LoginNotFound{}
	}
	return nil
}
