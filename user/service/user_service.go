package service

import (
	"github.com/abdimussa87/Intern-Seek-Version-1/entity"
	"github.com/abdimussa87/Intern-Seek-Version-1/user"
)

//UserServiceImpl implements user.UserService interface
type UserServiceImpl struct {
	userRepo user.UserRepository
}

//NewUserServiceImpl returns new UserServiceImpl
func NewUserServiceImpl(UserRepo user.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: UserRepo}
}

//StoreUser stores a user given a user
func (usi UserServiceImpl) StoreUser(user entity.User) error {
	err := usi.userRepo.StoreUser(user)
	if err != nil {
		return err
	}
	return nil
}

//UpdateUser updates a user given a user
func (usi UserServiceImpl) UpdateUser(user entity.User) error {
	err := usi.userRepo.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

//DeleteUser deletes a user given an id
func (usi UserServiceImpl) DeleteUser(id int) error {
	err := usi.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

//Users returns  list of users
func (usi UserServiceImpl) Users() ([]entity.User, error) {
	users, err := usi.userRepo.Users()
	if err != nil {
		return nil, err
	}
	return users, nil
}

//User returns a user given an id
func (usi UserServiceImpl) User(id int) (entity.User, error) {
	user, err := usi.userRepo.User(id)
	if err != nil {
		return user, err
	}
	return user, nil
}
