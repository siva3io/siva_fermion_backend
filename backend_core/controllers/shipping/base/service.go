package shipping_base

import (
	"errors"

	shipping_repo "fermion/backend_core/internal/repository/shipping"
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
	//ShippingOrder
	FavouriteShippingOrder(q map[string]interface{}) error
	UnFavouriteShippingOrder(q map[string]interface{}) error

	//ShippingPartner
	FavouriteUserShippingPartnerRegistration(q map[string]interface{}) error
	UnFavouriteUserShippingPartnerRegistration(q map[string]interface{}) error
}

type serviceBase struct {
	shippingRepository shipping_repo.ShippingBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	shippingRepository := shipping_repo.NewShippingBase()
	newServiceObj = &serviceBase{shippingRepository}
	return newServiceObj
}

func (s *serviceBase) UnFavouriteShippingOrder(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.shippingRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.ShippingOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "shipping_order_ids",
			"query":       "array_remove(shipping_order_ids,?)",
		}
		err = s.shippingRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteShippingOrder(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.shippingRepository.FindOne(query)
	samp := record.ShippingOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.ShippingOrderIds = pq.Int64Array([]int64{int64(id)})
		err := s.shippingRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.ShippingOrderIds = append(record.ShippingOrderIds, int64(id))
		query := map[string]interface{}{
			"user_id":            uint(user_id),
			"shipping_order_ids": record.ShippingOrderIds,
		}
		err := s.shippingRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteUserShippingPartnerRegistration(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.shippingRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.ShippingPartnerIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "shipping_partner_ids",
			"query":       "array_remove(shipping_partner_ids,?)",
		}
		err = s.shippingRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteUserShippingPartnerRegistration(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.shippingRepository.FindOne(query)
	samp := record.ShippingPartnerIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.ShippingPartnerIds = pq.Int64Array([]int64{int64(id)})
		err := s.shippingRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.ShippingPartnerIds = append(record.ShippingPartnerIds, int64(id))
		query := map[string]interface{}{
			"user_id":              uint(user_id),
			"shipping_partner_ids": record.ShippingPartnerIds,
		}
		err := s.shippingRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}
