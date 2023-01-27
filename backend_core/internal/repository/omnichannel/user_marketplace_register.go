package omnichannel

import (
	"fmt"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/pkg/util/helpers"

	"fermion/backend_core/internal/model/pagination"

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
type Marketplace interface {
	CreateMarketplace(data *omnichannel.User_Marketplace_Registration) error
	Updatemarketplace(id uint, data *omnichannel.User_Marketplace_Registration) error
	FindAllmarketplace(p *pagination.Paginatevalue) ([]omnichannel.User_Marketplace_Registration, error)
	FindByIDMarketplace(id uint) (interface{}, error)
	DeleteMarketplace(id uint, user_id uint) error
	AvailableMarketPlaces(p *pagination.Paginatevalue) ([]omnichannel.Marketplace, error)
}

type marketplace struct {
	db *gorm.DB
}

var marketplaceRepository *marketplace //singleton object

// singleton function
func NewMarketplace() *marketplace {
	if marketplaceRepository != nil {
		return marketplaceRepository
	}
	db := db.DbManager()
	marketplaceRepository = &marketplace{db}
	return marketplaceRepository
}

func (r *marketplace) CreateMarketplace(data *omnichannel.User_Marketplace_Registration) error {

	res := r.db.Model(&omnichannel.User_Marketplace_Registration{}).Create(&data)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *marketplace) Updatemarketplace(id uint, data *omnichannel.User_Marketplace_Registration) error {

	res := r.db.Model(&omnichannel.User_Marketplace_Registration{}).Where("id", id).Updates(&data)
	return res.Error

}

func (r *marketplace) FindAllmarketplace(p *pagination.Paginatevalue) ([]omnichannel.User_Marketplace_Registration, error) {

	var result []omnichannel.User_Marketplace_Registration
	res := r.db.Preload(clause.Associations).Model(&omnichannel.User_Marketplace_Registration{}).Preload(clause.Associations).Scopes(helpers.Paginate(&omnichannel.User_Marketplace_Registration{}, p, r.db)).Where("is_active = true").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}

func (r *marketplace) FindByIDMarketplace(id uint) (interface{}, error) {
	var result omnichannel.User_Marketplace_Registration
	res := r.db.Preload(clause.Associations+"."+clause.Associations).Model(&omnichannel.User_Marketplace_Registration{}).Where("id", id).Preload(clause.Associations).Find(&result)
	if res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected == 0 {
		return nil, fmt.Errorf("no Data or record Found")
	}
	return result, nil
}

func (r *marketplace) DeleteMarketplace(id uint, user_id uint) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": user_id,
		"deleted_at": time.Now().In(loc),
	}
	res := r.db.Model(&omnichannel.User_Marketplace_Registration{}).Where("id", id).Updates(data)
	r.db.Where("id", id).Delete(&omnichannel.Bank_info{})
	r.db.Where("id", id).Delete(&omnichannel.Kyc_info{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *marketplace) AvailableMarketPlaces(p *pagination.Paginatevalue) ([]omnichannel.Marketplace, error) {

	var result []omnichannel.Marketplace
	res := r.db.Model(&omnichannel.Marketplace{}).Preload(clause.Associations).Scopes(helpers.Paginate(&omnichannel.Marketplace{}, p, r.db)).Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}
	return result, nil
}
