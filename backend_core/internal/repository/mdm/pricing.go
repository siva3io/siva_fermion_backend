package mdm

import (
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
	CreatePricing(db *gorm.DB, data *shared_pricing_and_location.Pricing) error
	CreatePurchasePriceList(db *gorm.DB, data *shared_pricing_and_location.PurchasePriceList) error
	CreateTransferPriceList(db *gorm.DB, data *shared_pricing_and_location.TransferPriceList) error
	CreateSalePriceList(db *gorm.DB, data *shared_pricing_and_location.SalesPriceList) error
	FindPricing(query map[string]interface{}) (shared_pricing_and_location.Pricing, error)
	ListPricing(p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error)
	UpdatePriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.Pricing) error
	UpdateSalesPriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.SalesPriceList) error
	CreateSalesLineItems(db *gorm.DB, data *shared_pricing_and_location.SalesLineItems) error
	UpdateSalesLineItems(db *gorm.DB, query interface{}, data *shared_pricing_and_location.SalesLineItems) (int64, error)
	UpdatePurchasePriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.PurchasePriceList) error
	CreatePurchaseLineItems(db *gorm.DB, data *shared_pricing_and_location.PurchaseLineItems) error
	UpdatePurchaseLineItems(db *gorm.DB, query interface{}, data *shared_pricing_and_location.PurchaseLineItems) (int64, error)
	UpdateTransferPriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.TransferPriceList) error
	CreateTransferLineItems(db *gorm.DB, data *shared_pricing_and_location.TransferLineItems) error
	UpdateTransferLineItems(db *gorm.DB, query interface{}, data *shared_pricing_and_location.TransferLineItems) (int64, error)
	DeletePricing(db *gorm.DB, id uint, user_id uint) error
	DeleteSalesLineItems(db *gorm.DB, query interface{}) error
	DeletePurchaseLineItems(db *gorm.DB, query interface{}) error
	DeleteTransferLineItems(db *gorm.DB, query interface{}) error
	DeleteSalesPriceList(db *gorm.DB, id uint) error
	DeletePurchasePriceList(db *gorm.DB, id uint) error
	DeleteTransferPriceList(db *gorm.DB, id uint) error
	ListPurchasePricing(p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchasePriceList, error)

	FindOneSalesPriceLineItem(query map[string]interface{}) (shared_pricing_and_location.SalesLineItems, error)
	ListSalesPiceLineItems(p *pagination.Paginatevalue) ([]shared_pricing_and_location.SalesLineItems, error)
	ListPurchasePriceListLineItems(p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchaseLineItems, error)
}

type pricing struct {
	db *gorm.DB
}

func NewPricing() *pricing {
	db := db.DbManager()
	return &pricing{db}
}

// Create a Pricing
func (r *pricing) CreatePricing(db *gorm.DB, data *shared_pricing_and_location.Pricing) error {
	res := db.Model(&shared_pricing_and_location.Pricing{}).Create(data)
	return res.Error
}

// Delete Pricing
func (r *pricing) DeletePricing(db *gorm.DB, id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := db.Model(&shared_pricing_and_location.Pricing{}).Where("id", id).Updates(data)
	return res.Error
}

// Create PurchasePriceList
func (r *pricing) CreatePurchasePriceList(db *gorm.DB, data *shared_pricing_and_location.PurchasePriceList) error {
	res := db.Model(&shared_pricing_and_location.PurchasePriceList{}).Create(data)
	return res.Error
}

// Create TransferPriceList
func (r *pricing) CreateTransferPriceList(db *gorm.DB, data *shared_pricing_and_location.TransferPriceList) error {
	res := db.Model(&shared_pricing_and_location.TransferPriceList{}).Create(data)
	return res.Error
}

// Create Sales Price List
func (r *pricing) CreateSalePriceList(db *gorm.DB, data *shared_pricing_and_location.SalesPriceList) error {
	res := db.Model(&shared_pricing_and_location.SalesPriceList{}).Create(data)
	return res.Error
}

// Find Pricing
func (r *pricing) FindPricing(query map[string]interface{}) (shared_pricing_and_location.Pricing, error) {
	var result shared_pricing_and_location.Pricing
	res := r.db.Preload(clause.Associations + "." + clause.Associations + "." + clause.Associations + "." + clause.Associations).Model(&shared_pricing_and_location.Pricing{}).Preload(clause.Associations).Where(query).First(&result)

	if result.SalesPriceListId != nil {
		r.db.Preload("SalesLineItems.Product").Preload("SalesLineItems.UOM").Preload("SalesLineItems.Quantity_value_type").Preload(clause.Associations).Find(&result.SalesPriceList)
	} else if result.PurchasePriceListId != nil {
		r.db.Preload("PurchaseLineItems.Product").Preload("PurchasePriceList.Currency").Preload("PurchaseLineItems.Quantity_value_type").Preload(clause.Associations).Find(&result.PurchasePriceList)
	} else if result.TransferPriceListId != nil {
		r.db.Preload("TransferLineItems.Product").Preload(clause.Associations + "." + clause.Associations).Find(&result.TransferPriceList)
	}
	return result, res.Error
}

// List Pricing
func (r *pricing) ListPricing(p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error) {
	var result []shared_pricing_and_location.Pricing
	res := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.Pricing{}, p, r.db)).Where("is_active = true").Find(&result)
	return result, res.Error
}

// Update Pricing
func (r *pricing) UpdatePriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.Pricing) error {
	res := db.Model(&shared_pricing_and_location.Pricing{}).Where("id", id).Updates(data)
	return res.Error
}

// Update Sales Price List
func (r *pricing) UpdateSalesPriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.SalesPriceList) error {
	res := db.Model(&shared_pricing_and_location.SalesPriceList{}).Where("id", id).Updates(data)
	return res.Error
}

// Create sales Line Items
func (r *pricing) CreateSalesLineItems(db *gorm.DB, data *shared_pricing_and_location.SalesLineItems) error {
	res := db.Model(&shared_pricing_and_location.SalesLineItems{}).Create(data)
	return res.Error
}

// Update sales Line Items
func (r *pricing) UpdateSalesLineItems(db *gorm.DB, query interface{}, data *shared_pricing_and_location.SalesLineItems) (int64, error) {
	res := db.Model(&shared_pricing_and_location.SalesLineItems{}).Where(query).Updates(data)
	return res.RowsAffected, res.Error
}

// Update  PurchasePriceList
func (r *pricing) UpdatePurchasePriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.PurchasePriceList) error {
	res := db.Model(&shared_pricing_and_location.PurchasePriceList{}).Where("id", id).Updates(data)
	return res.Error
}

// Create Purchase Line Items
func (r *pricing) CreatePurchaseLineItems(db *gorm.DB, data *shared_pricing_and_location.PurchaseLineItems) error {
	res := db.Model(&shared_pricing_and_location.PurchaseLineItems{}).Create(data)
	return res.Error
}

// Update Purchase Line Items
func (r *pricing) UpdatePurchaseLineItems(db *gorm.DB, query interface{}, data *shared_pricing_and_location.PurchaseLineItems) (int64, error) {
	res := db.Model(&shared_pricing_and_location.PurchaseLineItems{}).Where(query).Updates(data)
	return res.RowsAffected, res.Error
}

// Update TransferPriceList
func (r *pricing) UpdateTransferPriceList(db *gorm.DB, id uint, data *shared_pricing_and_location.TransferPriceList) error {
	res := db.Model(&shared_pricing_and_location.TransferPriceList{}).Where("id", id).Updates(data)
	return res.Error
}

// Create Transfer Line Items
func (r *pricing) CreateTransferLineItems(db *gorm.DB, data *shared_pricing_and_location.TransferLineItems) error {
	res := db.Model(&shared_pricing_and_location.TransferLineItems{}).Create(data)
	return res.Error
}

// Update transfer Line Items
func (r *pricing) UpdateTransferLineItems(db *gorm.DB, query interface{}, data *shared_pricing_and_location.TransferLineItems) (int64, error) {
	res := db.Model(&shared_pricing_and_location.TransferLineItems{}).Where(query).Updates(data)
	return res.RowsAffected, res.Error
}

// Delete Sales Price List
func (r *pricing) DeleteSalesPriceList(db *gorm.DB, id uint) error {
	res := db.Model(&shared_pricing_and_location.SalesPriceList{}).Where("id", id).Delete(&shared_pricing_and_location.SalesPriceList{})
	return res.Error
}

// Delete Purchase Price List
func (r *pricing) DeletePurchasePriceList(db *gorm.DB, id uint) error {
	res := db.Model(&shared_pricing_and_location.PurchasePriceList{}).Where("id", id).Delete(&shared_pricing_and_location.PurchasePriceList{})
	return res.Error
}

// Delete Transfer Price List
func (r *pricing) DeleteTransferPriceList(db *gorm.DB, id uint) error {
	res := db.Model(&shared_pricing_and_location.TransferPriceList{}).Where("id", id).Delete(&shared_pricing_and_location.TransferPriceList{})
	return res.Error
}

// Delete Sales Line Items
func (r *pricing) DeleteSalesLineItems(db *gorm.DB, query interface{}) error {
	res := db.Model(&shared_pricing_and_location.SalesLineItems{}).Where(query).Delete(&shared_pricing_and_location.SalesLineItems{})
	return res.Error
}

// Delete Purchase Line Items
func (r *pricing) DeletePurchaseLineItems(db *gorm.DB, query interface{}) error {
	res := db.Model(&shared_pricing_and_location.PurchaseLineItems{}).Where(query).Delete(&shared_pricing_and_location.PurchaseLineItems{})
	return res.Error
}

// Delete Transfer Line Items
func (r *pricing) DeleteTransferLineItems(db *gorm.DB, query interface{}) error {
	res := db.Model(&shared_pricing_and_location.TransferLineItems{}).Where(query).Delete(&shared_pricing_and_location.TransferLineItems{})
	return res.Error
}

func (r *pricing) ListPurchasePricing(p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchasePriceList, error) {
	var result []shared_pricing_and_location.PurchasePriceList
	res := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.PurchasePriceList{}, p, r.db)).Where("is_active = true").Find(&result)
	return result, res.Error
}

func (r *pricing) FindOneSalesPriceLineItem(query map[string]interface{}) (shared_pricing_and_location.SalesLineItems, error) {
	var result shared_pricing_and_location.SalesLineItems
	res := r.db.Model(&shared_pricing_and_location.SalesLineItems{}).Where(query).First(&result)

	return result, res.Error
}
func (r *pricing) ListSalesPiceLineItems(p *pagination.Paginatevalue) ([]shared_pricing_and_location.SalesLineItems, error) {
	var result []shared_pricing_and_location.SalesLineItems
	res := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.SalesLineItems{}, p, r.db)).Where("is_active = true").Find(&result)
	return result, res.Error
}
func (r *pricing) ListPurchasePriceListLineItems(p *pagination.Paginatevalue) ([]shared_pricing_and_location.PurchaseLineItems, error) {
	var result []shared_pricing_and_location.PurchaseLineItems
	res := r.db.Preload(clause.Associations).Scopes(helpers.Paginate(&shared_pricing_and_location.PurchaseLineItems{}, p, r.db)).Where("is_active = true").Find(&result)
	return result, res.Error
}
