package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	core "fermion/backend_core/controllers/cores"
	eda "fermion/backend_core/controllers/eda"
	location_service "fermion/backend_core/controllers/mdm/locations"
	model_core "fermion/backend_core/internal/model/core"
	mdm "fermion/backend_core/internal/model/mdm"
	pagination "fermion/backend_core/internal/model/pagination"
	core_repo "fermion/backend_core/internal/repository"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	ipaas_utils "fermion/backend_core/ipaas_core/utils"
	access_checker "fermion/backend_core/pkg/util/access"
	helpers "fermion/backend_core/pkg/util/helpers"

	qrcode "github.com/skip2/go-qrcode"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License v3.0 as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License v3.0 for more details.
You should have received a copy of the GNU Lesser General Public License v3.0
along with this program.  If not, see <https://www.gnu.org/licenses/lgpl-3.0.html/>.
*/
type Service interface {

	//------------------Product Brand-----------------------------------------------------------------------------------
	CreateBrand(data *mdm.ProductBrand, token_id string, access_template_id string) error
	UpdateBrand(data *mdm.ProductBrand, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteBrand(query map[string]interface{}, token_id string, access_template_id string) error
	GetAllBrand(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]BrandRequestAndResponseDTO, error)
	SearchBrand(query string, token_id string, access_template_id string) (interface{}, error)

	//------------------Product Category--------------------------------------------------------------------------------------
	CreateCategory(data *mdm.ProductCategory, token_id string, access_template_id string) error
	UpdateCategory(data *mdm.ProductCategory, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteCategory(query map[string]interface{}, token_id string, access_template_id string) error
	GetAllCategory(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]CategoryAndSubcategoryRequestAndResponseDTO, error)
	GetAllSubCategory(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]CategoryAndSubcategoryRequestAndResponseDTO, error)
	SearchCategory(query string, token_id string, access_template_id string, access_action string) (interface{}, error)

	//------------------Product Base Attributes--------------------------------------------------------------------------------
	CreateBaseAttribute(data *mdm.ProductBaseAttributes, token_id string, access_template_id string) error
	UpdateBaseAttribute(data *mdm.ProductBaseAttributes, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteBaseAttribute(query map[string]interface{}, token_id string, access_template_id string) error
	GetAllBaseAttribute(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductBaseAttributesRequestAndResponseDTO, error)

	//------------------Product Base Attributes Values----------------------------------------------------------------------------
	CreateBaseAttributeValue(data *mdm.ProductBaseAttributesValues, token_id string, access_template_id string) error
	UpdateBaseAttributeValue(data *mdm.ProductBaseAttributesValues, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteBaseAttributeValue(query map[string]interface{}, token_id string, access_template_id string) error
	GetAllBaseAttributeValue(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductBaseAttributesValuesRequestAndResponseDTO, error)

	//------------------Product Selected Attributes----------------------------------------------------------------------------
	CreateSelectedAttribute(data *mdm.ProductSelectedAttributes, token_id string, access_template_id string) error
	UpdateSelectedAttribute(data *mdm.ProductSelectedAttributes, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteSelectedAttribute(query map[string]interface{}, token_id string, access_template_id string) error
	GetAllSelectedAttribute(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductSelectedAttributesRequestAndREsponseDTO, error)

	//------------------Product Selected Attributes Values----------------------------------------------------------------------
	CreateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues, token_id string, access_template_id string) error
	UpdateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteSelectedAttributeValue(query map[string]interface{}, token_id string, access_template_id string) error
	GetAllSelectedAttributeValue(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductSelectedAttributesValuesRequestAndResponseeDTO, error)

	//---------------------Product Template-------------------------------------------------------------------------------------
	CreateTemplate(data CreateProductTemplatePayload, query map[string]interface{}, host string, scheme string, user_token string) (interface{}, error)
	UpdateTemplate(data CreateProductTemplatePayload, query map[string]interface{}, token_id string, access_template_id string, query2 map[string]interface{}, host string, scheme string, user_token string) error
	DeleteTemplate(query map[string]interface{}, token_id string, access_template_id string) error
	GetTemplateView(templateQuery map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetAllTemplate(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)

	//------------------Product Variant---------------------------------------------------------------------------------
	CreateVariant(data CreateProductVariantDTO, token_id string, access_template_id string) error
	UpdateVariant(data CreateProductVariantDTO, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteVariant(query map[string]interface{}, token_id string, access_template_id string) error
	GetVariantView(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetAllVariant(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
	Qrcode(id uint, token_id string, access_template_id string) error
	Skucodecheck(skucode string, token_id string, access_template_id string) error

	//------------------Product Bundles------------------------------------------------------------------------------------
	CreateBundle(data CreateBundlePayload, token_id string, access_template_id string) (interface{}, error)
	UpdateBundle(data CreateBundlePayload, query map[string]interface{}, token_id string, access_template_id string) error
	DeleteBundle(query map[string]interface{}, token_id string, access_template_id string) error
	GetOneBundle(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetAllBundle(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]mdm.ProductBundles, error)

	//------------------Product BUlk-----------------------------------------------------------------------------------------
	BulkCreateProduct(data *[]mdm.ProductTemplate, token_id string, access_template_id string) error

	//---------------------Product UpsertAPi's--------------------------------------------------------------------------------------
	UpsertProductTemplate(data []interface{}, meta_data map[string]interface{}) (interface{}, error)
	UpsertProductVariant(data []interface{}, meta_data map[string]interface{}) (interface{}, error)
	GetAllChannelProducts(p *pagination.Paginatevalue, queryParams QueryParamsDTO, token_id string, access_template_id string, access_action string) ([]ChannelMapResponseDTO, error)

	//--------------------------Filter Tabs---------------------------------------------------------------------------------------------------
	GetProductVariantTab(product_variant_id string, tab_name string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)

	// -------------------------------hsn------------------------------
	CreateHsn(data mdm.HSNCodesData) (map[string]interface{}, error)
	GetAllHsn(page *pagination.Paginatevalue) ([]mdm.HSNCodesData, error)
	GetOneHsn(query map[string]interface{}) (mdm.HSNCodesData, error)
	UpdateHsn(query map[string]interface{}, data mdm.HSNCodesData) error
	DeleteHsn(query map[string]interface{}) error
}

type service struct {
	productRepository   mdm_repo.Products
	vendorRepository    mdm_repo.Vendors
	inventoryRepository mdm_repo.BasicInventory
	pricingRepository   mdm_repo.Pricing
	userRepository      core_repo.User
	coreService         core.Service
	location_service    location_service.Service
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	productRepository := mdm_repo.NewProducts()
	vendorRepository := mdm_repo.NewVendors()
	inventoryService := mdm_repo.NewBasicInventory()
	pricingRepository := mdm_repo.NewPricing()
	userRepository := core_repo.NewUser()
	coreService := core.NewService()
	location_service := location_service.NewService()

	newServiceObj = &service{productRepository, vendorRepository, inventoryService, pricingRepository, userRepository, coreService, location_service}
	return newServiceObj
}

// ------------------Product Brand--------------------------------------------------------------------------------------------------
func (s *service) CreateBrand(data *mdm.ProductBrand, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	err := s.productRepository.CreateBrand(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateBrand(data *mdm.ProductBrand, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	err := s.productRepository.UpdateBrand(data, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteBrand(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteBrand(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetAllBrand(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]BrandRequestAndResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	var brands []BrandRequestAndResponseDTO
	result, err := s.productRepository.GetAllBrand(p)
	if err != nil {
		return []BrandRequestAndResponseDTO{}, err
	}
	jsondata, _ := json.Marshal(result)
	err = json.Unmarshal(jsondata, &brands)
	if err != nil {
		return nil, err
	}
	return brands, nil
}
func (s *service) SearchBrand(query string, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	//--------Fetch the available search location details from the Location--------
	result, err := s.productRepository.SearchBrand(query)
	var brandData []BrandRequestAndResponseDTO
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []BrandRequestAndResponseDTO{}, nil
	}

	//------Convert the information result into a BrandRequestAndResponseDTO(id, name)---------
	for _, v := range result {
		var data BrandRequestAndResponseDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		brandData = append(brandData, data)
	}
	return brandData, nil
}

// ------------------Product Category-----------------------------------------------------------------------------------------------
func (s *service) CreateCategory(data *mdm.ProductCategory, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	err := s.productRepository.CreateCategory(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateCategory(data *mdm.ProductCategory, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	err := s.productRepository.UpdateCategory(data, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteCategory(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteCategory(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetAllCategory(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]CategoryAndSubcategoryRequestAndResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	categories := make([]CategoryAndSubcategoryRequestAndResponseDTO, 0)
	if p.Per_page == 0 {
		p.Per_page = 1000
	}

	result, err := s.productRepository.GetAllCategory(p)
	if err != nil {
		return nil, err
	}

	jsondata, _ := json.Marshal(result)
	err = json.Unmarshal(jsondata, &categories)
	if err != nil {
		return nil, err
	}

	response := make([]CategoryAndSubcategoryRequestAndResponseDTO, 0)
	for _, category := range categories {
		if category.ParentCategoryID == nil {
			response = append(response, category)
			continue
		}
		if len(category.ChildCategoryIds) != 0 {
			response = append(response, category)
		}
	}
	return response, nil
}
func (s *service) GetAllSubCategory(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]CategoryAndSubcategoryRequestAndResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	subCategories := make([]CategoryAndSubcategoryRequestAndResponseDTO, 0)
	response := make([]CategoryAndSubcategoryRequestAndResponseDTO, 0)
	if p.Per_page == 0 {
		p.Per_page = 1000
	}
	result, err := s.productRepository.GetAllSubCategory(p)
	if err != nil {
		return nil, err
	}

	jsondata, _ := json.Marshal(result)
	err = json.Unmarshal(jsondata, &subCategories)
	if err != nil {
		return nil, err
	}

	for _, subCategory := range subCategories {
		if subCategory.ParentCategoryID != nil {
			response = append(response, subCategory)
		}
	}
	return response, nil
}
func (s *service) SearchCategory(query string, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	//--------Fetch the available search location details from the Location--------
	result, err := s.productRepository.SearchCategory(query)
	var categoryData []CategoryAndSubcategoryRequestAndResponseDTO
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []CategoryAndSubcategoryRequestAndResponseDTO{}, nil
	}

	//------Convert the information result into a CategoryAndSubcategoryRequestAndResponseDTO(id, name)---------
	for _, v := range result {
		var data CategoryAndSubcategoryRequestAndResponseDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		categoryData = append(categoryData, data)
	}

	return categoryData, nil
}

// ------------------Product Base Attributes---------------------------------------------------------------------------------------------
func (s *service) CreateBaseAttribute(data *mdm.ProductBaseAttributes, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	err := s.productRepository.CreateBaseAttribute(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateBaseAttribute(data *mdm.ProductBaseAttributes, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	err := s.productRepository.UpdateBaseAttribute(data, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteBaseAttribute(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteBaseAttribute(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetAllBaseAttribute(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductBaseAttributesRequestAndResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	var query map[string]interface{}
	result, err := s.productRepository.GetAllBaseAttribute(query, p)
	if err != nil {
		return nil, err
	}
	var response []ProductBaseAttributesRequestAndResponseDTO
	jsondata, _ := json.Marshal(result)
	json.Unmarshal(jsondata, &response)
	return response, nil
}

// ------------------Product Base Attributes Values--------------------------------------------------------------------------------------------
func (s *service) CreateBaseAttributeValue(data *mdm.ProductBaseAttributesValues, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	err := s.productRepository.CreateBaseAttributeValue(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateBaseAttributeValue(data *mdm.ProductBaseAttributesValues, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	err := s.productRepository.UpdateBaseAttributeValue(data, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteBaseAttributeValue(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteBaseAttributeValue(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetAllBaseAttributeValue(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductBaseAttributesValuesRequestAndResponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	var query map[string]interface{}
	result, err := s.productRepository.GetAllBaseAttributeValue(query, p)
	if err != nil {
		return nil, err
	}
	var response []ProductBaseAttributesValuesRequestAndResponseDTO
	jsondata, _ := json.Marshal(result)
	json.Unmarshal(jsondata, &response)

	return response, nil
}

// ------------------Product Selected Attributes--------------------------------------------------------------------------------------------
func (s *service) CreateSelectedAttribute(data *mdm.ProductSelectedAttributes, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	err := s.productRepository.CreateSelectedAttribute(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateSelectedAttribute(data *mdm.ProductSelectedAttributes, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	err := s.productRepository.UpdateSelectedAttribute(data, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteSelectedAttribute(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteSelectedAttribute(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetAllSelectedAttribute(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductSelectedAttributesRequestAndREsponseDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	var query map[string]interface{}
	result, err := s.productRepository.GetAllSelectedAttribute(query, p)
	if err != nil {
		return nil, err
	}

	var response []ProductSelectedAttributesRequestAndREsponseDTO
	jsondata, _ := json.Marshal(result)
	json.Unmarshal(jsondata, &response)

	return response, nil
}

// ------------------Product Selected Attributes Values------------------------------------------------------------------------------------------
func (s *service) CreateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	err := s.productRepository.CreateSelectedAttributeValue(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	err := s.productRepository.UpdateSelectedAttributeValue(data, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteSelectedAttributeValue(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteSelectedAttributeValue(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetAllSelectedAttributeValue(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]ProductSelectedAttributesValuesRequestAndResponseeDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	var query map[string]interface{}
	result, err := s.productRepository.GetAllSelectedAttributeValue(query, p)
	if err != nil {
		return nil, err
	}

	var response []ProductSelectedAttributesValuesRequestAndResponseeDTO
	jsondata, _ := json.Marshal(result)
	json.Unmarshal(jsondata, &response)

	return response, nil
}

//====================================TEMPLATE & VARIANT=====================================================================================================

// ---------------------Product Template-------------------------------------------------------------------------------------------------------------
func (s *service) CreateTemplate(data CreateProductTemplatePayload, query map[string]interface{}, host string, scheme string, user_token string) (interface{}, error) {
	// ================ finding company preference for the user =================
	result, _ := s.userRepository.FindUser(query)
	var user_companyId = make(map[string]interface{})
	serviceProvider := "LOCAL"
	if result.CompanyId != nil {
		if *result.CompanyId > 0 {
			user_companyId["id"] = result.CompanyId
			company, _ := s.coreService.GetCompanyPreferences(user_companyId)
			if company.FilePreferenceID != nil && *company.FilePreferenceID > 0 {
				serviceProvider, _ = helpers.GetLookupcodeName(int(*company.FilePreferenceID))
			}
		}
	}

	var Template mdm.ProductTemplate
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Template)
	if err != nil {
		return nil, err
	}

	//===================default variant===================================
	if Template.ProductVariantIds == nil || len(Template.ProductVariantIds) == 0 {
		var defaultVariant mdm.ProductVariant
		defaultVariant.SkuId = Template.SkuCode
		Template.ProductVariantIds = append(Template.ProductVariantIds, defaultVariant)
	}

	//===================default status===================================
	defaultStatus, err := helpers.GetLookupcodeId("PRODUCT_STATUS", "DRAFT")
	if err != nil {
		return nil, err
	}
	Template.StatusId = &defaultStatus

	// ========== based on file prefernces updated product template image ===============
	if serviceProvider != "LOCAL" {
		response, _ := ipaas_utils.UploadFile(data.ImageOptions, data.Name, data.SKUCode, uint(*data.CreatedByID), scheme, host, ipaas_utils.AWS, "MDM", "Products", ".png", ipaas_utils.BASE64, user_token)
		mdata, _ := json.Marshal(&response)
		json.Unmarshal(mdata, &Template.ImageOptions)
	} else {
		response := helpers.BASE64optimization(data.ImageOptions)

		mdata, _ := json.Marshal(&response)
		json.Unmarshal(mdata, &Template.ImageOptions)
	}

	//=====================variant section ================================
	for index := range Template.ProductVariantIds {
		Template.ProductVariantIds[index].ParentSkuId = Template.SkuCode
		Template.ProductVariantIds[index].ProductName = Template.Name
		var images []interface{}
		mdata, _ := json.Marshal(Template.ProductVariantIds[index].ImageOptions)
		json.Unmarshal(mdata, &images)
		if Template.ProductVariantIds[index].ImageOptions == nil || len(images) == 0 {
			Template.ProductVariantIds[index].ImageOptions = Template.ImageOptions
		} else if serviceProvider != "LOCAL" {
			// ========== based on file prefernces updated product variant image ===============
			response, _ := ipaas_utils.UploadFile(Template.ProductVariantIds[index].ImageOptions, Template.ProductVariantIds[index].ProductName, Template.ProductVariantIds[index].SkuId, uint(*data.CreatedByID), scheme, host, ipaas_utils.AWS, "MDM", "Products", ".png", ipaas_utils.BASE64, user_token)
			mdata, _ := json.Marshal(&response)
			json.Unmarshal(mdata, &Template.ProductVariantIds[index].ImageOptions)
		} else {
			response := helpers.BASE64optimization(Template.ProductVariantIds[index].ImageOptions)
			mdata, _ := json.Marshal(&response)
			json.Unmarshal(mdata, &Template.ProductVariantIds[index].ImageOptions)
		}
		Template.ProductVariantIds[index].Description = Template.Description
		Template.ProductVariantIds[index].CategoryID = Template.PrimaryCategoryID
		Template.ProductVariantIds[index].LeafCategoryID = Template.SecondaryCategoryID
		Template.ProductVariantIds[index].ShippingOptions = Template.ShippingOptions
		Template.ProductVariantIds[index].PackageDimensions = Template.PackageDimensions
		Template.ProductVariantIds[index].PackageMaterialOptions = Template.PackageMaterialOptions
		Template.ProductVariantIds[index].StatusId = Template.StatusId

		Template.ProductVariantIds[index].ProductPricingDetails.CostPrice = Template.ProductPricingDetails.CostPrice
		Template.ProductVariantIds[index].ProductPricingDetails.CurrencyId = Template.ProductPricingDetails.CurrencyId
		Template.ProductVariantIds[index].ProductPricingDetails.MRP = Template.ProductPricingDetails.MRP
		if Template.ProductVariantIds[index].ProductPricingDetails.SalesPrice == 0 {
			Template.ProductVariantIds[index].ProductPricingDetails.SalesPrice = Template.ProductPricingDetails.SalesPrice
		}
		Template.ProductVariantIds[index].CreatedByID = Template.CreatedByID
		Template.ProductVariantIds[index].CompanyId = Template.CompanyId
		Template.ProductVariantIds[index].Location = Template.Location
		Template.ProductVariantIds[index].DeliverySlaDetails = Template.DeliverySlaDetails
		Template.ProductVariantIds[index].ProductCancellationTerms = Template.ProductCancellationTerms
		Template.ProductVariantIds[index].ProductReturnTerms = Template.ProductReturnTerms
		Template.ProductVariantIds[index].ProductReplacementTerms = Template.ProductReplacementTerms
		Template.ProductVariantIds[index].FoodItemDetails = Template.FoodItemDetails
		Template.ProductVariantIds[index].InventoryDetail = Template.InventoryDetail
		Template.ProductVariantIds[index].ManufacturerDetails = Template.ManufacturerDetails
		Template.ProductVariantIds[index].ProductCriticalDetails = Template.ProductCriticalDetails
		Template.ProductVariantIds[index].AdditivesInformation = Template.AdditivesInformation
		Template.ProductVariantIds[index].ShortDescription = Template.ShortDescription
		Template.ProductVariantIds[index].DomainId = Template.DomainId
		Template.ProductVariantIds[index].Channel = Template.Channel
	}
	err = s.productRepository.CreateTemplate(&Template)
	if err != nil {
		return nil, err
	}

	// if err != nil {
	// 	id := uuid.New().String()
	// 	defaultStatus, err := helpers.GetLookupcodeId("PRODUCT_STATUS", "REJECTED")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	Template.StatusId = &defaultStatus
	// 	SkuCode := strings.Split(Template.SkuCode, "_InValid")[0]
	// 	Template.SkuCode = SkuCode + "_InValid" + id
	// 	productVariantCount := len(Template.ProductVariantIds)
	// 	for index := 0; productVariantCount > index; index++ {
	// 		VariantSkuCode := strings.Split(Template.ProductVariantIds[index].SkuId, "_InValid")[0]
	// 		Template.ProductVariantIds[index].SkuId = VariantSkuCode + "_InValid" + id
	// 		Template.ProductVariantIds[index].StatusId = &defaultStatus
	// 	}
	// 	err = s.productRepository.CreateTemplate(&Template)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	//================attribute section ==================================
	for _, attribute := range data.AttributeKeyValues {
		//-----------create attribute---------------------
		var SelectedAttributes mdm.ProductSelectedAttributes
		SelectedAttributes.TemplateId = Template.ID
		SelectedAttributes.AttributeId = attribute.ID
		SelectedAttributes.Name = attribute.Name
		err = s.productRepository.CreateSelectedAttribute(&SelectedAttributes)
		if err != nil {
			return nil, err
		}
		//----------------create attribute value-------------------
		for _, attribute_value := range attribute.AttributeValues {
			var SelectedAttributeValues mdm.ProductSelectedAttributesValues
			SelectedAttributeValues.TemplateId = Template.ID
			SelectedAttributeValues.AttributeId = attribute.ID
			SelectedAttributeValues.AttributeValueId = attribute_value.ID
			SelectedAttributeValues.Name = attribute_value.Name
			err = s.productRepository.CreateSelectedAttributeValue(&SelectedAttributeValues)
			if err != nil {
				return nil, err
			}
		}
	}

	return Template, nil
}
func (s *service) UpdateTemplate(data CreateProductTemplatePayload, query map[string]interface{}, token_id string, access_template_id string, query2 map[string]interface{}, host string, scheme string, user_token string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	var Template mdm.ProductTemplate
	// ================ finding company preference for the user =================
	result, _ := s.userRepository.FindUser(query2)
	var user_companyId = make(map[string]interface{})
	serviceProvider := "LOCAL"

	if result.CompanyId != nil {
		if *result.CompanyId > 0 {
			user_companyId["id"] = result.CompanyId
			company, _ := s.coreService.GetCompanyPreferences(user_companyId)
			if company.FilePreferenceID != nil && *company.FilePreferenceID > 0 {
				serviceProvider, _ = helpers.GetLookupcodeName(int(*company.FilePreferenceID))
			}
		}
	}

	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Template)
	if err != nil {
		return err
	}

	if serviceProvider != "LOCAL" && Template.ImageOptions != nil {
		response, _ := ipaas_utils.UploadFile(data.ImageOptions, data.Name, data.SKUCode, uint(*data.CreatedByID), scheme, host, ipaas_utils.AWS, "MDM", "Products", ".png", ipaas_utils.BASE64, user_token)
		mdata, _ := json.Marshal(&response)
		json.Unmarshal(mdata, &Template.ImageOptions)
	} else {
		response := helpers.BASE64optimization(data.ImageOptions)
		mdata, _ := json.Marshal(&response)
		json.Unmarshal(mdata, &data.ImageOptions)
	}

	//=================template Update =================================
	err = s.productRepository.UpdateTemplate(&Template, query)
	if err != nil {
		return err
	}
	//==================varinat Upadte =================================
	for index, product_variant := range Template.ProductVariantIds {
		if product_variant.SkuId == "" || product_variant.ProductTemplateId == 0 {
			return errors.New("product_variant must contain sku_id and product_template_id")
		}
		updatequery := map[string]interface{}{
			"product_template_id": uint(query["id"].(int)),
			"sku_id":              product_variant.SkuId,
		}
		product_variant.SkuId = ""
		Template.ProductVariantIds[index].ProductPricingDetails.CostPrice = Template.ProductPricingDetails.CostPrice
		Template.ProductVariantIds[index].ProductPricingDetails.CurrencyId = Template.ProductPricingDetails.CurrencyId
		Template.ProductVariantIds[index].ProductPricingDetails.MRP = Template.ProductPricingDetails.MRP
		Template.ProductVariantIds[index].UpdatedByID = Template.UpdatedByID

		if serviceProvider != "LOCAL" && Template.ProductVariantIds[index].ImageOptions != nil {
			var imageLinkArray = make([]map[string]interface{}, 0)
			imageLinkArray = append(imageLinkArray, map[string]interface{}{"data": Template.ProductVariantIds[index].ImageOptions})
			response, _ := ipaas_utils.UploadFile(Template.ProductVariantIds[index].ImageOptions, "Name", product_variant.SkuId, uint(*data.CreatedByID), scheme, host, ipaas_utils.AWS, "MDM", "Products", ".png", ipaas_utils.BASE64, user_token)
			mdata, _ := json.Marshal(&response)
			json.Unmarshal(mdata, &Template.ProductVariantIds[index].ImageOptions)

		} else {
			response := helpers.BASE64optimization(Template.ProductVariantIds[index].ImageOptions)
			mdata, _ := json.Marshal(&response)
			json.Unmarshal(mdata, &Template.ProductVariantIds[index].ImageOptions)
		}

		if Template.ProductVariantIds[index].ProductPricingDetails.SalesPrice == 0 {
			Template.ProductVariantIds[index].ProductPricingDetails.SalesPrice = Template.ProductPricingDetails.SalesPrice
		}
		variantResponse, _ := s.productRepository.FindVariant(updatequery)
		if variantResponse.ID != 0 {

			err = s.productRepository.UpdateVariant(&product_variant, updatequery)
			if err != nil {
				return err
			}
		}
		if variantResponse.ID == 0 {
			product_variant.ProductTemplateId = uint(query["id"].(int))
			err := s.productRepository.CreateVariant(&product_variant)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *service) DeleteTemplate(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	findquery := map[string]interface{}{
		"id": query["id"].(int),
	}

	//===================Find Template=========================
	_, err := s.productRepository.FindTemplate(findquery)
	if err != nil {
		return errors.New("record not found")
	}

	//===================Delete Variants of the Template========================
	varaintQuery := map[string]interface{}{
		"product_template_id": query["id"].(int),
		"user_id":             query["user_id"].(int),
	}
	err = s.productRepository.DeleteVariant(varaintQuery)
	if err != nil {
		return err
	}

	//===================Delete Template========================
	err = s.productRepository.DeleteTemplate(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetTemplateView(templateQuery map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view products at data level")
	}
	p := new(pagination.Paginatevalue)
	res, err := s.productRepository.FindTemplate(templateQuery)
	if err != nil {
		return nil, err
	}

	var response TemplateReponseDTO
	dto, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &response)
	if err != nil {
		return nil, err
	}

	for _, variant := range response.ProductVariantIds {

		//-------------------Variant Attribute Details ---------------------------------
		attribute_details := []map[string]interface{}{}
		if variant.AttributeKeyValuesId != nil {
			for _, attribute_value_id := range variant.AttributeKeyValuesId {

				query := map[string]interface{}{"id": uint(attribute_value_id)}
				attribute_value, err := s.productRepository.GetAllBaseAttributeValue(query, p)
				if err != nil {
					return nil, err
				}

				if len(attribute_value) > 0 {
					query = map[string]interface{}{"id": attribute_value[0].BaseAttributeId}
					attribute, err := s.productRepository.GetAllBaseAttribute(query, p)
					if err != nil {
						return nil, err
					}
					attribute_details = append(attribute_details, map[string]interface{}{
						"attribute_name":  attribute[0].Name,
						"attribute_value": attribute_value[0].Name,
					})
				}
			}
		}
		variant.AttributeValues = attribute_details

		//-------------------Price List Details ---------------------------------
		price_list_details := []map[string]interface{}{}
		if variant.VendorPriceListIds != nil {
			for _, price_list_id := range variant.VendorPriceListIds {
				query := map[string]interface{}{"id": price_list_id}
				vendorPriceList, err := s.vendorRepository.FindOnePriceList(query)
				if err != nil {
					return nil, err
				}

				var priceListInterface map[string]interface{}
				mdata, err := json.Marshal(vendorPriceList)
				if err != nil {
					return nil, err
				}
				err = json.Unmarshal(mdata, &priceListInterface)
				if err != nil {
					return nil, err
				}
				price_list_details = append(price_list_details, priceListInterface)
			}
		}
		variant.PriceListDetails = price_list_details

		if variant.StandardProductTypes != nil {
			variant.StandardProductTypes = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(variant.StandardProductTypes)
		}
	}

	//-------------------Product Template Attributes--------------------------------------------------
	query := map[string]interface{}{
		"template_id": templateQuery["id"],
	}
	attributes, err := s.productRepository.GetAllSelectedAttribute(query, p)
	if err != nil {
		return nil, err
	}
	var AttributeKeyValues []TemplateSelectedAttributesResponse
	dto, err = json.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &AttributeKeyValues)
	if err != nil {
		return nil, err
	}

	for index, attribute := range AttributeKeyValues {
		query = map[string]interface{}{
			"template_id":  templateQuery["id"],
			"attribute_id": attribute.AttributeId,
		}
		attribute_values, err := s.productRepository.GetAllSelectedAttributeValue(query, p)
		if err != nil {
			return nil, err
		}
		var AttributeValues []TemplateSelectedAttributeValuesResponse
		dto, err = json.Marshal(attribute_values)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dto, &AttributeValues)
		if err != nil {
			return nil, err
		}
		AttributeKeyValues[index].AttributeValues = append(AttributeKeyValues[index].AttributeValues, AttributeValues...)
	}
	response.AttributeValues = AttributeKeyValues

	//-------------------Price List Details ---------------------------------
	template_price_list_details := []map[string]interface{}{}
	if response.VendorPriceListIds != nil {
		for _, price_list_id := range response.VendorPriceListIds {
			query := map[string]interface{}{"id": price_list_id}
			vendorPriceList, err := s.vendorRepository.FindOnePriceList(query)
			if err != nil {
				return nil, err
			}

			var priceListInterface map[string]interface{}
			mdata, err := json.Marshal(vendorPriceList)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(mdata, &priceListInterface)
			if err != nil {
				return nil, err
			}
			template_price_list_details = append(template_price_list_details, priceListInterface)
		}
	}
	response.PriceListDetails = template_price_list_details
	return response, err
}
func (s *service) GetAllTemplate(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	result, err := s.productRepository.GetAllTemplate(p)
	if err != nil {
		return nil, err
	}

	var response []TemplateReponseDTO
	mdata, _ := json.Marshal(result)
	json.Unmarshal(mdata, &response)

	return response, nil
}

// ------------------Product Variant------------------------------------------------------------------------------------------------------------------
func (s *service) CreateVariant(data CreateProductVariantDTO, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	var Variant mdm.ProductVariant
	mdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(mdata, &Variant)
	if err != nil {
		return err
	}

	if Variant.ProductTemplateId == 0 || Variant.ParentSkuId == "" || Variant.SkuId == "" {
		return errors.New("template_id, parent_sku_id, and sku_id are required")
	}

	//===============create varaint===============
	err = s.productRepository.CreateVariant(&Variant)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateVariant(data CreateProductVariantDTO, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	var Variant mdm.ProductVariant
	mdata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(mdata, &Variant)
	if err != nil {
		return err
	}

	//===============update variant=========================
	err = s.productRepository.UpdateVariant(&Variant, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteVariant(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteVariant(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetVariantView(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view products at data level")
	}
	p := new(pagination.Paginatevalue)

	//=============find variants =================
	res, err := s.productRepository.FindVariant(query)
	if err != nil {
		return nil, err
	}

	var response VariantResponseDTO
	dto, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &response)
	if err != nil {
		return nil, err
	}

	//===============attributes section==================
	attribute_details := []map[string]interface{}{}
	if response.AttributeKeyValuesId != nil {
		for _, attribute_value_id := range response.AttributeKeyValuesId {
			query = map[string]interface{}{"id": uint(attribute_value_id)}
			attribute_value, err := s.productRepository.GetAllBaseAttributeValue(query, p)
			if err != nil {
				return nil, err
			}
			query = map[string]interface{}{"id": attribute_value[0].BaseAttributeId}
			attribute, err := s.productRepository.GetAllBaseAttribute(query, p)
			if err != nil {
				return nil, err
			}
			attribute_details = append(attribute_details, map[string]interface{}{
				"attribute_name":  attribute[0].Name,
				"attribute_value": attribute_value[0].Name,
			})
		}
	}
	response.AttributeValues = attribute_details

	//===============Price List Details===================
	price_list_details := []map[string]interface{}{}
	if response.VendorPriceListIds != nil {
		for _, price_list_id := range response.VendorPriceListIds {
			query := map[string]interface{}{"id": price_list_id}
			vendorPriceList, err := s.vendorRepository.FindOnePriceList(query)
			if err != nil {
				return nil, err
			}

			var priceListInterface map[string]interface{}
			mdata, err := json.Marshal(vendorPriceList)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(mdata, &priceListInterface)
			if err != nil {
				return nil, err
			}
			price_list_details = append(price_list_details, priceListInterface)
		}
	}
	response.PriceListDetails = price_list_details

	return response, err
}
func (s *service) GetAllVariant(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	result, err := s.productRepository.GetAllVariant(p)
	if err != nil {
		return nil, err
	}

	var response []VariantResponseDTO
	mdata, _ := json.Marshal(result)
	json.Unmarshal(mdata, &response)

	for index, value := range response {
		query := map[string]interface{}{
			"id": value.ProductTemplateId,
		}
		template, _ := s.productRepository.FindTemplate(query)
		response[index].BrandID = template.BrandID
		helpers.JsonMarshaller(template.Brand, &response[index].Brand)
		response[index].HSNCODE = template.HSNCODE
		response[index].HSNCodesData = template.HSNCodesData
	}

	return response, nil
}
func (s *service) Qrcode(id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for view products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for view products at data level")
	}
	query := map[string]interface{}{
		"id": id,
	}
	res, err := s.productRepository.FindVariant(query)
	if err != nil {
		return err
	}
	content := fmt.Sprintf("Product Variant Name  - %v", res.ProductName)
	err = qrcode.WriteFile(content, qrcode.Medium, 256, "Product.png")
	return err
}
func (s *service) Skucodecheck(skucode string, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for view products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for view products at data level")
	}
	query := map[string]interface{}{
		"sku_code": skucode,
	}
	q := map[string]interface{}{
		"sku_id": skucode,
	}
	_, er := s.productRepository.FindTemplate(query)
	if er != nil {
		_, err := s.productRepository.FindVariant(q)
		if err != nil {
			return err
		}

	}
	return nil
}

// ---------------------Product Bundles---------------------------------------------------------------------------------------------------------------
func (s *service) CreateBundle(data CreateBundlePayload, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create products at data level")
	}
	var Bundle mdm.ProductBundles
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Bundle)
	if err != nil {
		return nil, err
	}
	err = s.productRepository.CreateBundle(&Bundle)
	if err != nil {
		return nil, err
	}
	return Bundle, nil
}
func (s *service) UpdateBundle(data CreateBundlePayload, query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update products at data level")
	}
	var Bundle mdm.ProductBundles
	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Bundle)
	if err != nil {
		return err
	}
	err = s.productRepository.UpdateBundle(&Bundle, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteBundle(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete products at data level")
	}
	err := s.productRepository.DeleteBundle(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetOneBundle(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view products at data level")
	}
	res, err := s.productRepository.FindBundle(query)
	if err != nil {
		return nil, err
	}
	var response BudleResponseDTO
	dto, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &response)
	if err != nil {
		return nil, err
	}
	for index, bundle_line_item := range response.Products {
		var variant VariantDetailsForBundleDTO
		dto, err := json.Marshal(bundle_line_item.ProductVariant)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(dto, &variant)
		if err != nil {
			return nil, err
		}
		response.Products[index].ProductVariant = variant
	}
	return response, err
}
func (s *service) GetAllBundle(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]mdm.ProductBundles, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list products at data level")
	}
	result, err := s.productRepository.GetAllBundle(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ------------------Product BUlk----------------------------------------------------------------------------------------------------------------------
func (s *service) BulkCreateProduct(data *[]mdm.ProductTemplate, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create products at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create products at data level")
	}
	err := s.productRepository.BulkCreateProduct(data)
	if err != nil {
		return err
	}
	return nil
}

// ---------------------Product UpsertAPi's----------------------------------------------------------------------------------------------------------------------
func (s *service) UpsertProductTemplate(data []interface{}, meta_data map[string]interface{}) (interface{}, error) {
	var success []interface{}
	var failures []interface{}

	token_id := helpers.ConvertStringToUint(meta_data["token_id"].(string))

	for index, payload := range data {
		product_payload := payload.(map[string]interface{})["product"]
		pricing_payload := payload.(map[string]interface{})["pricing"]
		inventory_payload := payload.(map[string]interface{})["inventory"]

		var templatePayload CreateProductTemplatePayload
		var Template mdm.ProductTemplate
		var channelProductLink mdm.ProductToChannelMap

		jsonPayloadByteArray, err := json.MarshalIndent(product_payload, "", "\t")
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "cannot proceed with this payload"})
			continue
		}
		err = json.Unmarshal(jsonPayloadByteArray, &templatePayload)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "cannot proceed with this payload"})
			continue
		}

		templatePayload.UpdatedByID = token_id

		dto, err := json.Marshal(templatePayload)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "cannot proceed with this payload"})
			continue
		}
		err = json.Unmarshal(dto, &Template)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "cannot proceed with this payload"})
			continue
		}

		err = json.Unmarshal(dto, &channelProductLink)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "cannot proceed with this payload"})
			continue
		}

		channelProductLink.TemplateSku = Template.SkuCode

		if Template.SkuCode == "" || channelProductLink.ChannelCode == "" {
			failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": "sku_code & channel_code should not be empty"})
			continue
		}

		channelProductLinkQuery := map[string]interface{}{"template_sku": Template.SkuCode, "product_variant_id": 0, "channel_code": channelProductLink.ChannelCode}

		res, _ := s.productRepository.FindChannelProductMap(channelProductLinkQuery)
		if res.TemplateSku != "" {
			//=============== update ===========================================
			templateQuery := map[string]interface{}{"sku_code": Template.SkuCode}
			err = s.productRepository.UpdateChannelProductMap(&channelProductLink, channelProductLinkQuery)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			Template.ProductVariantIds = nil
			err = s.productRepository.UpdateTemplate(&Template, templateQuery)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}

			var variantdata []interface{}
			variants := product_payload.(map[string]interface{})["product_variant_ids"]
			if variants != nil && len(variants.([]interface{})) > 0 {
				mdata, _ := json.Marshal(variants)
				json.Unmarshal(mdata, &variantdata)
			}

			//===================default variant for template==============================================

			defaultVariant := map[string]interface{}{}
			defaultVariant["sku_id"] = templatePayload.SKUCode
			DefaultProductVariant := map[string]interface{}{
				"product":   defaultVariant,
				"pricing":   pricing_payload,
				"inventory": inventory_payload,
			}
			if len(variantdata) == 0 {
				variantdata = append(variantdata, DefaultProductVariant)
			}

			Template, _ := s.productRepository.FindTemplate(templateQuery)
			for index := range variantdata {
				var variant CreateProductVariantDTO
				helpers.JsonMarshaller(variantdata[index].(map[string]interface{})["product"], &variant)
				variant.ParentSKUId = Template.SkuCode
				variant.ChannelCode = templatePayload.ChannelCode
				variant.ProductTemplateId = Template.ID
				variant.StatusId = Template.StatusId
				var images []interface{}
				helpers.JsonMarshaller(variant.ImageOptions, &images)
				if variant.ImageOptions == nil || len(images) == 0 {
					variant.ImageOptions = Template.ImageOptions
				}
				if variant.Description == nil {
					variant.Description = Template.Description
				}
				variant.CategoryID = Template.PrimaryCategoryID
				variant.LeafCategoryID = Template.SecondaryCategoryID
				variant.ShippingOptions = Template.ShippingOptions
				variant.PackageDimensions = Template.PackageDimensions
				variant.PackageMaterialOptions = Template.PackageMaterialOptions
				variant.ProductPricingDetails.CostPrice = Template.ProductPricingDetails.CostPrice
				variant.ProductPricingDetails.MRP = Template.ProductPricingDetails.MRP
				if variant.ProductPricingDetails.SalesPrice == 0 {
					variant.ProductPricingDetails.SalesPrice = Template.ProductPricingDetails.SalesPrice
				}

				variantdata[index].(map[string]interface{})["product"] = variant
			}

			res, err := s.UpsertProductVariant(variantdata, meta_data)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": "product_variant_error", "product_variants": err})
				continue
			}

			success = append(success, map[string]interface{}{
				"sku_code":         Template.SkuCode,
				"status":           true,
				"serial_number":    index + 1,
				"msg":              "updated",
				"product_variants": res,
			})
			continue
		}

		//=========== create ================================
		var variantdata []interface{}
		variants := product_payload.(map[string]interface{})["product_variant_ids"]
		if variants != nil && len(variants.([]interface{})) > 0 {
			helpers.JsonMarshaller(variants, &variantdata)
		}

		defaultStatus, err := helpers.GetLookupcodeId("PRODUCT_STATUS", "DRAFT")
		if err != nil {
			return nil, err
		}
		if templatePayload.StatusId == nil {
			templatePayload.StatusId = &defaultStatus
		}

		//===================default variant for template==============================================
		defaultVariant := map[string]interface{}{}
		defaultVariant["sku_id"] = templatePayload.SKUCode
		DefaultProductVariant := map[string]interface{}{
			"product":   defaultVariant,
			"pricing":   pricing_payload,
			"inventory": inventory_payload,
		}
		if len(variantdata) == 0 {
			variantdata = append(variantdata, DefaultProductVariant)
		}

		//============create or update product template===========================
		Template.ProductVariantIds = nil
		Template.CreatedByID = token_id
		TemplateFindQuery := map[string]interface{}{"sku_code": Template.SkuCode}
		TemplateResponse, _ := s.productRepository.FindTemplate(TemplateFindQuery)
		if TemplateResponse.SkuCode != "" {
			err = s.productRepository.UpdateTemplate(&Template, TemplateFindQuery)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}
		} else {
			err = s.productRepository.CreateTemplate(&Template)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}
		}

		//===============create product_channel_map===================
		channelProductLink.ProductTemplateId = Template.ID
		channelProductLink.CreatedByID = token_id
		err = s.productRepository.CreateChannelProductMap(&channelProductLink)
		if err != nil {
			failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		for index := range variantdata {
			var variant CreateProductVariantDTO
			helpers.JsonMarshaller(variantdata[index].(map[string]interface{})["product"], &variant)
			variant.ParentSKUId = Template.SkuCode
			variant.ChannelCode = templatePayload.ChannelCode
			variant.ProductTemplateId = Template.ID
			variant.StatusId = Template.StatusId
			var images []interface{}
			helpers.JsonMarshaller(variant.ImageOptions, &images)
			if variant.ImageOptions == nil || len(images) == 0 {
				variant.ImageOptions = Template.ImageOptions
			}
			if variant.Description == nil {
				variant.Description = Template.Description
			}
			variant.CategoryID = Template.PrimaryCategoryID
			variant.LeafCategoryID = Template.SecondaryCategoryID
			variant.ShippingOptions = Template.ShippingOptions
			variant.PackageDimensions = Template.PackageDimensions
			variant.PackageMaterialOptions = Template.PackageMaterialOptions
			variant.ProductPricingDetails.CostPrice = Template.ProductPricingDetails.CostPrice
			variant.ProductPricingDetails.MRP = Template.ProductPricingDetails.MRP
			if variant.ProductPricingDetails.SalesPrice == 0 {
				variant.ProductPricingDetails.SalesPrice = Template.ProductPricingDetails.SalesPrice
			}
			variantdata[index].(map[string]interface{})["product"] = variant
		}

		var jsonvariantPayload []interface{}
		helpers.JsonMarshaller(variantdata, &jsonvariantPayload)

		variant_response, err := s.UpsertProductVariant(jsonvariantPayload, meta_data)
		if err != nil {
			failures = append(failures, map[string]interface{}{"sku_code": Template.SkuCode, "status": false, "serial_number": index + 1, "msg": "product_variant_error", "product_variants": err})
			continue
		}

		success = append(success, map[string]interface{}{
			"sku_code":         Template.SkuCode,
			"status":           true,
			"serial_number":    index + 1,
			"msg":              "created",
			"product_variants": variant_response,
		})
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
}
func (s *service) UpsertProductVariant(data []interface{}, meta_data map[string]interface{}) (interface{}, error) {

	var success []interface{}
	var failures []interface{}
	token_id := helpers.ConvertStringToUint(meta_data["token_id"].(string))

	for index, payload := range data {
		product_payload := payload.(map[string]interface{})["product"]
		pricing_payload := payload.(map[string]interface{})["pricing"]
		inventory_payload := payload.(map[string]interface{})["inventory"]

		var Variant mdm.ProductVariant
		var channelProductLink mdm.ProductToChannelMap

		//=============== marshal payload to both Variant & ChannelProductLink ===================================================
		err := helpers.JsonMarshaller(product_payload, &Variant)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "cannot proceed with this payload", "error": err})
			continue
		}
		err = helpers.JsonMarshaller(product_payload, &channelProductLink)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "cannot proceed with this payload", "error": err})
			continue
		}

		channelProductLink.TemplateSku = Variant.ParentSkuId
		channelProductLink.ProductVariantSku = Variant.SkuId

		//================= check template and variant sku are available or not ===================================================
		if Variant.ParentSkuId == "" || Variant.SkuId == "" || channelProductLink.ChannelCode == "" {
			failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": "parent_sku_id and variant sku_id and channel_code are required"})
			continue
		}

		//================ find the product_template from product_to_channel_map using channel_product_link_query =============================
		channelProductLinkQuery := map[string]interface{}{"product_template_id": channelProductLink.ProductTemplateId, "channel_code": channelProductLink.ChannelCode, "product_variant_id": 0}
		res, _ := s.productRepository.FindChannelProductMap(channelProductLinkQuery)
		if res.TemplateSku == "" {
			failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": "Parent product is not available"})
			continue
		}

		var pricingResponse interface{}
		var inventoryResponse interface{}

		channelProductLinkQuery = map[string]interface{}{"product_variant_sku": channelProductLink.ProductVariantSku, "channel_code": channelProductLink.ChannelCode}
		res, _ = s.productRepository.FindChannelProductMap(channelProductLinkQuery)
		if Variant.ProductName == "" {
			query := map[string]interface{}{"id": Variant.ProductTemplateId}
			result, err := s.productRepository.FindTemplate(query)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			Variant.ProductName = result.Name
		}
		if res.ProductVariantSku != "" {
			Variant.UpdatedByID = token_id
			err = s.productRepository.UpdateChannelProductMap(&channelProductLink, channelProductLinkQuery)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			VariantQuery := map[string]interface{}{"sku_id": Variant.SkuId}
			err = s.productRepository.UpdateVariant(&Variant, VariantQuery)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": err})
				continue
			}

			findQuery := map[string]interface{}{"sku_id": Variant.SkuId}
			VariantInformation, _ := s.productRepository.FindVariant(findQuery)

			pricingResponse, inventoryResponse, err = s.UpsertPricingAndInventory(VariantInformation, pricing_payload, inventory_payload, meta_data)
			if err != nil {
				failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": err})
			}
			success = append(success, map[string]interface{}{
				"sku_id":        Variant.SkuId,
				"status":        true,
				"serial_number": index + 1,
				"msg":           "updated",
				"pricing":       pricingResponse,
				"inventory":     inventoryResponse,
			})
			continue
		}

		//====================== create product variant ========================================
		Variant.CreatedByID = token_id
		err = s.productRepository.CreateVariant(&Variant)
		if err != nil {
			failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": err})
			continue
		}

		//===================== create product_to_channel_map ===================================
		channelProductLink.CreatedByID = token_id
		channelProductLink.ProductVariantId = Variant.ID
		err = s.productRepository.CreateChannelProductMap(&channelProductLink)
		if err != nil {
			failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": err})
			continue
		}

		pricingResponse, inventoryResponse, err = s.UpsertPricingAndInventory(Variant, pricing_payload, inventory_payload, meta_data)
		if err != nil {
			failures = append(failures, map[string]interface{}{"sku_id": Variant.SkuId, "status": false, "serial_number": index + 1, "msg": err})
		}
		success = append(success, map[string]interface{}{
			"sku_id":        Variant.SkuId,
			"status":        true,
			"serial_number": index + 1,
			"msg":           "created",
			"pricing":       pricingResponse,
			"inventory":     inventoryResponse,
		})
	}

	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
}
func (s *service) UpsertPricingAndInventory(VariantInformation mdm.ProductVariant, pricing_payload interface{}, inventory_payload interface{}, meta_data map[string]interface{}) (interface{}, interface{}, error) {

	var metaData model_core.MetaData
	pricingPayload := map[string]interface{}{}
	inventoryPayload := map[string]interface{}{}
	salesPriceList := map[string]interface{}{}
	pricingPayloadFlag := true
	var ok bool

	//===============Pricing=================================================================================
	if pricing_payload == nil {
		pricingPayloadFlag = false
	}

	if pricingPayloadFlag {
		pricingPayload, ok = pricing_payload.(map[string]interface{})
		if !ok {
			return nil, nil, errors.New("pricingPayload - cannot convert payload to interface{} to map[string]interface{}")
		}

		if pricingPayload == nil || pricingPayload["price_list_id"] == nil {
			pricingPayloadFlag = false
		}
	}

	if pricingPayloadFlag {
		salesPriceList, ok = pricingPayload["sales_price_list"].(map[string]interface{})
		if !ok {
			return nil, nil, errors.New("sales_price_list - cannot convert payload to interface{} to map[string]interface{}")
		}
		if salesPriceList == nil || salesPriceList["sales_line_items"] == nil {
			pricingPayloadFlag = false
		}
	}

	if !pricingPayloadFlag {
		default_sales_line_item := map[string]interface{}{
			"product_id": VariantInformation.ID,
		}
		default_pricing_payload := map[string]interface{}{
			"price_list_id": 1,
			"sales_price_list": map[string]interface{}{
				"sales_line_items": []map[string]interface{}{
					default_sales_line_item,
				},
			},
		}
		_, ok = pricing_payload.(map[string]interface{})
		if !ok {
			pricing_payload = map[string]interface{}{}
		}

		for i := range default_pricing_payload {
			pricing_payload.(map[string]interface{})[i] = default_pricing_payload[i]
		}

	} else {
		if salesPriceList["sales_line_items"] != nil {
			lineItems, ok := salesPriceList["sales_line_items"].([]interface{})
			if !ok {
				return nil, nil, errors.New("sales_line_items is not an array")
			}
			for _, item := range lineItems {
				item.(map[string]interface{})["product_id"] = VariantInformation.ID
			}
		}
		pricingPayload["sales_price_list"] = salesPriceList
	}
	var pricing_request_payload_data []interface{}
	pricing_request_payload_data = append(pricing_request_payload_data, pricingPayload)
	pricing_request_payload := map[string]interface{}{
		"meta_data": meta_data,
		"data":      pricing_request_payload_data,
	}
	eda.Produce(eda.UPSERT_SALES_PRICE_LIST, pricing_request_payload)

	//=============== Inventory =================================================================================
	if inventory_payload != nil {
		inventoryPayload = inventory_payload.(map[string]interface{})
	}
	inventoryPayload["product_variant_id"] = VariantInformation.ID

	virtual_type_id, _ := helpers.GetLookupcodeId("LOCATION_TYPE", "VIRTUAL_LOCATION")

	p := new(pagination.Paginatevalue)
	p.Filters = fmt.Sprintf("[[\"name\",\"=\",\"Default Location\"],[\"location_type_id\",\"=\",\"%v\"],[\"company_id\",\"=\",\"%v\"]]", virtual_type_id, meta_data["company_id"])
	metaData.ModuleAccessAction = "LIST"
	location_data, _ := s.location_service.GetLocationList(metaData, p)
	if location_data != nil {
		location_arr := location_data.([]location_service.LocationListResponseDTO)
		if len(location_arr) > 0 {
			inventoryPayload["physical_location_id"] = location_arr[0].ID
		}
	}

	var inventory_request_payload_data []interface{}
	inventory_request_payload_data = append(inventory_request_payload_data, inventoryPayload)
	inventory_request_payload := map[string]interface{}{
		"meta_data": meta_data,
		"data":      inventory_request_payload_data,
	}
	eda.Produce(eda.UPSERT_DECENTRALIZED_INVENTORY, inventory_request_payload)

	return "Inprogress", "Inprogress", nil
}

func (s *service) GetAllChannelProducts(p *pagination.Paginatevalue, queryParams QueryParamsDTO, token_id string, access_template_id string, access_action string) ([]ChannelMapResponseDTO, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRODUCTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list channel products at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list channel products at data level")
	}
	var ChannelMapResponse []ChannelMapResponseDTO

	if queryParams.Ids != "" && queryParams.Ids != "[]" {
		ids := queryParams.Ids
		length := len(ids)
		ids = ids[1 : length-1]

		result, err := s.productRepository.GetAllChannelProductMapInArray(queryParams.ColumnName, queryParams.Condition, ids)
		if err != nil {
			return nil, err
		}
		mdata, _ := json.Marshal(result)
		json.Unmarshal(mdata, &ChannelMapResponse)

	} else {
		result, err := s.productRepository.GetAllChannelProductMap(p)
		if err != nil {
			return nil, err
		}

		mdata, _ := json.Marshal(result)
		json.Unmarshal(mdata, &ChannelMapResponse)
	}

	//---------get_product_inventory------------
	if queryParams.Inventory {
		for index, channel := range ChannelMapResponse {
			var inventoryResponse ChannelInventoryResponseDTO
			if channel.ProductVariantId != 0 {
				findQuery := map[string]interface{}{
					"product_variant_id": channel.ProductVariantId,
					"channel_code":       channel.ChannelCode,
				}
				result, err := s.inventoryRepository.FindOneDecentralizedInventory(findQuery)
				if err != nil {
					return nil, err
				}

				mdata, _ := json.Marshal(result)
				json.Unmarshal(mdata, &inventoryResponse)

				ChannelMapResponse[index].Inventory = inventoryResponse
				continue
			}
			if channel.ProductTemplateId != 0 {
				findQuery := map[string]interface{}{
					"template_sku":        channel.TemplateSku,
					"product_variant_sku": channel.TemplateSku,
				}
				result, err := s.productRepository.GetChannelProductMap(findQuery)
				if err != nil {
					return nil, err
				}
				findQuery = map[string]interface{}{
					"product_variant_id": result.ProductVariantId,
					"channel_code":       channel.ChannelCode,
				}
				resonse, err := s.inventoryRepository.FindOneDecentralizedInventory(findQuery)
				if err != nil {
					return nil, err
				}

				mdata, _ := json.Marshal(resonse)
				json.Unmarshal(mdata, &inventoryResponse)

				ChannelMapResponse[index].Inventory = inventoryResponse
			}

		}
	}
	//---------get_product_pricing------------
	if queryParams.Pricing {
		for index, channel := range ChannelMapResponse {
			var pricingResponse ChannelPricingResponseDTO
			if channel.ChannelCode == "" {
				ChannelMapResponse[index].Pricing = nil
				continue
			}
			if channel.ProductVariantId != 0 {

				findQuery := map[string]interface{}{
					"price_list_name": channel.ChannelCode,
				}
				result, err := s.pricingRepository.FindPricing(findQuery)
				if err != nil {
					return nil, err
				}

				findQuery = map[string]interface{}{
					"spl_id":     result.SalesPriceListId,
					"product_id": channel.ProductVariantId,
				}

				response, err := s.pricingRepository.FindOneSalesPriceLineItem(findQuery)
				if err != nil {
					return nil, err
				}

				mdata, _ := json.Marshal(response)
				json.Unmarshal(mdata, &pricingResponse)

				ChannelMapResponse[index].Pricing = pricingResponse
				continue
			}
			if channel.ProductTemplateId != 0 {
				findQuery := map[string]interface{}{
					"template_sku":        channel.TemplateSku,
					"product_variant_sku": channel.TemplateSku,
				}
				result, err := s.productRepository.GetChannelProductMap(findQuery)
				if err != nil {
					return nil, err
				}

				findQuery = map[string]interface{}{
					"price_list_name": result.ChannelCode,
				}
				pricingresult, err := s.pricingRepository.FindPricing(findQuery)
				if err != nil {
					return nil, err
				}

				findQuery = map[string]interface{}{
					"spl_id":     pricingresult.SalesPriceListId,
					"product_id": result.ProductVariantId,
				}

				response, err := s.pricingRepository.FindOneSalesPriceLineItem(findQuery)
				if err != nil {
					return nil, err
				}

				mdata, _ := json.Marshal(response)
				json.Unmarshal(mdata, &pricingResponse)

				ChannelMapResponse[index].Pricing = pricingResponse
			}
		}
	}
	//---------get_product_product_template------------
	if queryParams.ProductTemplate {
		for index, channel := range ChannelMapResponse {
			findQuery := map[string]interface{}{
				"sku_code": channel.TemplateSku,
				"id":       channel.ProductTemplateId,
			}

			result, err := s.GetTemplateView(findQuery, token_id, access_template_id)
			if err != nil {
				return nil, err
			}

			ChannelMapResponse[index].ProductTemplateDetails = result
		}
	}
	//---------get_product_variant------------
	if queryParams.ProductVariant {
		for index, channel := range ChannelMapResponse {
			findQuery := map[string]interface{}{
				"sku_id": channel.ProductVariantSku,
				"id":     channel.ProductVariantId,
			}
			if channel.ProductVariantId == 0 {
				findQuery = map[string]interface{}{
					"sku_id": channel.TemplateSku,
				}
			}

			result, err := s.GetVariantView(findQuery, token_id, access_template_id)
			if err != nil {
				return nil, err
			}

			ChannelMapResponse[index].ProductVariantDetails = result
		}
	}

	return ChannelMapResponse, nil
}

//-----------------------Filter Tabs----------------------------------------------------------------------------------------

func (s *service) GetProductVariantTab(product_variant_id string, tab_name string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	// var metaData core_model.MetaData
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "SALES_ORDERS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view product variant at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view product variant at data level")
	}

	_, err := strconv.Atoi(product_variant_id)
	if err != nil {
		return nil, err
	}

	if tab_name == "price_list" {
		var priceListTabResponse PriceListTabResponse
		var salesPriceListTabResponse []SalesPriceListTab
		var purchasePriceListTabResponse []PurchasePriceListTab
		query := map[string]interface{}{}

		helpers.AddMandatoryFilters(page, "product_variant_id", "=", product_variant_id)
		page.Per_page = 50
		//--------------------sales line items---------------------------
		// page.Filters = "[[\"product_variant_id\",\"=\",\"" + product_variant_id + "\"]]"

		salesLineItemResponse, err := s.pricingRepository.ListSalesPiceLineItems(query, page)
		if err != nil {
			return nil, err
		}
		jsondata, _ := json.Marshal(salesLineItemResponse)
		json.Unmarshal(jsondata, &salesPriceListTabResponse)

		for index, lineItem := range salesPriceListTabResponse {
			query := map[string]interface{}{
				"sales_price_list_id": lineItem.SPL_id,
			}
			res, err := s.pricingRepository.FindPricing(query)
			if err != nil {
				return nil, err
			}
			salesPriceListTabResponse[index].Price_list_name = res.PriceListName
			salesPriceListTabResponse[index].Add_channel_of_sale_id = res.SalesPriceList.AddChannelOfSaleId
			jsondata, _ := json.Marshal(res.SalesPriceList.AddChannelOfSale)
			var channel_of_sale core.LookupCodesDTO
			json.Unmarshal(jsondata, &channel_of_sale)
			salesPriceListTabResponse[index].Add_channel_of_sale = channel_of_sale
		}

		//------------------purchase line Items------------------------------
		// page.Filters = "[[\"product_id\",\"=\",\"" + product_variant_id + "\"]]"
		// page.Per_page = 50

		PurchaseLineItemResponse, err := s.pricingRepository.ListPurchasePriceListLineItems(query, page)
		if err != nil {
			return nil, err
		}
		jsondata, _ = json.MarshalIndent(PurchaseLineItemResponse, " ", "\t")
		json.Unmarshal(jsondata, &purchasePriceListTabResponse)

		for index, lineItem := range purchasePriceListTabResponse {
			query := map[string]interface{}{
				"purchase_price_list_id": lineItem.PPL_id,
			}
			res, err := s.pricingRepository.FindPricing(query)
			if err != nil {
				return nil, err
			}
			purchasePriceListTabResponse[index].Price_list_name = res.PriceListName
			purchasePriceListTabResponse[index].VendorName = res.PurchasePriceList.VendorDetails.Name
		}
		priceListTabResponse.SalesPriceList = salesPriceListTabResponse
		priceListTabResponse.PurchasePriceList = purchasePriceListTabResponse

		return priceListTabResponse, nil
	}

	if tab_name == "bundles" {
		var bundleTabResponse = make([]ProductBundlesTabResponse, 0)
		page.Filters = "[[\"product_variant_id\",\"=\",\"" + product_variant_id + "\"]]"
		page.Per_page = 50

		BundleLineItemResponse, err := s.productRepository.ListBundleLineItems(page)
		if err != nil {
			return nil, err
		}

		for _, lineItem := range BundleLineItemResponse {
			query := map[string]interface{}{
				"id": lineItem.BundleId,
			}
			res, err := s.productRepository.FindBundle(query)
			if err != nil {
				return nil, err
			}
			var bundleResponse ProductBundlesTabResponse
			jsondata, _ := json.MarshalIndent(res, " ", "\t")
			json.Unmarshal(jsondata, &bundleResponse)
			bundleResponse.Quantity = lineItem.Quantity
			bundleTabResponse = append(bundleTabResponse, bundleResponse)
		}
		return bundleTabResponse, nil
	}
	return nil, nil
}

func (s *service) CreateHsn(data mdm.HSNCodesData) (map[string]interface{}, error) {
	err := s.productRepository.SaveHsn(data)
	if err != nil {
		return map[string]interface{}{"success": "false"}, err
	}
	return map[string]interface{}{"success": "true"}, nil
}

func (s *service) GetAllHsn(page *pagination.Paginatevalue) ([]mdm.HSNCodesData, error) {
	result, err := s.productRepository.FindAllHsn(page)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *service) GetOneHsn(query map[string]interface{}) (mdm.HSNCodesData, error) {
	result, err := s.productRepository.FindOneHsn(query)
	if err != nil {
		return result, err
	}
	return result, nil

}
func (s *service) UpdateHsn(query map[string]interface{}, data mdm.HSNCodesData) error {
	err := s.productRepository.UpdateHsn(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteHsn(query map[string]interface{}) error {
	err := s.productRepository.DeleteHsn(query)
	if err != nil {
		return err
	}
	return nil
}
