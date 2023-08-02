package database

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type User struct {
	Username   string   `db:"usename" json:"username"`
	ValidUntil *string  `db:"valuntil" json:"valid_until"`
	RolesJson  string   `db:"roles" json:"-"`
	Roles      []string `db:"-" json:"roles"`
}

func (u *User) ParseValidUntil() (*string, error) {
	if u.ValidUntil == nil || *u.ValidUntil == "infinity" {
		ts := "∞"

		return &ts, nil
	}

	t, err := time.Parse(time.RFC3339, *u.ValidUntil)
	if err != nil {
		return nil, err
	}

	ts := t.Format("2006-01-02")

	return &ts, nil
}

func (conn *DBConn) GetAllUsers() ([]User, error) {
	var users []User
	err := conn.db.Select(&users, `
	SELECT a.usename  AS usename,
		   a.valuntil AS valuntil,
		   json_agg(c.rolname) AS roles
	FROM pg_user a
			 LEFT JOIN pg_auth_members b ON a.usesysid = b.member
			 LEFT JOIN pg_roles c ON b.roleid = c.oid
	GROUP BY 1, 2
	ORDER BY 1;
    `)
	if err != nil {
		return nil, err
	}

	for i, u := range users {
		validUntil, err := u.ParseValidUntil()
		if err != nil {
			return nil, err
		}

		users[i].ValidUntil = validUntil

		// convert postgres json aggregation to string slice
		var roles []string
		err = json.Unmarshal([]byte(users[i].RolesJson), &roles)
		if err != nil {
			return nil, err
		}

		users[i].Roles = roles
	}

	return users, nil
}

func (conn *DBConn) GetAllRoles() ([]string, error) {
	var roles []string
	err := conn.db.Select(&roles, "SELECT rolname FROM pg_roles ORDER BY rolname")
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (conn *DBConn) CreateUser(username string, validDuration time.Duration, roles []string) (string, time.Time, error) {
	tx := conn.db.MustBegin()

	password := RandString(15)

	validUntil := time.Now().Add(validDuration)

	sql := fmt.Sprintf("CREATE USER \"%s\" WITH PASSWORD '%s' VALID UNTIL '%s'", username, password, validUntil.Format("2006-01-02"))
	tx.MustExec(sql)

	err := conn.AddRole(tx, username, roles)
	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return "", time.Now(), err2
		}

		return "", time.Now(), err
	}

	return password, validUntil, tx.Commit()
}

func (conn *DBConn) ExtendUser(username string, validDuration time.Duration) (time.Time, error) {
	tx := conn.db.MustBegin()

	user, err := conn.GetUser(username)
	if err != nil {
		return time.Now(), err
	}

	current, err := time.Parse("2006-01-02", *user.ValidUntil)
	if err != nil {
		return time.Now(), err
	}

	// if users are expired, extend from now not their expiration
	if current.Before(time.Now()) {
		current = time.Now()
	}

	validUntil := current.Add(validDuration)

	// at maximum give extend the valid until to 2 times the asked extension from now
	if time.Until(validUntil) >= validDuration*2 {
		validUntil = time.Now().Add(validDuration * 2)
	}

	sql := fmt.Sprintf("ALTER USER \"%s\" WITH VALID UNTIL '%s'", username, validUntil.Format("2006-01-02"))
	tx.MustExec(sql)

	return validUntil, tx.Commit()
}

func (conn *DBConn) ResetPassword(username string) (string, error) {
	return conn.SetPassword(username, RandString(15))
}

func (conn *DBConn) SetPassword(username string, password string) (string, error) {
	tx := conn.db.MustBegin()

	sql := fmt.Sprintf("ALTER USER \"%s\" WITH PASSWORD '%s'", username, password)
	tx.MustExec(sql)

	return password, tx.Commit()
}

func (conn *DBConn) DeleteUser(username string) error {
	tx := conn.db.MustBegin()

	sql := fmt.Sprintf("DROP USER \"%s\"", username)

	tx.MustExec(sql)
	return tx.Commit()
}

func (conn *DBConn) UserExist(username string) (bool, error) {
	var exists bool

	err := conn.db.Get(&exists, "SELECT EXISTS(SELECT FROM pg_catalog.pg_user WHERE usesuper = false AND usename = $1)", username)

	return exists, err
}

func (conn *DBConn) GetUser(username string) (*User, error) {
	var user User

	err := conn.db.Get(&user, "SELECT usename, valuntil FROM pg_catalog.pg_user WHERE usesuper = false AND usename = $1", username)
	if err != nil {
		return nil, err
	}

	validUntil, err := user.ParseValidUntil()
	if err != nil {
		return nil, err
	}

	user.ValidUntil = validUntil

	return &user, err
}

func (conn *DBConn) AddRole(tx *sqlx.Tx, username string, roles []string) error {
	for _, role := range roles {

		roleExists, err := conn.RoleExist(role)
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				return err2
			}

			return err
		}

		if !roleExists {
			return fmt.Errorf("role does not exist: %s", role)
		}

		sql := fmt.Sprintf("GRANT \"%s\" to \"%s\"", role, username)
		tx.MustExec(sql)
	}

	return nil
}

func (conn *DBConn) RoleExist(role string) (bool, error) {
	var exists bool

	err := conn.db.Get(&exists, "SELECT EXISTS (SELECT rolname FROM pg_roles WHERE rolname = $1)", role)

	return exists, err
}
