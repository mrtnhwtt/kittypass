# Kittypass

<p align="center"> <img src="resource/image/kittypass-logo-xs.png" alt="Kittypass Logo" /> </p>

Kittypass is a CLI-based password manager written in Go that uses SQLite for storage.

Kittypass uses Vaults to organise and to securely store encrypted logins using a  masterpassword unique for each Vault. The passwords are encrypted using AES, and the Vault master password is securely hashed before storage with Bcrypt.

## Features

- **Secure Storage:** Passwords are encrypted using AES and stored in SQLite.
- **Organize Logins with folder like Vaults:** Easily add, retrieve, list, and delete logins from the command line.
- **Cross-Platform:** Supports Linux, WSL and Macos

## Requirements

On linux install one of the following package required to access the clipboard:

 `libx11-dev` or` xorg-dev` or `libX11-devel`

## Installation

### Build from source

If you have go 1.22 installed, you can directly build from source with the following command

```bash 
go build -o kittypass
```

WIP

## Usage

### Commands

WIP

### Example Usage

```sh
# Add a Vault with a master password
kittypass add vault -n myVault

# Add a new login and provide the password
kittypass add login --vault myVault --name github --username martin --password

# Add a new login and generate a new password
kittypass add login --vault myVault --name stackoverflow --username martin -sNU

# Retrieve a login
kittypass get --vault myVault --name github

# List all logins
kittypass list logins

# Update a login
kittypass update login --vault myVault --target stackoverflow --username martin@myemail.com

# Delete a login
kittypass delete login --name github
```

## Security

AES is used to encrypt your passwords before being stored in SQLite in hexadecimal format. A unique salt is used for each vaults.

## Future Updates

Planned Features:

- [ ] Instead of a unique salt per Vault, implement unique salt per Login.
- [ ] Develop a TUI
- [x] Integrate Viper for configuration management, allowing users to customize storage and logs location.
- [ ] Add support for other storage methods beyond SQLite
- [ ] Add support for other encryption algorithms

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.