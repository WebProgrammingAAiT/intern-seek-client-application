package user

import "github.com/abdimussa87/Intern-Seek-Version-1/entity"

// UserRepository specifies user related database operations
type UserRepository interface {
	StoreUser(user entity.User) error
	UpdateUser(user entity.User) error
	DeleteUser(id int) error
	Users() ([]entity.User, error)
	User(id int) (entity.User, error)
}
type ApplicationRepository interface{
	StoreApplication(application *entity.Application)(*entity.Application,[]error)
	UpdateApplication(application *entity.Application)(*entity.Application,[]error)
	DeleteApplication(id uint)(*entity.Application,[]error)
	Applications()([]entity.Application,[]error)
	Application(id uint)(*entity.Application,[]error)
}
