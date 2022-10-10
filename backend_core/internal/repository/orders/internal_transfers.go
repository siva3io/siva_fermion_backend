package orders

import (
	"errors"
	"fmt"
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
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type InternalTransfers interface {
	Save(data *orders.InternalTransfers) error
	FindAll(page *pagination.Paginatevalue) (interface{}, error)
	FindOne(query map[string]interface{}) (interface{}, error)
	Update(query map[string]interface{}, data *orders.InternalTransfers) error
	Delete(query map[string]interface{}) error

	SaveOrderLines(orders.InternalTransferLines) error
	UpdateOrderLines(map[string]interface{}, orders.InternalTransferLines) (int64, error)
	DeleteOrderLine(map[string]interface{}) error
	FindOrderLines(map[string]interface{}) (orders.InternalTransferLines, error)

	Search(query string) (interface{}, error)
}

type internal_transfers struct {
	db *gorm.DB
}

func NewInternalTransfer() *internal_transfers {
	db := db.DbManager()
	return &internal_transfers{db}

}

func (r *internal_transfers) Save(data *orders.InternalTransfers) error {
	err := r.db.Model(&orders.InternalTransfers{}).Create(data).Error

	if err != nil {

		return err

	}

	return nil
}

func (r *internal_transfers) FindAll(page *pagination.Paginatevalue) (interface{}, error) {
	var data []orders.InternalTransfers

	err := r.db.Model(&orders.InternalTransfers{}).Scopes(helpers.Paginate(&orders.InternalTransfers{}, page, r.db)).Preload(clause.Associations).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *internal_transfers) FindOne(query map[string]interface{}) (interface{}, error) {
	var data orders.InternalTransfers

	err := r.db.Preload(clause.Associations + "." + clause.Associations).Where(query).First(&data)

	if err.RowsAffected == 0 {
		return nil, errors.New("record not found")
	}

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}

func (r *internal_transfers) Update(query map[string]interface{}, data *orders.InternalTransfers) error {
	res := r.db.Model(&orders.InternalTransfers{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *internal_transfers) Delete(query map[string]interface{}) error {
	zone := os.Getenv("DB_TZ")
	loc, _ := time.LoadLocation(zone)
	data := map[string]interface{}{
		"deleted_by": query["user_id"].(int),
		"deleted_at": time.Now().In(loc),
	}
	delete(query, "user_id")
	res := r.db.Model(&orders.InternalTransfers{}).Where(query).Updates(data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *internal_transfers) SaveOrderLines(data orders.InternalTransferLines) error {

	res := r.db.Model(&orders.InternalTransferLines{}).Create(&data)

	if res.Error != nil {

		return res.Error

	}

	return nil
}

func (r *internal_transfers) FindOrderLines(query map[string]interface{}) (orders.InternalTransferLines, error) {
	var result orders.InternalTransferLines
	fmt.Println(query)
	res := r.db.Model(&orders.InternalTransferLines{}).Where(query).First(&result)

	if res.Error != nil {
		return result, res.Error
	}

	return result, nil
}

func (r *internal_transfers) UpdateOrderLines(query map[string]interface{}, data orders.InternalTransferLines) (int64, error) {
	res := r.db.Model(&orders.InternalTransferLines{}).Where(query).Updates(&data)

	if res.Error != nil {

		return res.RowsAffected, res.Error

	}

	return res.RowsAffected, nil
}

func (r *internal_transfers) DeleteOrderLine(query map[string]interface{}) error {
	res := r.db.Model(&orders.InternalTransferLines{}).Where(query).Delete(&orders.InternalTransferLines{})

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *internal_transfers) Search(query string) (interface{}, error) {
	var data []orders.InternalTransfers

	fields := []string{"reference_number", "ist_number"}

	fields_string, values := helpers.ApplySearch(query, fields)

	err := r.db.Model(&orders.InternalTransfers{}).Limit(2).Preload(clause.Associations).Where(fields_string, values...).Find(&data)

	if err.Error != nil {
		return nil, err.Error
	}

	return data, nil
}
