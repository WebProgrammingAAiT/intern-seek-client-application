package user

import "github.com/nebyubeyene/Intern-Seek-Version-1/entity"

//UserService specifies user related services
type UserService interface {
	StoreUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(id int) error
	Users() ([]entity.User, error)
	User(id int) (*entity.User, error)
}
type CompanyService interface {
	StoreCompany(user_id int, company *entity.CompanyDetail) error
	UpdateCompany(company *entity.CompanyDetail) error
	DeleteCompany(id uint) error
	Companies() ([]entity.CompanyDetail, error)
	Company(id uint) (*entity.CompanyDetail, error)
}
