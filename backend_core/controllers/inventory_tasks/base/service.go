package inventory_tasks_base

import (
	"errors"

	"fermion/backend_core/internal/model/inventory_tasks"
	"fermion/backend_core/internal/model/pagination"
	inventory_tasks_repo "fermion/backend_core/internal/repository/inventory_tasks"
	"fermion/backend_core/pkg/util/helpers"

	"github.com/lib/pq"
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
type ServiceBase interface {
	//Cycle Count
	FavouriteCycleCount(map[string]interface{}) error
	UnFavouriteCycleCount(map[string]interface{}) error
	GetArrCycleCount(user_id int) (pq.Int64Array, error)
	GetFavCycleCountList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.CycleCount, error)

	//Pick List
	FavouritePickList(map[string]interface{}) error
	UnFavouritePickList(map[string]interface{}) error
	GetArrPicList(user_id int) (pq.Int64Array, error)
	GetFavPicList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error)
}

type serviceBase struct {
	inventoryTasksRepository inventory_tasks_repo.InventoryTasksBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	inventoryTasksRepository := inventory_tasks_repo.NewInventoryTasksBase()
	newServiceObj = &serviceBase{inventoryTasksRepository}
	return newServiceObj
}
func (s *serviceBase) GetArrCycleCount(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.inventoryTasksRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("Cycle count favourites data Not found")
	}
	return record.CycleCountIds, er
}
func (s *serviceBase) GetFavCycleCountList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.CycleCount, error) {
	result, err := s.inventoryTasksRepository.GetCycleCountFavourites(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *serviceBase) GetArrPicList(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.inventoryTasksRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("Pick list favourites data Not found")
	}
	return record.PickListIds, er
}

func (s *serviceBase) GetFavPicList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_tasks.PickList, error) {
	result, err := s.inventoryTasksRepository.GetPickListFavourites(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *serviceBase) FavouriteCycleCount(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryTasksRepository.FindOne(query)
	samp := record.CycleCountIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.CycleCountIds = pq.Int64Array([]int64{int64(id)})
		err := s.inventoryTasksRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.CycleCountIds = append(record.CycleCountIds, int64(id))
		query := map[string]interface{}{
			"user_id":         uint(user_id),
			"cycle_count_ids": record.CycleCountIds,
		}
		err := s.inventoryTasksRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteCycleCount(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryTasksRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.CycleCountIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "cycle_count_ids",
			"query":       "array_remove(cycle_count_ids,?)",
		}
		err = s.inventoryTasksRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouritePickList(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryTasksRepository.FindOne(query)
	samp := record.PickListIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.PickListIds = pq.Int64Array([]int64{int64(id)})
		err := s.inventoryTasksRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.PickListIds = append(record.PickListIds, int64(id))
		query := map[string]interface{}{
			"user_id":       uint(user_id),
			"pick_list_ids": record.PickListIds,
		}
		err := s.inventoryTasksRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouritePickList(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryTasksRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.PickListIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "pick_list_ids",
			"query":       "array_remove(pick_list_ids,?)",
		}
		err = s.inventoryTasksRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}
