package user

import "github.com/abdimussa87/Intern-Seek-Version-1/entity"

//UserService specifies user related services
type UserService interface {
	StoreUser(user entity.User) error
	UpdateUser(user entity.User) error
	DeleteUser(id int) error
	Users() ([]entity.User, error)
	User(id int) (entity.User, error)
}
