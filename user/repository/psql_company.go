package repository

import (
	"database/sql"
	"errors"

	"github.com/nebyubeyene/Intern-Seek-Version-1/entity"
)

type CompanyRepositoryImpl struct {
	conn *sql.DB
}

func NewCompanyRepositoryImpl(Conn *sql.DB) *CompanyRepositoryImpl {
	return &CompanyRepositoryImpl{conn: Conn}
}

func (cri CompanyRepositoryImpl) StoreCompany(user_id int, company *entity.CompanyDetail) error {
	_, err := cri.conn.Exec("INSERT INTO company_detail(user_id,description,focus_area,country,city) values($1,$2,$3,$4,$5)", user_id, company.Description, company.FocusArea, company.Country, company.City)
	if err != nil {
		return errors.New("Storing user has failed")
	}
	return nil
}

func (cri CompanyRepositoryImpl) UpdateCompany(company *entity.CompanyDetail) error {
	_, err := cri.conn.Exec("UPDATE company_detail SET  focus_area=$1,description=$2,country=$3,city=$4 WHERE id=$5", company.FocusArea, company.Description, company.Country, company.City, company.ID)
	if err != nil {
		return errors.New("Updating user in the database has failed")
	}
	return nil
}
func (cri CompanyRepositoryImpl) Companies() ([]entity.CompanyDetail, error) {
	rows, err := cri.conn.Query("SELECT * FROM company_detail")
	if err != nil {
		return nil, errors.New("Could not query the database")
	}
	listOfCompanyD := []entity.CompanyDetail{}
	for rows.Next() {
		c := entity.CompanyDetail{}
		err = rows.Scan(&c.ID, &c.UserId, &c.Description, &c.FocusArea, &c.Country, &c.City)
		if err != nil {
			return nil, err
		}
		listOfCompanyD = append(listOfCompanyD, c)
	}

	return listOfCompanyD, nil
}

func (cri CompanyRepositoryImpl) Company(id uint) (*entity.CompanyDetail, error) {
	row := cri.conn.QueryRow("SELECT * FROM company_detail WHERE id=$1", id)
	c := entity.CompanyDetail{}
	err := row.Scan(&c.ID, &c.UserId, &c.Description, &c.FocusArea, &c.Country, &c.City)
	if err != nil {
		return &c, err
	}
	return &c, nil
}

func (cri CompanyRepositoryImpl) DeleteCompany(id uint) error {
	_, err := cri.conn.Exec("DELETE FROM company_detail WHERE id=$1", id)
	if err != nil {
		return errors.New("Deleting a company from database has failed")
	}
	return nil
}
