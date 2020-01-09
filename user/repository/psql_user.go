package repository

import (
	"database/sql"
	"errors"

	"github.com/nebyubeyene/Intern-Seek-Version-1/entity"
)

//UserRepositoryImpl implements the user.UserRepository interface
type UserRepositoryImpl struct {
	conn *sql.DB
}

//NewUserRepositoryImpl will create a new UserRepositoryImpl
func NewUserRepositoryImpl(Conn *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{conn: Conn}
}

//StoreUser stores a user in database given a user
func (uri UserRepositoryImpl) StoreUser(user *entity.User) error {
	_, err := uri.conn.Exec("INSERT INTO users(username,full_name,email,phone,password) values($1,$2,$3,$4,$5)", user.Username, user.Name, user.Email, user.Phone, user.Password)
	if err != nil {
		return errors.New("Storing user has failed")
	}
	return nil
}

//UpdateUser updates a user in database given a user
func (uri UserRepositoryImpl) UpdateUser(user *entity.User) error {
	_, err := uri.conn.Exec("UPDATE users SET username=$1,full_name=$2,email=$3,phone=$4,password=$5 WHERE id=$6", user.Username, user.Name, user.Email, user.Phone, user.Password, user.ID)
	if err != nil {
		return errors.New("Updating user in the database has failed")
	}
	return nil
}

//DeleteUser deletes a user in database given an id
func (uri UserRepositoryImpl) DeleteUser(id int) error {
	_, err := uri.conn.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return errors.New("Deleting a user from database has failed")
	}
	return nil
}

//Users returns users from a database
func (uri UserRepositoryImpl) Users() ([]entity.User, error) {
	rows, err := uri.conn.Query("SELECT * FROM users")
	if err != nil {
		return nil, errors.New("Could not query the database")
	}
	listOfUsers := []entity.User{}
	for rows.Next() {
		u := entity.User{}
		err = rows.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Phone, &u.Password)
		if err != nil {
			return nil, err
		}
		listOfUsers = append(listOfUsers, u)
	}

	return listOfUsers, nil
}

//User returns a user from database
func (uri UserRepositoryImpl) User(id int) (*entity.User, error) {
	row := uri.conn.QueryRow("SELECT * FROM users WHERE id=$1", id)
	u := entity.User{}
	err := row.Scan(&u.ID, &u.Username, &u.Name, &u.Email, &u.Phone, &u.Password)
	if err != nil {
		return &u, err
	}
	return &u, nil
}
