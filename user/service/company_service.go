package service

import (
	"github.com/nebyubeyene/Intern-Seek-Version-1/entity"
	"github.com/nebyubeyene/Intern-Seek-Version-1/user"
)

// UserService implements menu.UserService interface
type CompanyService struct {
	companyRepo user.CompanyRepository
}

// NewUserService  returns a new UserService object
func NewCompanyService(companyRepository user.CompanyRepository) *CompanyService {
	return &CompanyService{companyRepo: companyRepository}
}

// Users returns all stored application users
func (cs *CompanyService) Companies() ([]entity.CompanyDetail, error) {
	comps, err := cs.companyRepo.Companies()
	if err != nil {
		return nil, err
	}
	return comps, err
}

// User retrieves an application user by its id
func (cs *CompanyService) Company(id uint) (*entity.CompanyDetail, error) {
	comp, err := cs.companyRepo.Company(id)
	if err != nil {
		return nil, err
	}
	return comp, err
}

// UpdateUser updates  a given application user
func (cs *CompanyService) UpdateCompany(company *entity.CompanyDetail) error {
	err := cs.companyRepo.UpdateCompany(company)
	if err != nil {
		return err
	}
	return err
}

// DeleteUser deletes a given application user
func (cs *CompanyService) DeleteCompany(id uint) error {
	err := cs.companyRepo.DeleteCompany(id)
	if err != nil {
		return err
	}
	return err
}

// StoreUser stores a given application user
func (cs *CompanyService) StoreCompany(user_id int, company *entity.CompanyDetail) error {
	err := cs.companyRepo.StoreCompany(user_id, company)
	if err != nil {
		return err
	}
	return err
}
