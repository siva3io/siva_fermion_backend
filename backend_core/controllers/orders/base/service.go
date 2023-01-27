package orders_base

import (
	"errors"

	orders "fermion/backend_core/internal/repository/orders"
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
	//Internal Transfers Favourite
	FavouriteInternalTransfers(q map[string]interface{}) error
	//Internal Transfers UnFavourite
	UnFavouriteInternalTransfers(q map[string]interface{}) error

	//Sales Orders Favourite
	FavouriteSalesOrders(q map[string]interface{}) error
	//Sales Order UnFavourite
	UnFavouriteSalesOrders(q map[string]interface{}) error

	//Purchase Orders Favourite
	FavouritePurchaseOrders(q map[string]interface{}) error
	//Purchase Orders UnFavourite
	UnFavouritePurchaseOrders(q map[string]interface{}) error

	//Delivery Orders Favourite
	FavouriteDeliveryOrders(q map[string]interface{}) error
	//Delivery Orders UnFavourite
	UnFavouriteDeliveryOrders(q map[string]interface{}) error

	//Scrap Orders Favourite
	FavouriteScrapOrders(q map[string]interface{}) error
	//Scrap Orders UnFavourite
	UnFavouriteScrapOrders(q map[string]interface{}) error
}

type serviceBase struct {
	ordersRepository orders.OrdersBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	ordersRepository := orders.NewOrdersBase()
	newServiceObj = &serviceBase{ordersRepository}
	return newServiceObj
}

func (s *serviceBase) FavouriteInternalTransfers(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.ordersRepository.FindOne(query)
	samp := record.InternalTransferIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.InternalTransferIds = pq.Int64Array([]int64{int64(id)})
		err := s.ordersRepository.Create(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.InternalTransferIds = append(record.InternalTransferIds, int64(id))
		query := map[string]interface{}{
			"user_id":               uint(user_id),
			"internal_transfer_ids": record.InternalTransferIds,
		}
		err := s.ordersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteInternalTransfers(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.ordersRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.InternalTransferIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "internal_transfer_ids",
			"query":       "array_remove(internal_transfer_ids,?)",
		}
		err = s.ordersRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteSalesOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.ordersRepository.FindOne(query)
	samp := record.SalesOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.SalesOrderIds = pq.Int64Array([]int64{int64(id)})
		err := s.ordersRepository.Create(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.SalesOrderIds = append(record.SalesOrderIds, int64(id))
		query := map[string]interface{}{
			"user_id":         uint(user_id),
			"sales_order_ids": record.SalesOrderIds,
		}
		err := s.ordersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteSalesOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.ordersRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.SalesOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "sales_order_ids",
			"query":       "array_remove(sales_order_ids,?)",
		}
		err = s.ordersRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouritePurchaseOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.ordersRepository.FindOne(query)
	samp := record.PurchaseOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.PurchaseOrderIds = pq.Int64Array([]int64{int64(id)})
		err := s.ordersRepository.Create(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.PurchaseOrderIds = append(record.PurchaseOrderIds, int64(id))
		query := map[string]interface{}{
			"user_id":            uint(user_id),
			"purchase_order_ids": record.PurchaseOrderIds,
		}
		err := s.ordersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouritePurchaseOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.ordersRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.PurchaseOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "purchase_order_ids",
			"query":       "array_remove(purchase_order_ids,?)",
		}
		err = s.ordersRepository.UnFavourite(query)

		if err != nil {

			return err

		}

		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteDeliveryOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.ordersRepository.FindOne(query)
	samp := record.DeliveryOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.DeliveryOrderIds = pq.Int64Array([]int64{int64(id)})
		err := s.ordersRepository.Create(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.DeliveryOrderIds = append(record.DeliveryOrderIds, int64(id))
		query := map[string]interface{}{
			"user_id":            uint(user_id),
			"delivery_order_ids": record.DeliveryOrderIds,
		}
		err := s.ordersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteDeliveryOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.ordersRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.DeliveryOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "delivery_order_ids",
			"query":       "array_remove(delivery_order_ids,?)",
		}
		err = s.ordersRepository.UnFavourite(query)

		if err != nil {

			return err

		}

		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteScrapOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.ordersRepository.FindOne(query)
	samp := record.ScrapOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.ScrapOrderIds = pq.Int64Array([]int64{int64(id)})
		err := s.ordersRepository.Create(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.ScrapOrderIds = append(record.ScrapOrderIds, int64(id))
		query := map[string]interface{}{
			"user_id":         uint(user_id),
			"scrap_order_ids": record.ScrapOrderIds,
		}
		err := s.ordersRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteScrapOrders(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.ordersRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.ScrapOrderIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "scrap_order_ids",
			"query":       "array_remove(scrap_order_ids,?)",
		}
		err = s.ordersRepository.UnFavourite(query)

		if err != nil {

			return err

		}

		return nil
	}
	return errors.New("record not found in favourite list")
}
