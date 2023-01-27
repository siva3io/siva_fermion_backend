package inventory_orders_base

import (
	"errors"

	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/pagination"
	inventory_orders_repo "fermion/backend_core/internal/repository/inventory_orders"
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
	//Asn
	FavouriteAsn(map[string]interface{}) error
	UnFavouriteAsn(map[string]interface{}) error
	GetArrAsn(user_id int) (pq.Int64Array, error)
	GetFavAsnList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.ASN, error)

	//Inventory Adjustment
	FavouriteInvAdj(map[string]interface{}) error
	UnFavouriteInvAdj(map[string]interface{}) error
	GetArrInvAdj(user_id int) (pq.Int64Array, error)
	GetFavInvAdjList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error)

	//Grn
	FavouriteGrn(map[string]interface{}) error
	UnFavouriteGrn(map[string]interface{}) error
	GetArrGrn(user_id int) (pq.Int64Array, error)
	GetGrnFavList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.GRN, error)
}

type serviceBase struct {
	inventoryOrdersRepository inventory_orders_repo.InventoryOrdersBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	inventoryOrdersRepository := inventory_orders_repo.NewInventoryOrdersBase()
	newServiceObj = &serviceBase{inventoryOrdersRepository}
	return newServiceObj
}
func (s *serviceBase) GetArrInvAdj(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.inventoryOrdersRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("Inventory adjustment favourites data Not found")
	}

	return record.InvAdjIds, er
}
func (s *serviceBase) GetFavInvAdjList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.InventoryAdjustments, error) {
	result, err := s.inventoryOrdersRepository.GetInvAdjFavourites(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *serviceBase) GetArrGrn(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.inventoryOrdersRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("GRN favourites data Not found")
	}

	return record.GrnIds, er
}
func (s *serviceBase) GetGrnFavList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.GRN, error) {
	result, err := s.inventoryOrdersRepository.GetGrnFavourites(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *serviceBase) GetArrAsn(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.inventoryOrdersRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("ASN favourites data Not found")
	}

	return record.AsnIds, er
}
func (s *serviceBase) GetFavAsnList(data pq.Int64Array, p *pagination.Paginatevalue) ([]inventory_orders.ASN, error) {
	result, err := s.inventoryOrdersRepository.GetAsnFavourites(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *serviceBase) FavouriteAsn(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryOrdersRepository.FindOne(query)
	samp := record.AsnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.AsnIds = pq.Int64Array([]int64{int64(id)})
		err := s.inventoryOrdersRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.AsnIds = append(record.AsnIds, int64(id))
		query := map[string]interface{}{
			"user_id": uint(user_id),
			"asn_ids": record.AsnIds,
		}
		err := s.inventoryOrdersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteAsn(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryOrdersRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.AsnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "asn_ids",
			"query":       "array_remove(asn_ids,?)",
		}
		err = s.inventoryOrdersRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteInvAdj(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryOrdersRepository.FindOne(query)
	samp := record.InvAdjIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.InvAdjIds = pq.Int64Array([]int64{int64(id)})
		err := s.inventoryOrdersRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.InvAdjIds = append(record.InvAdjIds, int64(id))
		query := map[string]interface{}{
			"user_id":     uint(user_id),
			"inv_adj_ids": record.InvAdjIds,
		}
		err := s.inventoryOrdersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteInvAdj(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryOrdersRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.InvAdjIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "inv_adj_ids",
			"query":       "array_remove(inv_adj_ids,?)",
		}
		err = s.inventoryOrdersRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteGrn(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryOrdersRepository.FindOne(query)
	samp := record.GrnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.GrnIds = pq.Int64Array([]int64{int64(id)})
		err := s.inventoryOrdersRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.GrnIds = append(record.GrnIds, int64(id))
		query := map[string]interface{}{
			"user_id": uint(user_id),
			"grn_ids": record.GrnIds,
		}
		err := s.inventoryOrdersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteGrn(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.inventoryOrdersRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.GrnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "grn_ids",
			"query":       "array_remove(grn_ids,?)",
		}
		err = s.inventoryOrdersRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}
