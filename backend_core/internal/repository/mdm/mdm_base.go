package mdm

import (
	"errors"
	"fmt"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

	"github.com/lib/pq"
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
type MdmBase interface {
	CreateRecord(data mdm.UserMdmFav) error
	Favourite(query map[string]interface{}) error
	FindOne(query map[string]interface{}) (mdm.UserMdmFav, error)
	UnFavourite(query map[string]interface{}) error

	//contact
	GetFavContactList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Partner, error)

	//location
	GetFavLocationList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Locations, error)

	//vendor
	GetFavVendorList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Vendors, error)

	//Product
	GetFavProductList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.ProductVariant, error)

	//Pricing
	GetFavPricingList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error)
}

type mdmBase struct {
	db *gorm.DB
}

var mdmBaseRepository *mdmBase //singleton object

// singleton function
func NewMdmBase() *mdmBase {
	if mdmBaseRepository != nil {
		return mdmBaseRepository
	}
	db := db.DbManager()
	mdmBaseRepository = &mdmBase{db}
	return mdmBaseRepository

}

func (r *mdmBase) CreateRecord(data mdm.UserMdmFav) error {
	fmt.Println(data)
	err := r.db.Model(&mdm.UserMdmFav{}).Create(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *mdmBase) Favourite(query map[string]interface{}) error {
	//fmt.Println(data)
	err := r.db.Model(&mdm.UserMdmFav{}).Where("user_id = ?", query["user_id"]).Updates(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *mdmBase) UnFavourite(query map[string]interface{}) error {
	field := query["field_value"].(string)
	q := query["query"].(string)
	err := r.db.Model(&mdm.UserMdmFav{}).Where("user_id = ?", query["user_id"]).Update(field, gorm.Expr(q, query["id"])).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *mdmBase) FindOne(query map[string]interface{}) (mdm.UserMdmFav, error) {
	var data mdm.UserMdmFav

	err := r.db.Preload(clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}

	return data, nil
}

//UPDATE "user_mdm_favs" SET "contact_ids"=array_remove(contact_ids,8),"updated_date"='2022-06-03 15:12:27.387' WHERE user_id = 1 AND "user_mdm_favs"."deleted_at" IS NULL

func (r *mdmBase) GetFavContactList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Partner, error) {
	var result []mdm.Partner
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&mdm.Partner{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&mdm.Partner{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

//SELECT * FROM public.partners WHERE Id IN (2,3)

func (r *mdmBase) GetFavLocationList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Locations, error) {
	var result []shared_pricing_and_location.Locations
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&shared_pricing_and_location.Locations{}).Scopes(helpers.Paginate(&shared_pricing_and_location.Locations{}, p, r.db)).Where("Id IN ?", query).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mdmBase) GetFavVendorList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Vendors, error) {
	var result []mdm.Vendors
	query := []int64(data)
	err := r.db.Preload(clause.Associations).Model(&mdm.Vendors{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&mdm.Vendors{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mdmBase) GetFavProductList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.ProductVariant, error) {
	var result []mdm.ProductVariant
	query := []int64(data)
	err := r.db.Model(&mdm.ProductVariant{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&mdm.ProductVariant{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *mdmBase) GetFavPricingList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error) {
	var result []shared_pricing_and_location.Pricing
	query := []int64(data)
	err := r.db.Model(&shared_pricing_and_location.Pricing{}).Where("Id IN ?", query).Scopes(helpers.Paginate(&shared_pricing_and_location.Pricing{}, p, r.db)).Find(&result).Count(&p.TotalRows).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
