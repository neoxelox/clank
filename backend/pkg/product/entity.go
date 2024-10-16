package product

import (
	"fmt"
	"time"

	"github.com/neoxelox/kit/util"
)

const (
	PRODUCT_DEFAULT_PICTURE    = "https://clank.so/images/pictures/product.png"
	PRODUCT_MAX_CONTEXT_LENGTH = 2500
	PRODUCT_MAX_CATEGORIES     = 25
)

var (
	PRODUCT_LANGUAGES_SUPPORTED = []string{"ENGLISH"}
	PRODUCT_DEFAULT_CONTEXT     = map[string]string{
		"ENGLISH": "The product is called %s.",
	}
)

func IsLanguageSupported(language string) bool {
	for i := 0; i < len(PRODUCT_LANGUAGES_SUPPORTED); i++ {
		if language == PRODUCT_LANGUAGES_SUPPORTED[i] {
			return true
		}
	}

	return false
}

type ProductSettings struct {
}

type Product struct {
	ID             string
	OrganizationID string
	Name           string
	Picture        string
	Language       string
	Context        string
	Categories     []string
	Release        string
	Settings       ProductSettings
	Usage          int
	CreatedAt      time.Time
	DeletedAt      *time.Time
}

func NewProduct() *Product {
	return &Product{}
}

func (self Product) String() string {
	return fmt.Sprintf("<Product: %s (%s)>", self.Name, self.ID)
}

func (self Product) Equals(other Product) bool {
	return util.Equals(self, other)
}

func (self Product) Copy() *Product {
	return util.Copy(self)
}
