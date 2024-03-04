package services

import (
	"context"
	gql_model "server/graph/model"
	"server/pkg/model"
)

type ICompanyService interface {
	GetAllCompanies(ctx context.Context) ([]*model.Company, error)
	// GetCompanyByID(ctx context.Context, id string) (*model.Company, error)
	CreateCompany(ctx context.Context, company gql_model.CreateCompanyInput) (*model.Company, error)
	// UpdateCompany(ctx context.Context, company *model.Company) error
	// DeleteCompany(ctx context.Context, id string) error
}

type ICompanyRepo interface {
	Create(company *model.Company) error
	// Update(company *model.Company) error
	// Delete(id string) error
	GetAll() ([]*model.Company, error)
	// GetByID(id string) (*model.Company, error)
}
