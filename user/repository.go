package user

import "github.com/nebyubeyene/Intern-Seek-Version-1/entity"

// UserRepository specifies user related database operations
type UserRepository interface {
	StoreUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(id int) error
	Users() ([]entity.User, error)
	User(id int) (*entity.User, error)
}

type CompanyRepository interface {
	StoreCompany(userid int, company *entity.CompanyDetail) error
	UpdateCompany(company *entity.CompanyDetail) error
	DeleteCompany(id uint) error
	Companies() ([]entity.CompanyDetail, error)
	Company(id uint) (*entity.CompanyDetail, error)
}
