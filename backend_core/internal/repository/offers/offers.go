package repository

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/offers"
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
type Offer interface {
	ListOffers(page *pagination.Paginatevalue) (interface{}, error)
	ViewOffers(query map[string]interface{}) (interface{}, error)
	CreateOffers(data *offers.Offers) error
	UpdateOffers(query map[string]interface{}, data *offers.Offers) error
	DeleteOffers(query map[string]interface{}) error
	DeleteOfferLines(map[string]interface{}) error
	SaveOfferLines(*offers.OfferProductDetails) error
	UpdateOfferLines(map[string]interface{}, *offers.OfferProductDetails) (int64, error)
}

type Offers struct {
	db *gorm.DB
}

func NewOffers() *Offers {
	db := db.DbManager()
	return &Offers{db}
}

// ------------------------offers----------------------------------
func (r *Offers) ListOffers(page *pagination.Paginatevalue) (interface{}, error) {
	var data []offers.Offers
	err := r.db.Model(&offers.Offers{}).Scopes(helpers.Paginate(&offers.Offers{}, page, r.db)).Preload(clause.Associations).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}
func (r *Offers) ViewOffers(query map[string]interface{}) (interface{}, error) {
	var data offers.Offers
	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}
func (r *Offers) CreateOffers(data *offers.Offers) error {
	err := r.db.Model(&offers.Offers{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Offers) UpdateOffers(query map[string]interface{}, data *offers.Offers) error {
	res := r.db.Model(&offers.Offers{}).Where(query).Updates(data)
	if res.Error != nil {

		return res.Error

	}
	return nil
}
func (r *Offers) DeleteOffers(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&offers.Offers{}).Where(query).Updates(data)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *Offers) DeleteOfferLines(query map[string]interface{}) error {
	res := r.db.Model(&offers.OfferProductDetails{}).Where(query).Delete(&offers.OfferProductDetails{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}
func (r *Offers) SaveOfferLines(data *offers.OfferProductDetails) error {

	res := r.db.Model(&offers.OfferProductDetails{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}
func (r *Offers) UpdateOfferLines(query map[string]interface{}, data *offers.OfferProductDetails) (int64, error) {
	res := r.db.Model(&offers.OfferProductDetails{}).Where(query).Updates(&data)

	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}
