package offers

import (
	"time"

	products_model "fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/offers"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	repository "fermion/backend_core/internal/repository/offers"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"
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

type Service interface {

	//Ondc Apps
	ListOffers(page *pagination.Paginatevalue, token_id string) (interface{}, error)
	ViewOffers(query map[string]interface{}, token_id string) (interface{}, error)
	CreateOffers(data *offers.Offers, token_id string, company_id string) error
	UpdateOffers(query map[string]interface{}, data *offers.Offers, token_id string, company_id string) error
	DeleteOffers(query map[string]interface{}, token_id string) error
	DeleteOffersLines(query map[string]interface{}, token_id string) error
}

type service struct {
	offersRepository  repository.Offer
	productRepository mdm_repo.Products
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		offersRepository:  repository.NewOffers(),
		productRepository: mdm_repo.NewProducts(),
	}
	return newServiceObj
}

// func NewService() *service {

// 	offersRepository := repository.NewOffers()
// 	return &service{offersRepository}
// }

// ONDC apps
func (s *service) ListOffers(page *pagination.Paginatevalue, token_id string) (interface{}, error) {
	results, err := s.offersRepository.ListOffers(page)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return results, err
}
func (s *service) ViewOffers(query map[string]interface{}, token_id string) (interface{}, error) {
	data, err := s.offersRepository.ViewOffers(query)

	if err != nil && err.Error() == "record not found" {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}

	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return data, nil
}
func (s *service) CreateOffers(data *offers.Offers, token_id string, company_id string) error {
	data.CompanyId = *helpers.ConvertStringToUint(company_id)
	err := s.offersRepository.CreateOffers(data)
	if err != nil {
		return res.BuildError(res.ErrServerError, err)
	}
	//-----------------update offers in products----------------//
	product := new(products_model.ProductVariant)
	for _, offer_line := range data.OfferProductDetails {
		product.OfferDetails.OfferId = &offer_line.OfferId
		product.OfferDetails.DiscountTypeId = offer_line.DiscountTypeID
		product.OfferDetails.Title = data.PromotionalDetails.Name
		product.OfferDetails.Value = data.DiscountValue
		product.OfferDetails.ValidTo = time.Time(data.PromotionalDetails.EndDateAndTime)
		product.OfferDetails.ValidFrom = time.Time(data.PromotionalDetails.StartDateAndTime)

		query := map[string]interface{}{
			"id": offer_line.ProductId,
		}

		err = s.productRepository.UpdateVariant(product, query)
		if err != nil {
		}
	}
	// //updating the offer in products

	return nil
}
func (s *service) UpdateOffers(query map[string]interface{}, data *offers.Offers, token_id string, company_id string) error {
	_, err := s.offersRepository.ViewOffers(query)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	err = s.offersRepository.UpdateOffers(query, data)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	for _, offer_line := range data.OfferProductDetails {
		updated_query := map[string]interface{}{
			"product_id": offer_line.ProductId,
			"offer_id":   query["id"],
		}
		count, er := s.offersRepository.UpdateOfferLines(updated_query, &offer_line)
		if er != nil {
			return er
		} else if count == 0 {
			offer_line.OfferId = uint(query["id"].(uint))
			e := s.offersRepository.SaveOfferLines(&offer_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}
func (s *service) DeleteOffers(query map[string]interface{}, token_id string) error {
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.offersRepository.ViewOffers(q)
	if er != nil {
		return er
	}
	err := s.offersRepository.DeleteOffers(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteOffersLines(query map[string]interface{}, token_id string) error {
	err := s.offersRepository.DeleteOfferLines(query)
	if err != nil {
		return err
	}
	return nil
}
