package pricing

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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
	CreatePurchasePricing(pricing *shared_pricing_and_location.Pricing, purchase *shared_pricing_and_location.PurchasePriceList, token_id string, access_template_id string) error
	CreateTransferPricing(pricing *shared_pricing_and_location.Pricing, transfer *shared_pricing_and_location.TransferPriceList, token_id string, access_template_id string) error
	CreateSalesPricing(pricing *shared_pricing_and_location.Pricing, sales *shared_pricing_and_location.SalesPriceList, token_id string, access_template_id string) error
	UpdatePricing(id uint, reqBytes []byte, c echo.Context, token_id string, access_template_id string) error
	PricingByid(id uint, token_id string, access_template_id string) (interface{}, error)
	ListPricing(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]PricingDTO, error)
	DeletePricing(id uint, user_id uint, token_id string, access_template_id string) error
	DeleteLineItems(id uint, product_id uint, token_id string, access_template_id string) error
	ListPurchasePricing(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]shared_pricing_and_location.PurchasePriceList, error)
	ListSalesPiceLineItems(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]shared_pricing_and_location.SalesLineItems, error)

	//-------Channel Pricing Upsert----------------------------------------------------
	ChannelSalesPriceUpsert(data []interface{}, TokenUserId string) (interface{}, error)
	ChannelSalesPriceListLineItemsUpsert(data []shared_pricing_and_location.SalesLineItems, SalesPriceListId uint) (interface{}, error)
}

type service struct {
	pricingRepository mdm_repo.Pricing
	db                *gorm.DB
}

func NewService() *service {
	db := db.DbManager()
	return &service{mdm_repo.NewPricing(), db}
}

func (s *service) CreatePurchasePricing(pricing *shared_pricing_and_location.Pricing, purchase *shared_pricing_and_location.PurchasePriceList, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRICING", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pricing at data level")
	}
	tx := s.db.Begin(&sql.TxOptions{})
	purchase_err := s.pricingRepository.CreatePurchasePriceList(tx, purchase)
	if purchase_err != nil {
		tx.Rollback()
		return purchase_err
	}
	pricing.PurchasePriceListId = &purchase.ID
	defaultStatus, err := helpers.GetLookupcodeId("PRICING_STATUS", "ACTIVE")
	if err != nil {
		return err
	}
	pricing.StatusId = &defaultStatus
	pricing_err := s.pricingRepository.CreatePricing(tx, pricing)
	if pricing_err != nil {
		tx.Rollback()
		return pricing_err
	}
	tx.Commit()
	return nil
}

func (s *service) CreateTransferPricing(pricing *shared_pricing_and_location.Pricing, transfer *shared_pricing_and_location.TransferPriceList, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRICING", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pricing at data level")
	}
	tx := s.db.Begin(&sql.TxOptions{})
	transfer_err := s.pricingRepository.CreateTransferPriceList(tx, transfer)
	if transfer_err != nil {
		tx.Rollback()
		return transfer_err
	}
	pricing.TransferPriceListId = &transfer.ID
	defaultStatus, err := helpers.GetLookupcodeId("PRICING_STATUS", "ACTIVE")
	if err != nil {
		return err
	}
	pricing.StatusId = &defaultStatus
	pricing_err := s.pricingRepository.CreatePricing(tx, pricing)
	if pricing_err != nil {
		tx.Rollback()
		return pricing_err
	}
	tx.Commit()
	return nil
}

func (s *service) CreateSalesPricing(pricing *shared_pricing_and_location.Pricing, sales *shared_pricing_and_location.SalesPriceList, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "PRICING", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pricing at data level")
	}
	tx := s.db.Begin(&sql.TxOptions{})
	sales_err := s.pricingRepository.CreateSalePriceList(tx, sales)
	if sales_err != nil {
		tx.Rollback()
		return sales_err
	}
	pricing.SalesPriceListId = &sales.ID
	defaultStatus, err := helpers.GetLookupcodeId("PRICING_STATUS", "ACTIVE")
	if err != nil {
		return err
	}
	pricing.StatusId = &defaultStatus
	pricing_err := s.pricingRepository.CreatePricing(tx, pricing)
	if pricing_err != nil {
		tx.Rollback()
		return pricing_err
	}
	tx.Commit()
	return nil
}

func (s *service) PricingByid(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "PRICING", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view pricing at data level")
	}
	var data PricingDTO
	query := map[string]interface{}{
		"id": id,
	}
	pricing, err := s.pricingRepository.FindPricing(query)

	dto, _ := json.Marshal(pricing)
	_ = json.Unmarshal(dto, &data)
	return data, err
}

func (s *service) ListPricing(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]PricingDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "PRICING", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pricing at data level")
	}
	var pricingData []PricingDTO
	pricing_response, err := s.pricingRepository.ListPricing(p)
	fmt.Println(pricing_response)

	mdata, _ := json.Marshal(pricing_response)
	json.Unmarshal(mdata, &pricingData)

	return pricingData, err
}

func (s *service) UpdatePricing(id uint, reqBytes []byte, c echo.Context, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "PRICING", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update pricing at data level")
	}
	query := map[string]interface{}{
		"id": id,
	}
	pricing, err := s.pricingRepository.FindPricing(query)
	if err != nil {
		return err
	}
	pricing_data := c.Get("pricing").(*shared_pricing_and_location.Pricing)
	json.Unmarshal(reqBytes, pricing_data)
	st := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(st)
	t := uint(user_id)
	pricing_data.UpdatedByID = &t
	tx := s.db.Begin(&sql.TxOptions{})
	pricing_er := s.pricingRepository.UpdatePriceList(tx, id, pricing_data)
	if pricing_er != nil {
		tx.Rollback()
		return pricing_er
	}
	if pricing_data.Price_list_id == 1 {
		// sales_data := new(shared_pricing_and_location.SalesPriceList)
		// json.Unmarshal(reqBytes, sales_data)
		sales_data := c.Get("sales").(*shared_pricing_and_location.SalesPriceList)
		sales_err := s.pricingRepository.UpdateSalesPriceList(tx, *pricing.SalesPriceListId, sales_data)
		if sales_err != nil {
			return sales_err
		}
		for _, line_items := range sales_data.SalesLineItems {
			query := map[string]interface{}{
				"spl_id":     *pricing.SalesPriceListId,
				"product_id": line_items.Product_id,
			}
			count, er1 := s.pricingRepository.UpdateSalesLineItems(tx, query, &line_items)
			if er1 != nil {
				tx.Rollback()
				return er1
			} else if count == 0 {
				line_items.SPL_id = *pricing.SalesPriceListId
				er2 := s.pricingRepository.CreateSalesLineItems(tx, &line_items)
				if er2 != nil {
					tx.Rollback()
					return er2
				}
			}
		}
	} else if pricing_data.Price_list_id == 2 {
		purchase_data := c.Get("purchase").(*shared_pricing_and_location.PurchasePriceList)
		json.Unmarshal(reqBytes, purchase_data)
		purchase_err := s.pricingRepository.UpdatePurchasePriceList(tx, *pricing.PurchasePriceListId, purchase_data)
		if purchase_err != nil {
			return purchase_err
		}

		for _, line_items := range purchase_data.PurchaseLineItems {
			query := map[string]interface{}{
				"ppl_id":     *pricing.PurchasePriceListId,
				"product_id": line_items.Product_id,
			}
			count, er1 := s.pricingRepository.UpdatePurchaseLineItems(tx, query, &line_items)
			if er1 != nil {
				tx.Rollback()
				return er1
			} else if count == 0 {
				line_items.PPL_id = *pricing.PurchasePriceListId
				er2 := s.pricingRepository.CreatePurchaseLineItems(tx, &line_items)
				if er2 != nil {
					tx.Rollback()
					return er2
				}
			}
		}
	} else if pricing_data.Price_list_id == 3 {
		transfer_data := c.Get("transfer").(*shared_pricing_and_location.TransferPriceList)
		json.Unmarshal(reqBytes, transfer_data)
		transfer_err := s.pricingRepository.UpdateTransferPriceList(tx, *pricing.TransferPriceListId, transfer_data)
		if transfer_err != nil {
			return transfer_err
		}

		for _, line_items := range transfer_data.TransferLineItems {
			query := map[string]interface{}{
				"tpl_id":     *pricing.TransferPriceListId,
				"product_id": line_items.Product_id,
			}
			count, er1 := s.pricingRepository.UpdateTransferLineItems(tx, query, &line_items)
			if er1 != nil {
				tx.Rollback()
				return er1
			} else if count == 0 {
				line_items.TPL_id = *pricing.TransferPriceListId
				er2 := s.pricingRepository.CreateTransferLineItems(tx, &line_items)
				if er2 != nil {
					tx.Rollback()
					return er2
				}
			}
		}
	} else {
		return fmt.Errorf("invalid pricelist id")
	}
	tx.Commit()
	return nil
}

func (s *service) DeletePricing(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRICING", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pricing at data level")
	}
	query := map[string]interface{}{
		"id": id,
	}
	pricing, err := s.pricingRepository.FindPricing(query)
	if err != nil {
		return err
	}
	tx := s.db.Begin(&sql.TxOptions{})
	if pricing.Price_list_id == 1 {
		er1 := s.pricingRepository.DeleteSalesPriceList(tx, *pricing.SalesPriceListId)
		if er1 != nil {
			tx.Rollback()
			return er1
		}
		query := map[string]interface{}{
			"spl_id": *pricing.SalesPriceListId,
		}
		er2 := s.pricingRepository.DeleteSalesLineItems(tx, query)
		if er2 != nil {
			tx.Rollback()
			return er2
		}
	} else if pricing.Price_list_id == 2 {
		er1 := s.pricingRepository.DeletePurchasePriceList(tx, *pricing.PurchasePriceListId)
		if er1 != nil {
			tx.Rollback()
			return er1
		}
		query := map[string]interface{}{
			"ppl_id": *pricing.PurchasePriceListId,
		}
		er2 := s.pricingRepository.DeletePurchaseLineItems(tx, query)
		if er2 != nil {
			tx.Rollback()
			return er2
		}
	} else if pricing.Price_list_id == 3 {
		er1 := s.pricingRepository.DeleteTransferPriceList(tx, *pricing.TransferPriceListId)
		if er1 != nil {
			tx.Rollback()
			return er1
		}
		query := map[string]interface{}{
			"tpl_id": *pricing.TransferPriceListId,
		}
		er2 := s.pricingRepository.DeleteTransferLineItems(tx, query)
		if er2 != nil {
			tx.Rollback()
			return er2
		}
	} else {
		return fmt.Errorf("invalid pricelist id")
	}

	pricing_err := s.pricingRepository.DeletePricing(tx, id, user_id)
	if pricing_err != nil {
		tx.Rollback()
		return pricing_err
	}

	tx.Commit()
	return nil
}

func (s *service) DeleteLineItems(id uint, product_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "PRICING", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pricing at data level")
	}
	query := map[string]interface{}{
		"id": id,
	}
	pricing, err := s.pricingRepository.FindPricing(query)
	if err != nil {
		return err
	}
	tx := s.db.Begin(&sql.TxOptions{})
	if pricing.Price_list_id == 1 {
		query := map[string]interface{}{"spl_id": *pricing.SalesPriceListId, "product_id": product_id}
		er := s.pricingRepository.DeleteSalesLineItems(tx, query)
		if er != nil {
			tx.Rollback()
			return er
		}
	} else if pricing.Price_list_id == 2 {
		query := map[string]interface{}{"ppl_id": *pricing.PurchasePriceListId, "product_id": product_id}
		er := s.pricingRepository.DeletePurchaseLineItems(tx, query)
		if er != nil {
			tx.Rollback()
			return er
		}
	} else if pricing.Price_list_id == 3 {
		query := map[string]interface{}{"tpl_id": *pricing.TransferPriceListId, "product_id": product_id}
		er := s.pricingRepository.DeleteTransferLineItems(tx, query)
		if er != nil {
			tx.Rollback()
			return er
		}
	} else {
		return fmt.Errorf("invalid pricelist id")
	}
	tx.Commit()
	return nil
}

func (s *service) ListPurchasePricing(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]shared_pricing_and_location.PurchasePriceList, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PRICING", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pricing at data level")
	}
	pricing, err := s.pricingRepository.ListPurchasePricing(p)
	return pricing, err
}

func (s *service) ListSalesPiceLineItems(p *pagination.Paginatevalue, token_id string, access_template_id string) ([]shared_pricing_and_location.SalesLineItems, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "PRICING", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pricing at data level")
	}
	lineItems, err := s.pricingRepository.ListSalesPiceLineItems(p)
	return lineItems, err
}

//----------------Channel Pricing UPsert --------------------------------

func (s *service) ChannelSalesPriceUpsert(data []interface{}, TokenUserId string) (interface{}, error) {
	var success []interface{}
	var failures []interface{}

	for index, payload := range data {
		var pricing shared_pricing_and_location.Pricing

		jsonPayloadByteArray, err := json.Marshal(payload)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		json.Unmarshal(jsonPayloadByteArray, &pricing)
		// if err != nil {
		// 	failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
		// 	continue
		// }

		pricing.UpdatedByID = helpers.ConvertStringToUint(TokenUserId)
		if pricing.Price_list_name == "" {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "pricelist name shouldn't be empty"})
			continue
		}
		priceListQuery := map[string]interface{}{
			"price_list_name": pricing.Price_list_name,
		}
		result, err := s.pricingRepository.FindPricing(priceListQuery)

		if result.Price_list_name == "" {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		if result.SalesPriceListId != nil {
			result, err := s.ChannelSalesPriceListLineItemsUpsert(pricing.SalesPriceList.SalesLineItems, *result.SalesPriceListId)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			success = append(success, map[string]interface{}{"status": true, "serial_number": index + 1, "msg": "updated", "sales_price_line_items": result})
			continue
		}
		failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "sales_price_list_does_not_exist"})
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
}

func (s *service) ChannelSalesPriceListLineItemsUpsert(data []shared_pricing_and_location.SalesLineItems, SalesPriceListId uint) (interface{}, error) {
	var success []interface{}
	var failures []interface{}

	if SalesPriceListId == 0 {
		return nil, errors.New("missing sales price list id")
	}
	for index, lineItem := range data {
		query := map[string]interface{}{
			"spl_id":     SalesPriceListId,
			"product_id": lineItem.Product_id,
		}
		result, _ := s.pricingRepository.FindOneSalesPriceLineItem(query)
		if result.Product_id != 0 {
			_, err := s.pricingRepository.UpdateSalesLineItems(s.db, query, &lineItem)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			success = append(success, map[string]interface{}{"status": true, "serial_number": index + 1, "msg": "updated"})
			continue
		}
		lineItem.SPL_id = SalesPriceListId
		err := s.pricingRepository.CreateSalesLineItems(s.db, &lineItem)
		if err != nil {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		success = append(success, map[string]interface{}{"status": true, "serial_number": index + 1, "msg": "created"})
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
}
