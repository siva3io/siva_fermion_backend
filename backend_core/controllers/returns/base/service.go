package mdm_base

import (
	"errors"
	"fmt"

	returns_repo "fermion/backend_core/internal/repository/returns"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

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
	//SalesReturns
	FavouriteSalesReturns(q map[string]interface{}) error
	UnFavouriteSalesReturns(q map[string]interface{}) error

	//PurchaseReturns
	FavouritePurchaseReturns(q map[string]interface{}) error
	UnFavouritePurchaseReturns(q map[string]interface{}) error
}

type serviceBase struct {
	returnsRepository returns_repo.ReturnsBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	returnsRepository := returns_repo.NewReturnsBase()
	newServiceObj = &serviceBase{returnsRepository}
	return newServiceObj
}

func (s *serviceBase) FavouriteSalesReturns(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.returnsRepository.FindOne(query)
	samp := record.SalesReturnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.SalesReturnIds = pq.Int64Array([]int64{int64(id)})
		err := s.returnsRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.SalesReturnIds = append(record.SalesReturnIds, int64(id))
		query := map[string]interface{}{
			"user_id":          uint(user_id),
			"sales_return_ids": record.SalesReturnIds,
		}
		err := s.returnsRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteSalesReturns(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.returnsRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.SalesReturnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "sales_return_ids",
			"query":       "array_remove(sales_return_ids,?)",
		}
		err = s.returnsRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouritePurchaseReturns(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	fmt.Println(id, user_id)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.returnsRepository.FindOne(query)
	samp := record.PurchaseReturnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.PurchaseReturnIds = pq.Int64Array([]int64{int64(id)})
		err := s.returnsRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.PurchaseReturnIds = append(record.PurchaseReturnIds, int64(id))
		query := map[string]interface{}{
			"user_id":             uint(user_id),
			"purchase_return_ids": record.PurchaseReturnIds,
		}
		err := s.returnsRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouritePurchaseReturns(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.returnsRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.PurchaseReturnIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "purchase_return_ids",
			"query":       "array_remove(purchase_return_ids,?)",
		}
		err = s.returnsRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}
