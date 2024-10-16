package product

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/neoxelox/kit"
	kitUtil "github.com/neoxelox/kit/util"
	"github.com/rs/xid"
	"github.com/scylladb/go-set/strset"

	"backend/pkg/config"
	"backend/pkg/engine"
	"backend/pkg/organization"
	"backend/pkg/util"
)

type ProductEndpoints struct {
	config            config.Config
	observer          *kit.Observer
	productRepository *ProductRepository
}

func NewProductEndpoints(observer *kit.Observer, productRepository *ProductRepository,
	config config.Config) *ProductEndpoints {
	return &ProductEndpoints{
		config:            config,
		observer:          observer,
		productRepository: productRepository,
	}
}

type ProductEndpointsListProductsResponse struct {
	Products []ProductPayload `json:"products"`
}

func (self *ProductEndpoints) ListProducts(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)

	products, err := self.productRepository.ListByOrganizationID(requestCtx, requestOrganization.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := ProductEndpointsListProductsResponse{}
	response.Products = make([]ProductPayload, 0, len(products))
	for _, product := range products {
		if product.DeletedAt != nil {
			continue
		}

		response.Products = append(response.Products, *NewProductPayload(product))
	}

	return ctx.JSON(http.StatusOK, &response)
}

type ProductEndpointsPostProductRequest struct {
	Name       string    `json:"name"`
	Picture    *string   `json:"picture"`
	Language   string    `json:"language"`
	Context    *string   `json:"context"`
	Categories *[]string `json:"categories"`
	Release    *string   `json:"release"`
}

type ProductEndpointsPostProductResponse struct {
	ProductPayload
}

func (self *ProductEndpoints) PostProduct(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestOrganization := organization.RequestOrganization(requestCtx)
	request := ProductEndpointsPostProductRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	// TODO: Create an organization middleware that limits the usage of features by plan
	if requestOrganization.Plan == organization.OrganizationPlanTrial ||
		requestOrganization.Plan == organization.OrganizationPlanDemo {
		exists, err := self.productRepository.ExistsByOrganizationID(requestCtx, requestOrganization.ID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if exists {
			return kit.HTTPErrInvalidRequest
		}
	}

	// TODO: Create an organization middleware that limits the usage of features by plan
	if requestOrganization.Plan == organization.OrganizationPlanStarter {
		count, err := self.productRepository.CountNotDeletedByOrganizationID(requestCtx, requestOrganization.ID)
		if err != nil {
			return kit.HTTPErrServerGeneric.Cause(err)
		}

		if count >= organization.ORGANIZATION_PLAN_STARTER_MAX_PRODUCTS {
			return kit.HTTPErrInvalidRequest
		}
	}

	if len(request.Name) == 0 {
		return kit.HTTPErrInvalidRequest
	}

	if request.Picture != nil {
		if !util.SameOrigin(*request.Picture, self.config.CDN.BaseURL) {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Picture = kitUtil.Pointer(PRODUCT_DEFAULT_PICTURE)
	}

	if !IsLanguageSupported(request.Language) {
		return kit.HTTPErrInvalidRequest
	}

	if request.Context != nil {
		if len(*request.Context) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		if len(*request.Context) > PRODUCT_MAX_CONTEXT_LENGTH {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Context = kitUtil.Pointer(fmt.Sprintf(PRODUCT_DEFAULT_CONTEXT[request.Language], request.Name))
	}

	if request.Categories != nil {
		if len(*request.Categories) > PRODUCT_MAX_CATEGORIES {
			return kit.HTTPErrInvalidRequest
		}

		for i := 0; i < len(*request.Categories); i++ {
			if len((*request.Categories)[i]) == 0 {
				return kit.HTTPErrInvalidRequest
			}

			if (*request.Categories)[i] == engine.OPTION_UNKNOWN {
				return kit.HTTPErrInvalidRequest
			}

			(*request.Categories)[i] = strings.ToUpper(strings.ReplaceAll((*request.Categories)[i], " ", "_"))
		}
		*request.Categories = strset.New(*request.Categories...).List()
	} else {
		request.Categories = kitUtil.Pointer([]string{})
	}

	if request.Release != nil {
		if len(*request.Release) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		if *request.Release == engine.OPTION_UNKNOWN {
			return kit.HTTPErrInvalidRequest
		}
	} else {
		request.Release = kitUtil.Pointer(engine.OPTION_UNKNOWN)
	}

	product := NewProduct()
	product.ID = xid.New().String()
	product.OrganizationID = requestOrganization.ID
	product.Name = request.Name
	product.Picture = *request.Picture
	product.Language = request.Language
	product.Context = *request.Context
	product.Categories = *request.Categories
	product.Release = *request.Release
	product.Settings = ProductSettings{}
	product.Usage = 0
	product.CreatedAt = time.Now()
	product.DeletedAt = nil

	product, err = self.productRepository.Create(requestCtx, *product)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := ProductEndpointsPostProductResponse{}
	response.ProductPayload = *NewProductPayload(*product)

	return ctx.JSON(http.StatusOK, &response)
}

type ProductEndpointsGetProductResponse struct {
	ProductPayload
}

func (self *ProductEndpoints) GetProduct(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := RequestProduct(requestCtx)

	response := ProductEndpointsGetProductResponse{}
	response.ProductPayload = *NewProductPayload(*requestProduct)

	return ctx.JSON(http.StatusOK, &response)
}

type ProductEndpointsPutProductRequest struct {
	Name       *string   `json:"name"`
	Picture    *string   `json:"picture"`
	Context    *string   `json:"context"`
	Categories *[]string `json:"categories"`
	Release    *string   `json:"release"`
}

type ProductEndpointsPutProductResponse struct {
	ProductPayload
}

func (self *ProductEndpoints) PutProduct(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := RequestProduct(requestCtx)
	request := ProductEndpointsPutProductRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	if request.Name != nil {
		if len(*request.Name) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		requestProduct.Name = *request.Name
	}

	if request.Picture != nil {
		if !util.SameOrigin(*request.Picture, self.config.CDN.BaseURL) {
			return kit.HTTPErrInvalidRequest
		}

		requestProduct.Picture = *request.Picture
	}

	if request.Context != nil {
		if len(*request.Context) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		if len(*request.Context) > PRODUCT_MAX_CONTEXT_LENGTH {
			return kit.HTTPErrInvalidRequest
		}

		requestProduct.Context = *request.Context
	}

	if request.Categories != nil {
		if len(*request.Categories) > PRODUCT_MAX_CATEGORIES {
			return kit.HTTPErrInvalidRequest
		}

		for i := 0; i < len(*request.Categories); i++ {
			if len((*request.Categories)[i]) == 0 {
				return kit.HTTPErrInvalidRequest
			}

			if (*request.Categories)[i] == engine.OPTION_UNKNOWN {
				return kit.HTTPErrInvalidRequest
			}

			(*request.Categories)[i] = strings.ToUpper(strings.ReplaceAll((*request.Categories)[i], " ", "_"))
		}
		*request.Categories = strset.New(*request.Categories...).List()

		requestProduct.Categories = *request.Categories
	}

	if request.Release != nil {
		if len(*request.Release) == 0 {
			return kit.HTTPErrInvalidRequest
		}

		if *request.Release == engine.OPTION_UNKNOWN {
			return kit.HTTPErrInvalidRequest
		}

		requestProduct.Release = *request.Release
	}

	err = self.productRepository.UpdateInfo(requestCtx, *requestProduct)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := ProductEndpointsPutProductResponse{}
	response.ProductPayload = *NewProductPayload(*requestProduct)

	return ctx.JSON(http.StatusOK, &response)
}

func (self *ProductEndpoints) DeleteProduct(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := RequestProduct(requestCtx)

	err := self.productRepository.Delete(requestCtx, requestProduct.ID)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	return ctx.JSON(http.StatusOK, struct{}{})
}

type ProductEndpointsGetProductSettingsResponse struct {
	ProductPayloadSettings
}

func (self *ProductEndpoints) GetProductSettings(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := RequestProduct(requestCtx)

	response := ProductEndpointsGetProductSettingsResponse{}
	response.ProductPayloadSettings = NewProductPayload(*requestProduct).Settings

	return ctx.JSON(http.StatusOK, &response)
}

type ProductEndpointsPutProductSettingsRequest struct {
}

type ProductEndpointsPutProductSettingsResponse struct {
	ProductPayloadSettings
}

func (self *ProductEndpoints) PutProductSettings(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := RequestProduct(requestCtx)
	request := ProductEndpointsPutProductSettingsRequest{}

	err := ctx.Bind(&request)
	if err != nil {
		return kit.HTTPErrInvalidRequest.Cause(err)
	}

	err = self.productRepository.UpdateSettings(requestCtx, *requestProduct)
	if err != nil {
		return kit.HTTPErrServerGeneric.Cause(err)
	}

	response := ProductEndpointsPutProductSettingsResponse{}
	response.ProductPayloadSettings = NewProductPayload(*requestProduct).Settings

	return ctx.JSON(http.StatusOK, &response)
}

type ProductEndpointsGetProductUsageResponse struct {
	Usage int `json:"usage"`
}

func (self *ProductEndpoints) GetProductUsage(ctx echo.Context) error {
	requestCtx := ctx.Request().Context()
	requestProduct := RequestProduct(requestCtx)

	response := ProductEndpointsGetProductUsageResponse{}
	response.Usage = requestProduct.Usage

	return ctx.JSON(http.StatusOK, &response)
}
