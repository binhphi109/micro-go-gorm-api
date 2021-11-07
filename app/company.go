package app

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sample/sample-server/core"
	"github.com/sample/sample-server/model"
	"github.com/sample/sample-server/store"
)

type CompanyApp struct {
	Store  *store.Store
	Logger *logrus.Logger
	Config *core.Config
}

func (a *App) Company() *CompanyApp {
	if a.company == nil {
		a.company = &CompanyApp{
			Store:  a.Store,
			Logger: a.Logger,
			Config: a.Config,
		}
	}
	return a.company
}

func (ua *CompanyApp) GetAllCompanies() ([]model.Company, error) {
	companies, err := ua.Store.Company().GetAll()
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (ua *CompanyApp) GetCompany(companyId string) (*model.Company, error) {
	company, err := ua.Store.Company().Get(companyId)
	if err != nil {
		return nil, err
	}

	return &company, nil
}

func (ua *CompanyApp) CreateCompany(company model.Company) (*model.Company, error) {
	oldCompany, err := ua.Store.Company().GetByName(company.Name)
	if err != nil {
		return nil, err
	}

	//checks if name is already register or not
	if oldCompany.Name != "" {
		return nil, fmt.Errorf("name already in use")
	}

	//insert company details in database
	newCompany, err := ua.Store.Company().Create(company)
	if err != nil {
		return nil, err
	}

	return &newCompany, nil
}

func (ua *CompanyApp) UpdateCompany(company model.Company) (model.Company, error) {
	oldCompany, err := ua.Store.Company().Get(company.ID)
	if err != nil {
		return company, err
	}

	oldCompany.Name = company.Name
	oldCompany.UpdatedAt = time.Now()

	//update company details in database
	oldCompany, err = ua.Store.Company().Update(oldCompany)
	if err != nil {
		return company, err
	}

	return oldCompany, nil
}

func (ua *CompanyApp) DeleteCompany(companyId string) error {
	company, err := ua.Store.Company().Get(companyId)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	company.Deleted = true
	company.DeletedAt = &timeNow

	//soft delete company in database
	_, err = ua.Store.Company().Update(company)
	if err != nil {
		return err
	}

	return nil
}
