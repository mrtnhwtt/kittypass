package storage

type StorageAccessError struct {}

func (e StorageAccessError) Error() string {
	return "failed to access storage"
}

type StorageInitError struct {}

func (e StorageInitError) Error() string {
	return "failed to initialize storage"
}

type StorageUpdateError struct {}

func (e StorageUpdateError) Error() string {
	return "failed to update storage. Changes were not saved"
}

type StorageReadError struct {}

func (e StorageReadError) Error() string {
	return "failed to read a result from the storage"
}

type VaultNotFound struct {}

func (e VaultNotFound) Error() string {
	return "vault not found"
}

type LoginNotFound struct {}

func (e LoginNotFound) Error() string {
	return "login not found"
}