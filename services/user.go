package services

import "database/sql"

type User struct {
	Db *sql.DB
}

func (u *User) CreateUser(name string) {}
func (u *User) GetUserID(name string) {

	// return userId maybe? because this could then be used to query attendance?
}
