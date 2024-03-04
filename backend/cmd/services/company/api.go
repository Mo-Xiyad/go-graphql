package services

import (
	"context"
	gql_model "server/graph/model"
	"server/pkg/model"
)

type CompanyService struct {
	CompanyRepo ICompanyRepo
}

func NewCompanyService(companyRepo ICompanyRepo) *CompanyService {
	return &CompanyService{
		CompanyRepo: companyRepo,
	}
}

func (c *CompanyService) GetAllCompanies(ctx context.Context) ([]*model.Company, error) {
	companies, err := c.CompanyRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return companies, nil
}

// func (c *CompanyService) GetCompanyByID(ctx context.Context, id string) (*model.Company, error) {
// 	company, err := c.CompanyRepo.GetByID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return company, nil
// }

func (c *CompanyService) CreateCompany(ctx context.Context, company gql_model.CreateCompanyInput) (*model.Company, error) {
	companyModel := model.Company{
		Name:           company.Name,
		Email:          company.Email,
		PhoneNumber:    company.PhoneNumber,
		OrganizationID: company.OrganizationID,
	}
	err := c.CompanyRepo.Create(&companyModel)
	if err != nil {
		return nil, err
	}
	return &companyModel, nil
}
