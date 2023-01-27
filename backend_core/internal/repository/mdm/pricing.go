package mdm

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
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
type Pricing interface {
	//====== pricing =====================================================
	CreatePricing(data *shared_pricing_and_location.Pricing) error
	UpdatePriceList(query map[string]interface{}, data *shared_pricing_and_location.Pricing) error
	DeletePricing(query map[string]interface{}) error
	FindPricing(query map[string]interface{}) (shared_pricing_and_location.Pricing, error)
	ListPricing(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error)

	//====== sales price list =============================================
	CreateSalePriceList(data *shared_pricing_and_location.SalesPriceList) error
	CreateSalesLineItems(data *shared_pricing_and_location.SalesLineItems) error
	UpdateSalesPriceList(query map[string]interface{}, data *shared_pricing_and_location.SalesPriceList) error
	UpdateSalesLineItems(query map[string]interface{}, data *shared_pricing_and_location.SalesLineItems) error
	DeleteSalesPriceList(query map[string]interface{}) error
	DeleteSalesLineItems(query map[string]interface{}) error
	FindOneSalesPriceLineItem(query map[string]interface{}) (shared_pricing_and_location.SalesLineItems, error)
	ListSalesPiceLineItems(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.SalesLineItems, error)

	//======= purchase price list =========================================
	CreatePurchasePriceList(data *shared_pricing_and_location.PurchasePriceList) error
	CreatePurchaseLineItems(data *shared_pricing_and_location.PurchaseLineItems) error
	UpdatePurchasePriceList(query map[string]interface{}, data *shared_pricing_and_location.PurchasePriceList) error
	UpdatePurchaseLineItems(query map[string]interface{}, data *shared_pricing_and_location.PurchaseLineItems) error
	DeletePurchasePriceList(query map[string]interface{}) error
	DeletePurchaseLineItems(query map[string]interface{}) error
	ListPurchasePricing(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchasePriceList, error)
	FindOnePurchaseLineItem(query map[string]interface{}) (shared_pricing_and_location.PurchaseLineItems, error)
	ListPurchasePriceListLineItems(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchaseLineItems, error)

	//=======  transfer price list ========================================
	CreateTransferPriceList(data *shared_pricing_and_location.TransferPriceList) error
	CreateTransferLineItems(data *shared_pricing_and_location.TransferLineItems) error
	UpdateTransferPriceList(query map[string]interface{}, data *shared_pricing_and_location.TransferPriceList) error
	UpdateTransferLineItems(query map[string]interface{}, data *shared_pricing_and_location.TransferLineItems) error
	DeleteTransferPriceList(query map[string]interface{}) error
	DeleteTransferLineItems(query map[string]interface{}) error
	FindOneTransferLineItem(query map[string]interface{}) (shared_pricing_and_location.TransferLineItems, error)
}

type pricing struct {
	db *gorm.DB
}

var pricingRepository *pricing //singleton object

// singleton function
func NewPricing() *pricing {
	if pricingRepository != nil {
		return pricingRepository
	}
	db := db.DbManager()
	pricingRepository = &pricing{db}
	return pricingRepository
}

// ========================================= pricings =================================================================
func (r *pricing) CreatePricing(data *shared_pricing_and_location.Pricing) error {
	err := r.db.Model(&shared_pricing_and_location.Pricing{}).Create(data).Error
	if err != nil {
		return err
	}
	return err
}
func (r *pricing) UpdatePriceList(query map[string]interface{}, data *shared_pricing_and_location.Pricing) error {
	err := r.db.Model(&shared_pricing_and_location.Pricing{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) DeletePricing(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.Pricing{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) FindPricing(query map[string]interface{}) (shared_pricing_and_location.Pricing, error) {
	var data shared_pricing_and_location.Pricing
	err := r.db.Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations).Model(&shared_pricing_and_location.Pricing{}).Preload(clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *pricing) ListPricing(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error) {
	var data []shared_pricing_and_location.Pricing
	err := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.Pricing{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// ========================================= sales price list =================================================================
func (r *pricing) CreateSalePriceList(data *shared_pricing_and_location.SalesPriceList) error {
	err := r.db.Model(&shared_pricing_and_location.SalesPriceList{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *pricing) CreateSalesLineItems(data *shared_pricing_and_location.SalesLineItems) error {
	err := r.db.Model(&shared_pricing_and_location.SalesLineItems{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *pricing) UpdateSalesPriceList(query map[string]interface{}, data *shared_pricing_and_location.SalesPriceList) error {
	err := r.db.Model(&shared_pricing_and_location.SalesPriceList{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) UpdateSalesLineItems(query map[string]interface{}, data *shared_pricing_and_location.SalesLineItems) error {
	err := r.db.Model(&shared_pricing_and_location.SalesLineItems{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) DeleteSalesPriceList(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.SalesPriceList{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) DeleteSalesLineItems(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.SalesLineItems{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) FindOneSalesPriceLineItem(query map[string]interface{}) (shared_pricing_and_location.SalesLineItems, error) {
	var data shared_pricing_and_location.SalesLineItems
	err := r.db.Model(&shared_pricing_and_location.SalesLineItems{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *pricing) ListSalesPiceLineItems(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.SalesLineItems, error) {
	var data []shared_pricing_and_location.SalesLineItems
	err := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.SalesLineItems{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// ========================================= purchase price list =================================================================
func (r *pricing) CreatePurchasePriceList(data *shared_pricing_and_location.PurchasePriceList) error {
	err := r.db.Model(&shared_pricing_and_location.PurchasePriceList{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *pricing) CreatePurchaseLineItems(data *shared_pricing_and_location.PurchaseLineItems) error {
	err := r.db.Model(&shared_pricing_and_location.PurchaseLineItems{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *pricing) UpdatePurchasePriceList(query map[string]interface{}, data *shared_pricing_and_location.PurchasePriceList) error {
	err := r.db.Model(&shared_pricing_and_location.PurchasePriceList{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) UpdatePurchaseLineItems(query map[string]interface{}, data *shared_pricing_and_location.PurchaseLineItems) error {
	err := r.db.Model(&shared_pricing_and_location.PurchaseLineItems{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) DeletePurchasePriceList(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.PurchasePriceList{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) DeletePurchaseLineItems(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.PurchaseLineItems{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) ListPurchasePricing(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchasePriceList, error) {
	var data []shared_pricing_and_location.PurchasePriceList
	err := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.PurchasePriceList{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *pricing) FindOnePurchaseLineItem(query map[string]interface{}) (shared_pricing_and_location.PurchaseLineItems, error) {
	var data shared_pricing_and_location.PurchaseLineItems
	err := r.db.Model(&shared_pricing_and_location.PurchaseLineItems{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
func (r *pricing) ListPurchasePriceListLineItems(query map[string]interface{}, p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchaseLineItems, error) {
	var data []shared_pricing_and_location.PurchaseLineItems
	err := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.PurchaseLineItems{}, p, r.db)).Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

// ========================================= transfer price list =================================================================
func (r *pricing) CreateTransferPriceList(data *shared_pricing_and_location.TransferPriceList) error {
	err := r.db.Model(&shared_pricing_and_location.TransferPriceList{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *pricing) CreateTransferLineItems(data *shared_pricing_and_location.TransferLineItems) error {
	err := r.db.Model(&shared_pricing_and_location.TransferLineItems{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *pricing) UpdateTransferPriceList(query map[string]interface{}, data *shared_pricing_and_location.TransferPriceList) error {
	err := r.db.Model(&shared_pricing_and_location.TransferPriceList{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) UpdateTransferLineItems(query map[string]interface{}, data *shared_pricing_and_location.TransferLineItems) error {
	err := r.db.Model(&shared_pricing_and_location.TransferLineItems{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) DeleteTransferPriceList(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.TransferPriceList{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) DeleteTransferLineItems(query map[string]interface{}) error {
	timeZone := os.Getenv("DB_TZ")
	timeLocation, _ := time.LoadLocation(timeZone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(timeLocation),
	}
	delete(query, "user_id")
	err := r.db.Model(&shared_pricing_and_location.TransferLineItems{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("oops! record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
func (r *pricing) FindOneTransferLineItem(query map[string]interface{}) (shared_pricing_and_location.TransferLineItems, error) {
	var data shared_pricing_and_location.TransferLineItems
	err := r.db.Model(&shared_pricing_and_location.TransferLineItems{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("oops! record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}
