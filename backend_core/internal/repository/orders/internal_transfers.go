package orders

import (
	"errors"
	"os"
	"time"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/orders"
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
type InternalTransfers interface {
	Save(data *orders.InternalTransfers) error
	SaveOrderLines(data *orders.InternalTransferLines) error

	Update(query map[string]interface{}, data *orders.InternalTransfers) error
	UpdateOrderLines(query map[string]interface{}, data *orders.InternalTransferLines) error

	Delete(query map[string]interface{}) error
	DeleteOrderLine(query map[string]interface{}) error

	FindOne(query map[string]interface{}) (orders.InternalTransfers, error)

	FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.InternalTransfers, error)

	FindOrderLines(query map[string]interface{}) (orders.InternalTransferLines, error)
}

type internal_transfers struct {
	db *gorm.DB
}

var internalTransfersRepository *internal_transfers //singleton object

// singleton function
func NewInternalTransfer() *internal_transfers {
	if internalTransfersRepository != nil {
		return internalTransfersRepository
	}
	internalTransfersRepository = &internal_transfers{
		db.DbManager(),
	}
	return internalTransfersRepository

}

func (r *internal_transfers) Save(data *orders.InternalTransfers) error {

	err := r.db.Model(&orders.InternalTransfers{}).Create(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *internal_transfers) FindAll(query map[string]interface{}, page *pagination.Paginatevalue) ([]orders.InternalTransfers, error) {

	var data []orders.InternalTransfers
	err := r.db.Model(&orders.InternalTransfers{}).Preload(clause.Associations).Scopes(helpers.Paginate(&orders.InternalTransfers{}, page, r.db)).Where(query).Find(&data)
	if err.Error != nil {
		return nil, err.Error
	}
	return data, nil
}

func (r *internal_transfers) FindOne(query map[string]interface{}) (orders.InternalTransfers, error) {

	var data orders.InternalTransfers
	err := r.db.Model(&orders.InternalTransfers{}).Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *internal_transfers) Update(query map[string]interface{}, data *orders.InternalTransfers) error {

	err := r.db.Model(&orders.InternalTransfers{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *internal_transfers) Delete(query map[string]interface{}) error {

	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"],
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	err := r.db.Model(&orders.InternalTransfers{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *internal_transfers) SaveOrderLines(data *orders.InternalTransferLines) error {

	err := r.db.Model(&orders.InternalTransferLines{}).Create(data)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *internal_transfers) FindOrderLines(query map[string]interface{}) (orders.InternalTransferLines, error) {
	var data orders.InternalTransferLines
	err := r.db.Model(&orders.InternalTransferLines{}).Where(query).First(&data)
	if err.RowsAffected == 0 {
		return data, errors.New("record not found")
	}
	if err.Error != nil {
		return data, err.Error
	}
	return data, nil
}

func (r *internal_transfers) UpdateOrderLines(query map[string]interface{}, data *orders.InternalTransferLines) error {

	err := r.db.Model(&orders.InternalTransferLines{}).Where(query).Updates(data)
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *internal_transfers) DeleteOrderLine(query map[string]interface{}) error {

	err := r.db.Model(&orders.InternalTransferLines{}).Where(query).Delete(&orders.InternalTransferLines{})
	if err.RowsAffected == 0 {
		return errors.New("record not found")
	}
	if err.Error != nil {
		return err.Error
	}
	return nil
}
