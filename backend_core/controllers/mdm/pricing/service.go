package pricing

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"
	"fermion/backend_core/internal/model/pagination"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
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
	CreatePricing(metaData core.MetaData, data *shared_pricing_and_location.Pricing) error
	CreateSalesPricing(metaData core.MetaData, data *shared_pricing_and_location.SalesPriceList) error
	CreatePurchasePricing(metaData core.MetaData, data *shared_pricing_and_location.PurchasePriceList) error
	CreateTransferPricing(metaData core.MetaData, data *shared_pricing_and_location.TransferPriceList) error

	UpdatePricing(metaData core.MetaData, data interface{}) error
	FindPricing(metaData core.MetaData) (interface{}, error)
	ListPricing(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeletePricing(metaData core.MetaData) error
	DeleteLineItems(metaData core.MetaData) error
	ListPurchasePricing(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	ListSalesPiceLineItems(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)

	//-------Channel Pricing Upsert----------------------------------------------------
	ChannelSalesPriceUpsert(metaData core.MetaData, data []interface{}) (interface{}, error)
	ChannelSalesPriceListLineItemsUpsert(metaData core.MetaData, data []shared_pricing_and_location.SalesLineItems) (interface{}, error)
}

type service struct {
	pricingRepository mdm_repo.Pricing
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{mdm_repo.NewPricing()}
	return newServiceObj
}

func (s *service) CreatePricing(metaData core.MetaData, data *shared_pricing_and_location.Pricing) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create uom at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create uom at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId
	defaultStatus, err := helpers.GetLookupcodeId("PRICING_STATUS", "ACTIVE")
	if err != nil {
		return err
	}
	data.StatusId = &defaultStatus

	err = s.pricingRepository.CreatePricing(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) CreateSalesPricing(metaData core.MetaData, data *shared_pricing_and_location.SalesPriceList) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pricing at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.pricingRepository.CreateSalePriceList(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) CreatePurchasePricing(metaData core.MetaData, data *shared_pricing_and_location.PurchasePriceList) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pricing at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.pricingRepository.CreatePurchasePriceList(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) CreateTransferPricing(metaData core.MetaData, data *shared_pricing_and_location.TransferPriceList) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create pricing at data level")
	}

	data.CompanyId = metaData.CompanyId
	data.CreatedByID = &metaData.TokenUserId

	err := s.pricingRepository.CreateTransferPriceList(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) FindPricing(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "READ", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view pricing at data level")
	}
	response, err := s.pricingRepository.FindPricing(metaData.Query)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s *service) ListPricing(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pricing at data level")
	}

	result, err := s.pricingRepository.ListPricing(metaData.Query, p)
	if err != nil {
		return result, err
	}
	return result, nil
}
func (s *service) UpdatePricing(metaData core.MetaData, data interface{}) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update pricing at data level")
	}

	pricingDetails, err := s.pricingRepository.FindPricing(metaData.Query)
	if err != nil {
		return err
	}

	pricingPayload := new(shared_pricing_and_location.Pricing)
	helpers.JsonMarshaller(data, pricingPayload)

	pricingPayload.UpdatedByID = &metaData.TokenUserId
	err = s.pricingRepository.UpdatePriceList(metaData.Query, pricingPayload)
	if err != nil {
		return err
	}
	fmt.Println("pricing updated successfully")

	if pricingPayload.PriceListId == 1 {
		salesPriceList := new(shared_pricing_and_location.SalesPriceList)
		helpers.JsonMarshaller(data, salesPriceList)

		priceListquery := map[string]interface{}{
			"id":         *pricingDetails.SalesPriceListId,
			"company_id": metaData.CompanyId,
		}
		salesPriceList.UpdatedByID = &metaData.TokenUserId
		err := s.pricingRepository.UpdateSalesPriceList(priceListquery, salesPriceList)
		if err != nil {
			return err
		}
		fmt.Println("sales_price_list updated successfully")

		for _, line_items := range salesPriceList.SalesLineItems {
			lineItemQuery := map[string]interface{}{
				"sales_price_list_id": *pricingDetails.SalesPriceListId,
				"product_variant_id":  *line_items.ProductVariantId,
			}
			salesPriceLineItemDetails, _ := s.pricingRepository.FindOneSalesPriceLineItem(lineItemQuery)
			if salesPriceLineItemDetails.ID != 0 {
				line_items.UpdatedByID = &metaData.TokenUserId
				err := s.pricingRepository.UpdateSalesLineItems(lineItemQuery, &line_items)
				if err != nil {
					return err
				}
				fmt.Println("sales_price_list-line_item updated successfully")
			} else {
				line_items.CreatedByID = &metaData.TokenUserId
				line_items.SalesPriceListId = *pricingDetails.SalesPriceListId
				err := s.pricingRepository.CreateSalesLineItems(&line_items)
				if err != nil {
					return err
				}
				fmt.Println("sales_price_list-line_item created successfully")
			}
		}
		return nil
	}
	if pricingPayload.PriceListId == 2 {
		purchasePriceList := new(shared_pricing_and_location.PurchasePriceList)
		helpers.JsonMarshaller(data, purchasePriceList)

		priceListquery := map[string]interface{}{
			"id":         *pricingDetails.PurchasePriceListId,
			"company_id": metaData.CompanyId,
		}
		purchasePriceList.UpdatedByID = &metaData.TokenUserId
		err := s.pricingRepository.UpdatePurchasePriceList(priceListquery, purchasePriceList)
		if err != nil {
			return err
		}
		fmt.Println("purchase_price_list updated successfully")

		for _, line_items := range purchasePriceList.PurchaseLineItems {
			lineItemQuery := map[string]interface{}{
				"purchase_price_list_id": *pricingDetails.PurchasePriceListId,
				"product_variant_id":     *line_items.ProductVariantId,
			}
			purchasePriceLineItemDetails, _ := s.pricingRepository.FindOnePurchaseLineItem(lineItemQuery)
			if purchasePriceLineItemDetails.ID != 0 {
				line_items.UpdatedByID = &metaData.TokenUserId
				err := s.pricingRepository.UpdatePurchaseLineItems(lineItemQuery, &line_items)
				if err != nil {
					return err
				}
				fmt.Println("purchase_price_list-line_item updated successfully")
			} else {
				line_items.CreatedByID = &metaData.TokenUserId
				line_items.PurchasePriceListId = *pricingDetails.PurchasePriceListId
				err := s.pricingRepository.CreatePurchaseLineItems(&line_items)
				if err != nil {
					return err
				}
				fmt.Println("purchase_price_list-line_item created successfully")
			}
		}
		return nil
	}
	if pricingPayload.PriceListId == 3 {
		transferPriceList := new(shared_pricing_and_location.TransferPriceList)
		helpers.JsonMarshaller(data, transferPriceList)

		priceListquery := map[string]interface{}{
			"id":         pricingDetails.TransferPriceListId,
			"company_id": metaData.CompanyId,
		}
		transferPriceList.UpdatedByID = &metaData.TokenUserId
		err := s.pricingRepository.UpdateTransferPriceList(priceListquery, transferPriceList)
		if err != nil {
			return err
		}
		fmt.Println("transfer_price_list updated successfully")

		for _, line_items := range transferPriceList.TransferLineItems {
			lineItemQuery := map[string]interface{}{
				"transfer_price_list_id": *pricingDetails.TransferPriceListId,
				"product_variant_id":     *line_items.ProductVariantId,
			}
			transferPriceLineItemDetails, _ := s.pricingRepository.FindOneTransferLineItem(lineItemQuery)
			if transferPriceLineItemDetails.ID != 0 {
				line_items.UpdatedByID = &metaData.TokenUserId
				err := s.pricingRepository.UpdateTransferLineItems(lineItemQuery, &line_items)
				if err != nil {
					return err
				}
				fmt.Println("transfer_price_list-line_item updated successfully")
			} else {
				line_items.CreatedByID = &metaData.TokenUserId
				line_items.TransferPriceListId = *pricingDetails.TransferPriceListId
				err := s.pricingRepository.CreateTransferLineItems(&line_items)
				if err != nil {
					return err
				}
				fmt.Println("transfer_price_list-line_item created successfully")
			}
		}
		return nil
	}
	return nil
}
func (s *service) DeletePricing(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pricing at data level")
	}

	query := map[string]interface{}{
		"id":         metaData.Query["id"],
		"company_id": metaData.CompanyId,
	}
	pricing, _ := s.pricingRepository.FindPricing(query)
	if pricing.ID == 0 {
		return errors.New("oops! record not found")
	}
	if pricing.PriceListId == 1 {
		query := map[string]interface{}{
			"id":      *pricing.SalesPriceListId,
			"user_id": metaData.TokenUserId,
		}
		err := s.pricingRepository.DeleteSalesPriceList(query)
		if err != nil {
			return err
		}
		fmt.Println("sales_price_list deleted successfully")
		lineItemQuery := map[string]interface{}{
			"sales_price_list_id": *pricing.SalesPriceListId,
			"user_id":             metaData.TokenUserId,
		}
		err = s.pricingRepository.DeleteSalesLineItems(lineItemQuery)
		if err != nil {
			return err
		}
		fmt.Println("sales_price_list-line_item deleted successfully")
	} else if pricing.PriceListId == 2 {
		query := map[string]interface{}{
			"id":      *pricing.PurchasePriceListId,
			"user_id": metaData.TokenUserId,
		}
		err := s.pricingRepository.DeletePurchasePriceList(query)
		if err != nil {
			return err
		}
		fmt.Println("purchase_price_list deleted successfully")
		lineItemQuery := map[string]interface{}{
			"purchase_price_list_id": *pricing.PurchasePriceListId,
			"user_id":                metaData.TokenUserId,
		}
		err = s.pricingRepository.DeletePurchaseLineItems(lineItemQuery)
		if err != nil {
			return err
		}
		fmt.Println("purchase_price_list-line_item deleted successfully")
	} else if pricing.PriceListId == 3 {
		query := map[string]interface{}{
			"id":      *pricing.TransferPriceListId,
			"user_id": metaData.TokenUserId,
		}
		err := s.pricingRepository.DeleteTransferPriceList(query)
		if err != nil {
			return err
		}
		fmt.Println("transfer_price_list deleted successfully")
		lineItemQuery := map[string]interface{}{
			"transfer_price_list_id": *pricing.TransferPriceListId,
			"user_id":                metaData.TokenUserId,
		}
		err = s.pricingRepository.DeleteTransferLineItems(lineItemQuery)
		if err != nil {
			return err
		}
		fmt.Println("transfer_price_list-line_item deleted successfully")
	}

	err := s.pricingRepository.DeletePricing(metaData.Query)
	if err != nil {
		return err
	}
	fmt.Println("pricing deleted successfully")
	return nil
}
func (s *service) DeleteLineItems(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete pricing at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete pricing at data level")
	}

	productVariantId := metaData.AdditionalFields["product_variant_id"]

	query := map[string]interface{}{
		"id":         metaData.Query["id"],
		"company_id": metaData.CompanyId,
	}

	pricing, _ := s.pricingRepository.FindPricing(query)
	if pricing.ID == 0 {
		return errors.New("oops! record not found")
	}
	if pricing.PriceListId == 1 {
		lineItemQuery := map[string]interface{}{
			"sales_price_list_id": *pricing.SalesPriceListId,
			"product_variant_id":  productVariantId,
			"user_id":             metaData.TokenUserId,
		}
		err := s.pricingRepository.DeleteSalesLineItems(lineItemQuery)
		if err != nil {
			return err
		}
	} else if pricing.PriceListId == 2 {
		lineItemQuery := map[string]interface{}{
			"purchase_price_list_id": *pricing.PurchasePriceListId,
			"product_variant_id":     productVariantId,
			"user_id":                metaData.TokenUserId,
		}
		err := s.pricingRepository.DeletePurchaseLineItems(lineItemQuery)
		if err != nil {
			return err
		}
	} else if pricing.PriceListId == 3 {
		lineItemQuery := map[string]interface{}{
			"transfer_price_list_id": *pricing.TransferPriceListId,
			"product_variant_id":     productVariantId,
			"user_id":                metaData.TokenUserId,
		}
		err := s.pricingRepository.DeleteTransferLineItems(lineItemQuery)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *service) ListPurchasePricing(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pricing at data level")
	}
	result, err := s.pricingRepository.ListPurchasePricing(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (s *service) ListSalesPiceLineItems(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "PRICING", metaData.TokenUserId)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list pricing at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list pricing at data level")
	}
	lineItems, err := s.pricingRepository.ListSalesPiceLineItems(metaData.Query, p)
	return lineItems, err
}

// ----------------Channel Pricing UPsert --------------------------------
func (s *service) ChannelSalesPriceUpsert(metaData core.MetaData, data []interface{}) (interface{}, error) {
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

		pricing.UpdatedByID = &metaData.TokenUserId
		if pricing.PriceListName == "" {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": "pricelist name shouldn't be empty"})
			continue
		}
		priceListQuery := map[string]interface{}{
			"price_list_name": pricing.PriceListName,
			"company_id":      metaData.CompanyId,
		}
		result, err := s.pricingRepository.FindPricing(priceListQuery)

		if result.PriceListName == "" {
			failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			continue
		}
		if result.SalesPriceListId != nil {
			metaData.AdditionalFields = map[string]interface{}{
				"sales_price_list_id": *result.SalesPriceListId,
			}
			result, err := s.ChannelSalesPriceListLineItemsUpsert(metaData, pricing.SalesPriceList.SalesLineItems)
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

func (s *service) ChannelSalesPriceListLineItemsUpsert(metaData core.MetaData, data []shared_pricing_and_location.SalesLineItems) (interface{}, error) {
	var success []interface{}
	var failures []interface{}

	SalesPriceListId := metaData.AdditionalFields["sales_price_list_id"]

	if SalesPriceListId == 0 {
		return nil, errors.New("missing sales price list id")
	}
	for index, lineItem := range data {
		query := map[string]interface{}{
			"sales_price_list_id": SalesPriceListId,
			"product_variant_id":  lineItem.ProductVariantId,
		}
		result, _ := s.pricingRepository.FindOneSalesPriceLineItem(query)
		if result.ProductVariantId != nil && *result.ProductVariantId != 0 {
			lineItem.UpdatedByID = &metaData.TokenUserId
			err := s.pricingRepository.UpdateSalesLineItems(query, &lineItem)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			success = append(success, map[string]interface{}{"status": true, "serial_number": index + 1, "msg": "updated"})
			continue
		}
		lineItem.SalesPriceListId = SalesPriceListId.(uint)
		lineItem.CreatedByID = &metaData.TokenUserId
		err := s.pricingRepository.CreateSalesLineItems(&lineItem)
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
