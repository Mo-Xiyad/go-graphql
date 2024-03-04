package services

import (
	"server/pkg/model"

	"gorm.io/gorm"
)

type CompanyRepo struct {
	Db *gorm.DB
}

func NewCompanyRepo(db *gorm.DB) *CompanyRepo {
	return &CompanyRepo{
		Db: db,
	}
}

func (c *CompanyRepo) Create(company *model.Company) error {
	if err := c.Db.Create(&company).Error; err != nil {
		return err
	}

	return nil
}

func (c *CompanyRepo) GetAll() ([]*model.Company, error) {
	var companies []*model.Company
	if err := c.Db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}
