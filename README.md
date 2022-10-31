# pg_user

Is a cli application for managing postgres users.

Features

- Creating new users with roles
- Deleting existing users
- Extending the validity of existing users
- Resetting passwords of existing users
- Listing of existing users
- Listing all configured hosts

## Installation

1. Download the latest release for your distribution or architecture at the [releases page](https://github.com/hiperdk/pg_user/releases/latest). This can only be done while authenticated with github, and thus a simple curl in the terminal is not sufficient.

2. Unzip the compressed binary and move to a location in path
   3. eg. `$ sudo mv pg_user /usr/local/bin/pg_user`

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

## Creating a new release

To create a new release of the application, a tag needs to be created first.

Look up existing tags

`$ git tag -l`

Choose the correct semantic versioning of the to be created release

Tag commit for release

`$ git tag vX.X.X`

Push tag

`$ git push origin vX.X.X`

Now that the tag has been pushed, a release needs to be create using GitHub. After creating a new release the [release](.github/workflows/release.yml) workflow will compile binaries for major platforms and architectures and attach them to the release.
