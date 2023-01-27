package mdm

import (
	"errors"
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
type Products interface {

	//--------------------Product Brand------------------------------------------------------------------------------------------------------------------
	CreateBrand(data *mdm.ProductBrand) error
	UpdateBrand(data *mdm.ProductBrand, query map[string]interface{}) error
	DeleteBrand(query map[string]interface{}) error
	GetAllBrand(p *pagination.Paginatevalue) ([]mdm.ProductBrand, error)
	SearchBrand(query string) ([]mdm.ProductBrand, error)

	//------------------Product Category-----------------------------------------------------------------------------------------------------------------
	CreateCategory(data *mdm.ProductCategory) error
	UpdateCategory(data *mdm.ProductCategory, query map[string]interface{}) error
	DeleteCategory(query map[string]interface{}) error
	GetAllCategory(p *pagination.Paginatevalue) ([]mdm.ProductCategory, error)
	GetAllSubCategory(p *pagination.Paginatevalue) ([]mdm.ProductCategory, error)
	SearchCategory(query string) ([]mdm.ProductCategory, error)

	//------------------Product Base Attributes----------------------------------------------------------------------------------------------------------
	CreateBaseAttribute(data *mdm.ProductBaseAttributes) error
	UpdateBaseAttribute(data *mdm.ProductBaseAttributes, query map[string]interface{}) error
	DeleteBaseAttribute(query map[string]interface{}) error
	GetAllBaseAttribute(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductBaseAttributes, error)

	//------------------Product Base Attributes Values----------------------------------------------------------------------------------------------------
	CreateBaseAttributeValue(data *mdm.ProductBaseAttributesValues) error
	UpdateBaseAttributeValue(data *mdm.ProductBaseAttributesValues, query map[string]interface{}) error
	DeleteBaseAttributeValue(query map[string]interface{}) error
	GetAllBaseAttributeValue(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductBaseAttributesValues, error)

	//------------------Product Selected Base Attributes--------------------------------------------------------------------------------------------------
	CreateSelectedAttribute(data *mdm.ProductSelectedAttributes) error
	UpdateSelectedAttribute(data *mdm.ProductSelectedAttributes, query map[string]interface{}) error
	DeleteSelectedAttribute(query map[string]interface{}) error
	GetAllSelectedAttribute(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductSelectedAttributes, error)

	//------------------Product Selected Base Attributes Values--------------------------------------------------------------------------------------------
	CreateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues) error
	UpdateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues, query map[string]interface{}) error
	DeleteSelectedAttributeValue(query map[string]interface{}) error
	GetAllSelectedAttributeValue(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductSelectedAttributesValues, error)

	//---------------------Product Template----------------------------------------------------------------------------------------------------------------
	CreateTemplate(data *mdm.ProductTemplate) error
	UpdateTemplate(data *mdm.ProductTemplate, query map[string]interface{}) error
	DeleteTemplate(query map[string]interface{}) error
	FindTemplate(query map[string]interface{}) (mdm.ProductTemplate, error)
	GetAllTemplate(p *pagination.Paginatevalue) ([]mdm.ProductTemplate, error)

	//---------------------Product Variant-----------------------------------------------------------------------------------------------------------------
	CreateVariant(data *mdm.ProductVariant) error
	UpdateVariant(data *mdm.ProductVariant, query map[string]interface{}) error
	DeleteVariant(query map[string]interface{}) error
	FindVariant(query map[string]interface{}) (mdm.ProductVariant, error)
	GetAllVariant(p *pagination.Paginatevalue) ([]mdm.ProductVariant, error)

	//-------------------Product Bundles------------------------------------------------------------------------------------------------------------------
	CreateBundle(data *mdm.ProductBundles) error
	UpdateBundle(data *mdm.ProductBundles, query map[string]interface{}) error
	DeleteBundle(query map[string]interface{}) error
	FindBundle(query map[string]interface{}) (mdm.ProductBundles, error)
	GetAllBundle(p *pagination.Paginatevalue) ([]mdm.ProductBundles, error)
	ListBundleLineItems(p *pagination.Paginatevalue) ([]mdm.BundleLineItems, error)

	//--------------------Product BUlk--------------------------------------------------------------------------------------------------------------------
	BulkCreateProduct(data *[]mdm.ProductTemplate) error

	//--------------------ProductToChannelMap------------------------------------------------------------------------------------------------------------
	CreateChannelProductMap(data *mdm.ProductToChannelMap) error
	UpdateChannelProductMap(data *mdm.ProductToChannelMap, query map[string]interface{}) error
	DeleteChannelProductMap(query map[string]interface{}) error
	FindChannelProductMap(query map[string]interface{}) (mdm.ProductToChannelMap, error)
	GetAllChannelProductMap(p *pagination.Paginatevalue) ([]mdm.ProductToChannelMap, error)
	GetAllChannelProductMapInArray(column string, condition string, array string) ([]mdm.ProductToChannelMap, error)
	GetChannelProductMap(query interface{}) (mdm.ProductToChannelMap, error)

	//-----------------Product Pricing Details -------------------------------------------------------------------------------------------------------------
	UpdateProductPricingDetails(data *mdm.ProductPricingDetails, query map[string]interface{}) error
	CreateProductPricingDetails(data *mdm.ProductPricingDetails) error

	GetAllProducts(query map[string]interface{}, search_type string) ([]mdm.ProductVariant, error)

	// -----------------------------hsn----------------------------------

	SaveHsn(data mdm.HSNCodesData) error
	FindOneHsn(query map[string]interface{}) (mdm.HSNCodesData, error)
	FindAllHsn(page *pagination.Paginatevalue) ([]mdm.HSNCodesData, error)

	UpdateHsn(query map[string]interface{}, data mdm.HSNCodesData) error

	DeleteHsn(query map[string]interface{}) error
}
type products struct {
	db *gorm.DB
}

var productsRepository *products //singleton object

// singleton function
func NewProducts() *products {
	if productsRepository != nil {
		return productsRepository
	}
	db := db.DbManager()
	productsRepository = &products{db}
	return productsRepository
}

// -------------------------Product Brand--------------------------------------------------------------------------------------------------
func (r *products) CreateBrand(data *mdm.ProductBrand) error {
	err := r.db.Model(&mdm.ProductBrand{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateBrand(data *mdm.ProductBrand, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductBrand{}).Where(query).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteBrand(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductBrand{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) GetAllBrand(p *pagination.Paginatevalue) ([]mdm.ProductBrand, error) {
	var result []mdm.ProductBrand
	err := r.db.Model(&mdm.ProductBrand{}).Scopes(helpers.Paginate(&mdm.ProductBrand{}, p, r.db)).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (r *products) SearchBrand(query string) ([]mdm.ProductBrand, error) {
	var data []mdm.ProductBrand
	err := r.db.Model(&mdm.ProductBrand{}).Find(&data, "name ILIKE ? ", "%"+query+"%").Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// ------------------------Product Category------------------------------------------------------------------------------------------------
func (r *products) CreateCategory(data *mdm.ProductCategory) error {
	err := r.db.Model(&mdm.ProductCategory{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateCategory(data *mdm.ProductCategory, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductCategory{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteCategory(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductCategory{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) GetAllCategory(p *pagination.Paginatevalue) ([]mdm.ProductCategory, error) {
	var result []mdm.ProductCategory
	err := r.db.Model(&mdm.ProductCategory{}).Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.ProductCategory{}, p, r.db)).Find(&result).Error
	return result, err
}
func (r *products) GetAllSubCategory(p *pagination.Paginatevalue) ([]mdm.ProductCategory, error) {
	var result []mdm.ProductCategory
	err := r.db.Model(&mdm.ProductCategory{}).Scopes(helpers.Paginate(&mdm.ProductCategory{}, p, r.db)).Where("parent_category_id is NOT NULL").Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (r *products) SearchCategory(query string) ([]mdm.ProductCategory, error) {
	var data []mdm.ProductCategory
	err := r.db.Model(&mdm.ProductCategory{}).Find(&data, "name ILIKE ? ", "%"+query+"%").Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// ---------------------Product Base Attributes-------------------------------------------------------------------------------------------
func (r *products) CreateBaseAttribute(data *mdm.ProductBaseAttributes) error {
	err := r.db.Model(&mdm.ProductBaseAttributes{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateBaseAttribute(data *mdm.ProductBaseAttributes, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductBaseAttributes{}).Where(query).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteBaseAttribute(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductBaseAttributes{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) GetAllBaseAttribute(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductBaseAttributes, error) {
	var result []mdm.ProductBaseAttributes
	err := r.db.Model(&mdm.ProductBaseAttributes{}).Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.ProductBaseAttributes{}, p, r.db)).Where(query).Find(&result).Error
	return result, err
}

// -------------------Product Base Attributes Values---------------------------------------------------------------------------------------
func (r *products) CreateBaseAttributeValue(data *mdm.ProductBaseAttributesValues) error {
	err := r.db.Model(&mdm.ProductBaseAttributesValues{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateBaseAttributeValue(data *mdm.ProductBaseAttributesValues, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductBaseAttributesValues{}).Where(query).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteBaseAttributeValue(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductBaseAttributesValues{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) GetAllBaseAttributeValue(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductBaseAttributesValues, error) {
	var result []mdm.ProductBaseAttributesValues
	err := r.db.Model(&mdm.ProductBaseAttributesValues{}).Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.ProductBaseAttributesValues{}, p, r.db)).Where(query).Find(&result).Error
	return result, err
}

// --------------------Product Selected Base Attributes--------------------------------------------------------------------------------------
func (r *products) CreateSelectedAttribute(data *mdm.ProductSelectedAttributes) error {
	err := r.db.Model(&mdm.ProductSelectedAttributes{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateSelectedAttribute(data *mdm.ProductSelectedAttributes, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductSelectedAttributes{}).Where(query).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteSelectedAttribute(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductSelectedAttributes{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) GetAllSelectedAttribute(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductSelectedAttributes, error) {
	var result []mdm.ProductSelectedAttributes
	err := r.db.Model(&mdm.ProductSelectedAttributes{}).Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.ProductSelectedAttributes{}, p, r.db)).Where(query).Find(&result).Error
	return result, err
}

// --------------------Product Selected Attributes Values------------------------------------------------------------------------------------
func (r *products) CreateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues) error {
	err := r.db.Model(&mdm.ProductSelectedAttributesValues{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateSelectedAttributeValue(data *mdm.ProductSelectedAttributesValues, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductSelectedAttributesValues{}).Where(query).Updates(&data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteSelectedAttributeValue(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductSelectedAttributesValues{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) GetAllSelectedAttributeValue(query map[string]interface{}, p *pagination.Paginatevalue) ([]mdm.ProductSelectedAttributesValues, error) {
	var result []mdm.ProductSelectedAttributesValues
	err := r.db.Model(&mdm.ProductSelectedAttributesValues{}).Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.ProductSelectedAttributesValues{}, p, r.db)).Where(query).Find(&result).Error
	return result, err
}

// -------------------------Product Template---------------------------------------------------------------------------------------------------
func (r *products) CreateTemplate(data *mdm.ProductTemplate) error {
	err := r.db.Model(&mdm.ProductTemplate{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateTemplate(data *mdm.ProductTemplate, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductTemplate{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteTemplate(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductTemplate{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) FindTemplate(query map[string]interface{}) (mdm.ProductTemplate, error) {
	var result mdm.ProductTemplate
	err := r.db.Model(&mdm.ProductTemplate{}).Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations).Where(query).First(&result).Error
	result.StockTreatmentIds = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(result.StockTreatmentIds)
	result.ProductProcurementTreatmentIds = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(result.ProductProcurementTreatmentIds)
	return result, err
}
func (r *products) GetAllTemplate(p *pagination.Paginatevalue) ([]mdm.ProductTemplate, error) {
	var result []mdm.ProductTemplate
	err := r.db.Model(&mdm.ProductTemplate{}).Preload(clause.Associations + "." + clause.Associations).Scopes(helpers.Paginate(&mdm.ProductTemplate{}, p, r.db)).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ------------------------Product Variant---------------------------------------------------------------------------------------------------------
func (r *products) CreateVariant(data *mdm.ProductVariant) error {
	err := r.db.Model(&mdm.ProductVariant{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateVariant(data *mdm.ProductVariant, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductVariant{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}

	return nil
}
func (r *products) DeleteVariant(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductVariant{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) FindVariant(query map[string]interface{}) (mdm.ProductVariant, error) {
	var result mdm.ProductVariant
	err := r.db.Model(&mdm.ProductVariant{}).Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations).Where(query).First(&result).Error
	result.StandardProductTypes = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(result.StandardProductTypes)
	return result, err
}
func (r *products) GetAllVariant(p *pagination.Paginatevalue) ([]mdm.ProductVariant, error) {
	var result []mdm.ProductVariant
	err := r.db.Model(&mdm.ProductVariant{}).Preload(clause.Associations + "." + clause.Associations).Preload("CreatedBy.Company").Scopes(helpers.Paginate(&mdm.ProductVariant{}, p, r.db)).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *products) GetAllProducts(query map[string]interface{}, search_type string) ([]mdm.ProductVariant, error) {
	var result []mdm.ProductVariant
	var err error
	if search_type == "product_name" {
		err = r.db.Model(&mdm.ProductVariant{}).Find(&result, "product_name ILIKE ? ", "%"+query["name"].(string)+"%").Error
	} else if search_type == "price" {
		err = r.db.Model(&mdm.ProductVariant{}).Where("sales_price BETWEEN ? AND ?", query["minimum_value"], query["maximum_value"]).Find(&result).Error
	} else if search_type == "sku_id" {
		err = r.db.Model(&mdm.ProductVariant{}).Find(&result, "sku_id ILIKE ? ", "%"+query["sku_id"].(string)+"%").Error
	}
	// .Find(&result, "sales_price BETWEEN ? AND ?", "50", "100").Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ----------------------------Product Bundles-------------------------------------------------------------------------------------------------------
func (r *products) CreateBundle(data *mdm.ProductBundles) error {
	err := r.db.Model(&mdm.ProductBundles{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateBundle(data *mdm.ProductBundles, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductBundles{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteBundle(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&mdm.ProductBundles{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) FindBundle(query map[string]interface{}) (mdm.ProductBundles, error) {
	var result mdm.ProductBundles
	err := r.db.Model(&mdm.ProductBundles{}).Preload(clause.Associations + "." + clause.Associations).Preload(clause.Associations).Where(query).First(&result).Error
	return result, err
}
func (r *products) GetAllBundle(p *pagination.Paginatevalue) ([]mdm.ProductBundles, error) {
	var result []mdm.ProductBundles
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Model(&mdm.ProductBundles{}).Scopes(helpers.Paginate(&mdm.ProductBundles{}, p, r.db)).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (r *products) ListBundleLineItems(p *pagination.Paginatevalue) ([]mdm.BundleLineItems, error) {
	var result []mdm.BundleLineItems
	res := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.BundleLineItems{}, p, r.db)).Where("is_active = true").Find(&result)
	return result, res.Error
}

// ------------------Product BUlk--------------------------------------------------------------------------------------------------------------------
func (r *products) BulkCreateProduct(data *[]mdm.ProductTemplate) error {
	err := r.db.Model(&mdm.ProductTemplate{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

// --------------------ProductToChannelMap-----------------------------------------------------------------------------------------------------------
func (r *products) CreateChannelProductMap(data *mdm.ProductToChannelMap) error {
	err := r.db.Model(&mdm.ProductToChannelMap{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) UpdateChannelProductMap(data *mdm.ProductToChannelMap, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductToChannelMap{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *products) DeleteChannelProductMap(query map[string]interface{}) error {
	var data mdm.ProductToChannelMap
	err := r.db.Model(&mdm.ProductToChannelMap{}).Where(query).Delete(&data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *products) FindChannelProductMap(query map[string]interface{}) (mdm.ProductToChannelMap, error) {
	var result mdm.ProductToChannelMap
	err := r.db.Model(&mdm.ProductToChannelMap{}).Where(query).First(&result).Error
	return result, err
}
func (r *products) GetAllChannelProductMap(p *pagination.Paginatevalue) ([]mdm.ProductToChannelMap, error) {
	var result []mdm.ProductToChannelMap
	err := r.db.Model(&mdm.ProductToChannelMap{}).Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.ProductToChannelMap{}, p, r.db)).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (r *products) GetAllChannelProductMapInArray(column string, condition string, array string) ([]mdm.ProductToChannelMap, error) {
	var result []mdm.ProductToChannelMap
	abc := fmt.Sprintf("%v %v (%v)", column, condition, array)
	err := r.db.Model(&mdm.ProductToChannelMap{}).Preload(clause.Associations).Where(abc).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (r *products) GetChannelProductMap(query interface{}) (mdm.ProductToChannelMap, error) {
	var result mdm.ProductToChannelMap
	err := r.db.Model(&mdm.ProductToChannelMap{}).Preload(clause.Associations).Where(query).First(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

// ----------------------Product Pricing Details-----------------------------------------------------------------------------------------------------
func (r *products) CreateProductPricingDetails(data *mdm.ProductPricingDetails) error {
	res := r.db.Model(&mdm.ProductPricingDetails{}).Create(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *products) UpdateProductPricingDetails(data *mdm.ProductPricingDetails, query map[string]interface{}) error {
	err := r.db.Model(&mdm.ProductPricingDetails{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

//------------hsn-------------

func (r *products) SaveHsn(data mdm.HSNCodesData) error {
	fmt.Println(data)
	err := r.db.Model(&mdm.HSNCodesData{}).Create(&data).Error
	if err != nil {
		fmt.Println("------------error---------------", err)
		return nil
	}
	return nil
}

func (r *products) FindOneHsn(query map[string]interface{}) (mdm.HSNCodesData, error) {

	var data mdm.HSNCodesData
	err := r.db.Model(&mdm.HSNCodesData{}).Preload(clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *products) FindAllHsn(page *pagination.Paginatevalue) ([]mdm.HSNCodesData, error) {

	var data []mdm.HSNCodesData
	err := r.db.Model(&mdm.HSNCodesData{}).Preload(clause.Associations).Scopes(helpers.Paginate(&mdm.HSNCodesData{}, page, r.db)).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *products) UpdateHsn(query map[string]interface{}, data mdm.HSNCodesData) error {

	err := r.db.Model(&mdm.HSNCodesData{}).Where(query).Updates(&data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *products) DeleteHsn(query map[string]interface{}) error {

	err := r.db.Model(&mdm.HSNCodesData{}).Where(query).Delete(&mdm.HSNCodesData{})
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
