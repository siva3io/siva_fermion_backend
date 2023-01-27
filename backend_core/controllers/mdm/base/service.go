package mdm_base

import (
	"errors"
	"strconv"

	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
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
	//Contacts
	FavouriteContacts(q map[string]interface{}) error
	UnFavouriteContacts(q map[string]interface{}) error
	GetFavContactList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Partner, error)
	GetArrContact(user_id int) (pq.Int64Array, error)

	//Locations
	FavouriteLocations(q map[string]interface{}) error
	UnFavouriteLocations(q map[string]interface{}) error
	GetFavLocationList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Locations, error)
	GetArrLocation(user_id int) (pq.Int64Array, error)

	//Vendors
	FavouriteVendors(q map[string]interface{}) error
	UnfavouriteVendors(q map[string]interface{}) error
	GetFavVendorList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Vendors, error)
	GetArrVendor(user_id int) (pq.Int64Array, error)

	//Products
	FavouriteProducts(q map[string]interface{}) error
	UnFavouriteProducts(q map[string]interface{}) error
	GetFavProductList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.ProductVariant, error)
	GetArrProduct(user_id int) (pq.Int64Array, error)

	//Pricing
	FavouritePricings(q map[string]interface{}) error
	UnFavouritePricings(q map[string]interface{}) error
	GetFavPricingList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error)
	GetArrPricing(user_id int) (pq.Int64Array, error)
}

type serviceBase struct {
	mdmRepository mdm_repo.MdmBase
}

var newServiceObj *serviceBase //singleton object

// singleton function
func NewServiceBase() *serviceBase {
	if newServiceObj != nil {
		return newServiceObj
	}
	mdmRepository := mdm_repo.NewMdmBase()
	newServiceObj = &serviceBase{mdmRepository}
	return newServiceObj
}

//--------------------------------Contacts------------------------------------

func (s *serviceBase) UnFavouriteContacts(q map[string]interface{}) error {
	id, _ := strconv.Atoi(q["ID"].(string))
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.mdmRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.ContactIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "contact_ids",
			"query":       "array_remove(contact_ids,?)",
		}
		err = s.mdmRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) FavouriteContacts(q map[string]interface{}) error {
	//fmt.Println("start")
	id, _ := strconv.Atoi(q["ID"].(string))
	user_id := q["user_id"].(int)
	//fmt.Println(id, user_id)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	samp := record.ContactIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		//fmt.Println("-----------------------------Create Record---------------------------------------------")
		record.UserID = uint(user_id)
		record.ContactIds = pq.Int64Array([]int64{int64(id)}) //append(record.ContactIds, id)
		//fmt.Println("contact id:----------------------------->", record.ContactIds)
		err := s.mdmRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		//fmt.Println("----------------------------Add id in Record---")
		record.ContactIds = append(record.ContactIds, int64(id))
		//fmt.Println("contact id:----------------------------->", record.ContactIds)
		query := map[string]interface{}{
			"user_id":     uint(user_id),
			"contact_ids": record.ContactIds,
		}
		err := s.mdmRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) GetFavContactList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Partner, error) {
	result, err := s.mdmRepository.GetFavContactList(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *serviceBase) GetArrContact(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("user data Not found")
	}
	return record.ContactIds, er
}

//----------------------------------Products---------------------------------------

func (s *serviceBase) FavouriteProducts(q map[string]interface{}) error {
	//fmt.Println("start")
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	//fmt.Println(id, user_id)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	samp := record.ProductIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		//fmt.Println("-----------------------------Create Record---------------------------------------------")
		record.UserID = uint(user_id)
		record.ProductIds = pq.Int64Array([]int64{int64(id)}) //append(record.ContactIds, id)
		//fmt.Println("contact id:----------------------------->", record.ContactIds)
		err := s.mdmRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		//fmt.Println("----------------------------Add id in Record---")
		record.ProductIds = append(record.ProductIds, int64(id))
		//fmt.Println("contact id:----------------------------->", record.ContactIds)
		query := map[string]interface{}{
			"user_id":     uint(user_id),
			"product_ids": record.ProductIds,
		}
		err := s.mdmRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteProducts(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.mdmRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.ProductIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "product_ids",
			"query":       "array_remove(product_ids,?)",
		}
		err = s.mdmRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) GetFavProductList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.ProductVariant, error) {
	result, err := s.mdmRepository.GetFavProductList(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *serviceBase) GetArrProduct(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("user data Not found")
	}
	return record.ProductIds, er
}

//-----------------------------------Locations------------------------------------

func (s *serviceBase) FavouriteLocations(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.mdmRepository.FindOne(query)
	samp := record.LocationIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.LocationIds = pq.Int64Array([]int64{int64(id)})
		err := s.mdmRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.LocationIds = append(record.LocationIds, int64(id))
		query := map[string]interface{}{
			"user_id":      uint(user_id),
			"location_ids": record.LocationIds,
		}
		err := s.mdmRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnFavouriteLocations(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.mdmRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.LocationIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "location_ids",
			"query":       "array_remove(location_ids,?)",
		}
		err = s.mdmRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) GetArrLocation(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("user data Not found")
	}
	return record.LocationIds, er
}

func (s *serviceBase) GetFavLocationList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Locations, error) {
	result, err := s.mdmRepository.GetFavLocationList(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//------------------------------------Vendors----------------------------------------

func (s *serviceBase) FavouriteVendors(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, er := s.mdmRepository.FindOne(query)
	samp := record.VendorIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		record.UserID = uint(user_id)
		record.VendorIds = pq.Int64Array([]int64{int64(id)})
		err := s.mdmRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		record.VendorIds = append(record.VendorIds, int64(id))
		query := map[string]interface{}{
			"user_id":    uint(user_id),
			"vendor_ids": record.VendorIds,
		}
		err := s.mdmRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) UnfavouriteVendors(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.mdmRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.VendorIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "vendor_ids",
			"query":       "array_remove(vendor_ids,?)",
		}
		err = s.mdmRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s *serviceBase) GetFavVendorList(data pq.Int64Array, p *pagination.Paginatevalue) ([]mdm.Vendors, error) {
	result, err := s.mdmRepository.GetFavVendorList(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *serviceBase) GetArrVendor(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("user data Not found")
	}
	return record.VendorIds, er
}

//-------------------------------------Pricing-----------------------------------------

func (s serviceBase) UnFavouritePricings(q map[string]interface{}) error {
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	record, err := s.mdmRepository.FindOne(query)
	if err != nil {
		return errors.New("record not found")
	}
	samp := record.PricingIds
	if helpers.Contains([]int64(samp), int64(id)) {
		query = map[string]interface{}{
			"user_id":     uint(user_id),
			"id":          int64(id),
			"field_value": "pricing_ids",
			"query":       "array_remove(pricing_ids,?)",
		}
		err = s.mdmRepository.UnFavourite(query)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("record not found in favourite list")
}

func (s serviceBase) FavouritePricings(q map[string]interface{}) error {
	//fmt.Println("start")
	id := q["ID"].(int)
	user_id := q["user_id"].(int)
	//fmt.Println(id, user_id)
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	samp := record.PricingIds
	if helpers.Contains([]int64(samp), int64(id)) {
		return errors.New("record is already added in favourites list")
	}
	if er != nil {
		//fmt.Println("-----------------------------Create Record---------------------------------------------")
		record.UserID = uint(user_id)
		record.PricingIds = pq.Int64Array([]int64{int64(id)}) //append(record.PricingIds, id)
		//fmt.Println("Pricing id:----------------------------->", record.PricingIds)
		err := s.mdmRepository.CreateRecord(record)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
	} else {
		//fmt.Println("----------------------------Add id in Record---")
		record.PricingIds = append(record.PricingIds, int64(id))
		//fmt.Println("Pricing id:----------------------------->", record.PricingIds)
		query := map[string]interface{}{
			"user_id":     uint(user_id),
			"pricing_ids": record.PricingIds,
		}
		err := s.mdmRepository.Favourite(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *serviceBase) GetFavPricingList(data pq.Int64Array, p *pagination.Paginatevalue) ([]shared_pricing_and_location.Pricing, error) {
	result, err := s.mdmRepository.GetFavPricingList(data, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *serviceBase) GetArrPricing(user_id int) (pq.Int64Array, error) {
	query := map[string]interface{}{
		"user_id": user_id,
	}
	//fmt.Println("second")
	record, er := s.mdmRepository.FindOne(query)
	if er != nil {
		return nil, errors.New("user data Not found")
	}
	return record.PricingIds, er
}
