package product

type ProductPayloadSettings struct {
}

type ProductPayload struct {
	ID             string                 `json:"id"`
	OrganizationID string                 `json:"organization_id"`
	Name           string                 `json:"name"`
	Picture        string                 `json:"picture"`
	Language       string                 `json:"language"`
	Context        string                 `json:"context"`
	Categories     []string               `json:"categories"`
	Release        string                 `json:"release"`
	Settings       ProductPayloadSettings `json:"settings"`
	Usage          int                    `json:"usage"`
}

func NewProductPayload(product Product) *ProductPayload {
	return &ProductPayload{
		ID:             product.ID,
		OrganizationID: product.OrganizationID,
		Name:           product.Name,
		Picture:        product.Picture,
		Language:       product.Language,
		Context:        product.Context,
		Categories:     product.Categories,
		Release:        product.Release,
		Settings:       ProductPayloadSettings{},
		Usage:          product.Usage,
	}
}
