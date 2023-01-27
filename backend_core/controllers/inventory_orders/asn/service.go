package asn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/inventory_orders/grn"
	"fermion/backend_core/controllers/inventory_orders/inventory_adjustments"
	"fermion/backend_core/controllers/orders/scrap_orders"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/pagination"
	inventory_orders_repo "fermion/backend_core/internal/repository/inventory_orders"
	returns_repo "fermion/backend_core/internal/repository/returns"
	access_checker "fermion/backend_core/pkg/util/access"
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
	CreateAsn(metaData core.MetaData, data *inventory_orders.ASN) error
	BulkCreateAsn(metaData core.MetaData, data *[]inventory_orders.ASN) error
	UpdateAsn(metaData core.MetaData, data *inventory_orders.ASN) error
	GetAsn(metaData core.MetaData) (interface{}, error)
	GetAllAsn(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	DeleteAsn(metaData core.MetaData) error
	DeleteAsnLines(metaData core.MetaData) error

	SendMailAsn(metaData core.MetaData, q *SendMailAsn) error
	GetAsnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	asnRepository            inventory_orders_repo.Asn
	grnService               grn.Service
	inventoryService         inventory_adjustments.Service
	scrapOrdersService       scrap_orders.Service
	PurchaseReturnRepository returns_repo.PurchaseReturn
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	AsnRepository := inventory_orders_repo.NewAsn()
	newServiceObj = &service{AsnRepository,
		grn.NewService(),
		inventory_adjustments.NewService(),
		scrap_orders.NewService(), returns_repo.NewPurchaseReturn()}
	return newServiceObj
}

func (s *service) CreateAsn(metaData core.MetaData, data *inventory_orders.ASN) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "ASN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create asn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create asn at data level")
	}
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId
	var dtoData AsnRequest
	if dtoData.AutoCreateAsnNumber {
		data.AsnNumber = helpers.GenerateSequence("ASN", fmt.Sprint(metaData.TokenUserId), "asns")
	}
	if dtoData.AutoGenerateReferenceNumber {
		data.ReferenceNumber = helpers.GenerateSequence("REF", fmt.Sprint(metaData.TokenUserId), "asns")
	}

	defaultStatus, err := helpers.GetLookupcodeId("GRN_STATUS", "DRAFT")
	if err != nil {
		return err
	}
	data.StatusID = defaultStatus
	err = s.asnRepository.CreateAsn(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) BulkCreateAsn(metaData core.MetaData, data *[]inventory_orders.ASN) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "ASN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create asn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create asn at data level")
	}
	bulk_data := []map[string]interface{}{}
	var dtoData AsnRequest
	for index, value := range *data {
		v := map[string]interface{}{}
		count := helpers.GetCount("SELECT COUNT(*) FROM asns") + 1 + index
		if dtoData.AutoCreateAsnNumber {
			value.AsnNumber = "ASN/" + fmt.Sprint(metaData.TokenUserId) + "/000" + strconv.Itoa(count)
		}
		if dtoData.AutoGenerateReferenceNumber {
			value.ReferenceNumber = "REF/" + fmt.Sprint(metaData.TokenUserId) + "/000" + strconv.Itoa(count)
		}

		val, err := json.Marshal(value)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(val, &v)
		if err != nil {
			return res.BuildError(res.ErrUnprocessableEntity, err)
		}

		bulk_data = append(bulk_data, v)
	}
	var AsnData []inventory_orders.ASN
	dto, err := json.Marshal(bulk_data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &AsnData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	err = s.asnRepository.BulkCreateAsn(&AsnData)
	return err
}

func (s *service) UpdateAsn(metaData core.MetaData, data *inventory_orders.ASN) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "ASN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update asn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update asn at data level")
	}
	var find_query = map[string]interface{}{
		"id": metaData.Query["id"],
	}
	old_data, er := s.asnRepository.GetAsn(find_query)
	if er != nil {
		return er
	}
	data.UpdatedByID = &metaData.TokenUserId
	old_status := old_data.StatusID
	new_status := data.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusID)
		data.StatusHistory = result
	}

	err := s.asnRepository.UpdateAsn(metaData.Query, data)
	if er != nil {
		return err
	}
	for _, order_line := range data.AsnOrderLines {
		query := map[string]interface{}{
			"asn_id":     old_data.ID,
			"product_id": order_line.ProductID,
		}
		count, er := s.asnRepository.UpdateAsnLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.CreatedByID = &metaData.TokenUserId
			order_line.CompanyId = metaData.CompanyId
			order_line.UpdatedByID = nil
			order_line.Asn_id = metaData.TokenUserId
			e := s.asnRepository.CreateAsnLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (s *service) GetAsn(metaData core.MetaData) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "ASN", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view asn at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view asn at data level")
	}
	result, er := s.asnRepository.GetAsn(metaData.Query)
	if er != nil {
		return result, er
	}
	// result_order_lines, err := s.asnRepository.GetAsnLines(metaData.Query)
	// result.AsnOrderLines = result_order_lines
	// if err != nil {
	// 	return result, err
	// }

	return result, nil
}

func (s *service) GetAllAsn(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "ASN", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list asn at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list asn at data level")
	}
	result, err := s.asnRepository.GetAllAsn(metaData.Query, p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteAsn(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "ASN", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete asn at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete asn at data level")
	}
	err := s.asnRepository.DeleteAsn(metaData.Query)
	if err != nil {
		return err
	} else {
		query := map[string]interface{}{"asn_id": metaData.Query["id"]}
		er := s.asnRepository.DeleteAsnLines(query)
		if er != nil {
			return er
		}
	}
	return nil
}

func (s *service) DeleteAsnLines(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	access_module_flag, data_access := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "ASN", metaData.TokenUserId)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete asn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete asn at data level")
	}
	err := s.asnRepository.DeleteAsnLines(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SendMailAsn(metaData core.MetaData, q *SendMailAsn) error {
	query := map[string]interface{}{
		"id1": q.ID,
	}
	result, er := s.asnRepository.GetAsn(query)
	if er != nil {
		return er
	}
	err := helpers.SendEmail("ASN", q.ReceiverEmail, "pkg/util/response/static/inventory_orders/asn_template.html", result)

	return err
}
func (s *service) GetAsnTab(metaData core.MetaData, page *pagination.Paginatevalue) (interface{}, error) {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "LIST", "ASN", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list asn at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list asn at data level")
	}

	tab := metaData.AdditionalFields["tab"].(string)
	asnId := metaData.AdditionalFields["id"]

	if tab == "grn" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("GRN_SOURCE_DOCUMENT_TYPES", "ASN")
		grnPage := page
		grnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, asnId)
		data, err := s.grnService.GetAllGRN(metaData, grnPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "inventory_adjustments" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("INVENTORY_ADJUSTMENTS_SOURCE_DOCUMENT_TYPES", "ASN")

		inventoryPage := page
		inventoryPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, asnId)
		data, err := s.inventoryService.GetAllInvAdj(metaData, inventoryPage)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "purchase_returns" {
		asnPagination := new(pagination.Paginatevalue)
		source_document_type_id, _ := helpers.GetLookupcodeId("ASN_SOURCE_DOCUMENT_TYPES", "PURCHASE_RETURNS")
		asnPagination.Filters = fmt.Sprintf("[[\"id\",\"=\",%v],[\"source_document_type_id\",\"=\",%v]]", asnId, source_document_type_id)
		asn_interface, err := s.GetAllAsn(metaData, asnPagination)
		if err != nil {
			return nil, err
		}
		asn := asn_interface.([]inventory_orders.ASN)
		if len(asn) == 0 {
			return nil, errors.New("no source document")
		}

		var asnSourceDoc map[string]interface{}
		asnSourceDocJson := asn[0].SourceDocuments
		dto, err := json.Marshal(asnSourceDocJson)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}
		err = json.Unmarshal(dto, &asnSourceDoc)
		if err != nil {
			return 0, res.BuildError(res.ErrUnprocessableEntity, err)
		}

		asnSourceDocId := asnSourceDoc["id"]
		if asnSourceDocId == nil {
			return nil, errors.New("no source document")
		}

		PurchaseReturnPage := page
		PurchaseReturnPage.Filters = fmt.Sprintf("[[\"id\",\"=\",%v]]", asnSourceDoc["id"])
		data, err := s.PurchaseReturnRepository.FindAll(nil, page)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	return nil, nil
}
