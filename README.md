# pg_user

[![Release](https://img.shields.io/github/v/release/frederikhs/pg_user.svg)](https://github.com/frederikhs/pg_user/releases/latest)
[![GoDoc](https://godoc.org/github.com/frederikhs/pg_user?status.svg)](https://godoc.org/github.com/frederikhs/pg_user)
[![Quality](https://goreportcard.com/badge/github.com/frederikhs/pg_user)](https://goreportcard.com/report/github.com/frederikhs/pg_user)
[![Test](https://github.com/frederikhs/pg_user/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/frederikhs/pg_user/actions/workflows/test.yml)

Is a cli application for managing postgres users.

Features

- Creating new users with roles
- Deleting existing users
- Extending the validity of existing users
- Resetting passwords (random or specific) of existing users
- Listing existing users and their associated roles
- Listing all configured hosts
- Listing existing roles 

## Installation

### Linux amd64

```bash
# install
curl -L https://github.com/frederikhs/pg_user/releases/latest/download/pg_user_Linux_x86_64.tar.gz -o pg_user.tar.gz
tar -xvf pg_user.tar.gz
sudo chmod +x pg_user
sudo mv pg_user /usr/local/bin/pg_user

# clean up
rm pg_user.tar.gz
```

### Other

Other distributions or OS visit the [releases page](https://github.com/frederikhs/pg_user/releases/latest)

## Configuration

`pg_user` uses the format and location of the `.pgpass` file the resides in the `$HOME` directory of the user running the program.

---

Example of a host details in the `.pgpass` format

```
<HOST>:<PORT>:<DATABASE>:<USERNAME>:<PASSWORD>
```

## Examples of usage

`$ pg_user add test@example.org --host mydatabase.com`

This requires the hostname `mydatabase.com` must be configured in the `.pgpass` file

To also add roles to the user when creating use

`$ pg_user add test@example.org --host mydatabase.com --roles role1,role2`

## Development

This repository is very plain and simple and development of the application only required `go`.

`$ go run main.go hosts` will compile the application and list all configured hosts in your `.pgpass` file.

## Release

release a new version by creating a tag and pushing it, then goreleaser will do the rest

```bash
git tag -a v0.x.0 -m "v0.x.0"
git push origin v0.x.0
```
