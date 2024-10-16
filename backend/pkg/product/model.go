package product

import (
	"encoding/json"
	"time"
)

const (
	PRODUCT_MODEL_TABLE = "\"product\""
)

type ProductModel struct {
	ID             string     `db:"id"`
	OrganizationID string     `db:"organization_id"`
	Name           string     `db:"name"`
	Picture        string     `db:"picture"`
	Language       string     `db:"language"`
	Context        string     `db:"context"`
	Categories     []string   `db:"categories"`
	Release        string     `db:"release"`
	Settings       []byte     `db:"settings"`
	Usage          int        `db:"usage"`
	CreatedAt      time.Time  `db:"created_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

func NewProductModel(product Product) *ProductModel {
	settings, err := json.Marshal(product.Settings)
	if err != nil {
		panic(err)
	}

	return &ProductModel{
		ID:             product.ID,
		OrganizationID: product.OrganizationID,
		Name:           product.Name,
		Picture:        product.Picture,
		Language:       product.Language,
		Context:        product.Context,
		Categories:     product.Categories,
		Release:        product.Release,
		Settings:       settings,
		Usage:          product.Usage,
		CreatedAt:      product.CreatedAt,
		DeletedAt:      product.DeletedAt,
	}
}

func (self *ProductModel) ToEntity() *Product {
	var settings ProductSettings
	err := json.Unmarshal(self.Settings, &settings)
	if err != nil {
		panic(err)
	}

	return &Product{
		ID:             self.ID,
		OrganizationID: self.OrganizationID,
		Name:           self.Name,
		Picture:        self.Picture,
		Language:       self.Language,
		Context:        self.Context,
		Categories:     self.Categories,
		Release:        self.Release,
		Settings:       settings,
		Usage:          self.Usage,
		CreatedAt:      self.CreatedAt,
		DeletedAt:      self.DeletedAt,
	}
}
