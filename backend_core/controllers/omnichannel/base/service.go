package marketplace_base

import (
	"errors"
	"fmt"

	omnichannel_repo "fermion/backend_core/internal/repository/omnichannel"
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
	//MarketPlace
	FavouriteMarketPlaceFav(map[string]interface{}) error
	UnFavouriteMarketPlaceFav(map[string]interface{}) error

	//Webstore
	FavouriteWebstoreFav(map[string]interface{}) error
	UnFavouriteWebstoreFav(map[string]interface{}) error

	//VirtualWarehouse
	FavouriteVirtualWarehouseFav(map[string]interface{}) error
	UnFavouriteVirtualWarehouseFav(map[string]interface{}) error
}

type serviceBase struct {
	omnichannelRepository omnichannel_repo.OmnichannelBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	omnichannelRepository := omnichannel_repo.NewmarketPlaceBase()
	newServiceObj = &serviceBase{omnichannelRepository}
	return newServiceObj

}

func (s *serviceBase) FavouriteMarketPlaceFav(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	fmt.Println("Service FAV -----------------------1")
	record, err := s.omnichannelRepository.FindOne(query)
	fmt.Println("findOne FAV -----------------------2", record)
	samp := record.MarketplaceIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.MarketplaceIds = pq.Int64Array([]int64{int64(id)})
		err := s.omnichannelRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.MarketplaceIds = append(record.MarketplaceIds, int64(id))
		query := map[string]interface{}{
			"user_id":         uint(user_id),
			"marketplace_ids": record.MarketplaceIds,
		}
		err := s.omnichannelRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteMarketPlaceFav(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.omnichannelRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.MarketplaceIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "marketplace_ids",
			"query":       "array_remove(marketplace_ids,?)",
		}
		err = s.omnichannelRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteWebstoreFav(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	fmt.Println("Service FAV -----------------------1")
	record, err := s.omnichannelRepository.FindOne(query)
	fmt.Println("findOne FAV -----------------------2", record)
	samp := record.WebstoreIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.WebstoreIds = pq.Int64Array([]int64{int64(id)})
		err := s.omnichannelRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.WebstoreIds = append(record.WebstoreIds, int64(id))
		query := map[string]interface{}{
			"user_id":      uint(user_id),
			"webstore_ids": record.WebstoreIds,
		}
		err := s.omnichannelRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteWebstoreFav(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.omnichannelRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.WebstoreIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "webstore_ids",
			"query":       "array_remove(webstore_ids,?)",
		}
		err = s.omnichannelRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

//---------VirtualWarehouse-----------

func (s *serviceBase) FavouriteVirtualWarehouseFav(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)

	query := map[string]interface{}{
		"user_id": user_id,
	}
	fmt.Println("Service FAV -----------------------1")
	record, err := s.omnichannelRepository.FindOne(query)
	fmt.Println("findOne FAV -----------------------2", record)
	samp := record.VirtualWarehouseIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if err != nil {
		record.UserID = uint(user_id)
		record.VirtualWarehouseIds = pq.Int64Array([]int64{int64(id)})
		err := s.omnichannelRepository.CreateRecord(record)
		if err != nil {
			return err
		}
	} else {
		record.VirtualWarehouseIds = append(record.VirtualWarehouseIds, int64(id))
		query := map[string]interface{}{
			"user_id":               uint(user_id),
			"virtual_warehouse_ids": record.VirtualWarehouseIds,
		}
		err := s.omnichannelRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteVirtualWarehouseFav(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.omnichannelRepository.FindOne(query)
	if err != nil {
		return err
	}
	samp := record.VirtualWarehouseIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "virtual_warehouse_ids",
			"query":       "array_remove(virtual_warehouse_ids,?)",
		}
		err = s.omnichannelRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}
