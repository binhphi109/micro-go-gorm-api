package store

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/sample/sample-server/model"
)

type CompanyStore struct {
	Db     *gorm.DB
	Logger *logrus.Logger
}

func (s *Store) Company() *CompanyStore {
	if s.company == nil {
		s.company = &CompanyStore{
			Db:     s.Db,
			Logger: s.Logger,
		}
	}
	return s.company
}

func (us *CompanyStore) GetAll() ([]model.Company, error) {
	var companies []model.Company
	result := us.Db.Table("Companies").Find(&companies)

	if result.Error != nil {
		return companies, result.Error
	}

	return companies, nil
}

func (us *CompanyStore) Get(id string) (model.Company, error) {
	var company model.Company
	us.Db.Table("Companies").Where("id = ?", id).First(&company)

	if company.ID == "" {
		return company, fmt.Errorf("no Company found")
	}

	return company, nil
}

func (us *CompanyStore) GetByName(name string) (model.Company, error) {
	var company model.Company
	us.Db.Table("Companies").Where("name = ?", name).First(&company)

	if company.ID == "" {
		return company, fmt.Errorf("no Company found")
	}

	return company, nil
}

func (us *CompanyStore) Create(company model.Company) (model.Company, error) {
	var db = us.Db.Table("Companies").Create(&company)

	if db.Error != nil {
		us.Logger.WithError(db.Error).Fatal("Error creating company", company.Name)
		return model.Company{}, db.Error
	}

	return company, nil
}

func (us *CompanyStore) Update(company model.Company) (model.Company, error) {
	var db = us.Db.Table("Companies").Where("id = ?", company.ID).Updates(model.Company{
		Name:      company.Name,
		Deleted:   company.Deleted,
		UpdatedAt: company.UpdatedAt,
		DeletedAt: company.DeletedAt,
	})

	if db.Error != nil {
		us.Logger.WithError(db.Error).Fatal("Error updating company", company.Name)
		return model.Company{}, db.Error
	}

	return company, nil
}

func (us *CompanyStore) Delete(companyId string) error {
	var db = us.Db.Table("Companies").Where("id = ?", companyId).Delete(&model.Company{})

	if db.Error != nil {
		us.Logger.WithError(db.Error).Fatal("Error deleting company", companyId)
		return db.Error
	}

	return nil
}
