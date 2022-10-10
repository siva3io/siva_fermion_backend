package asn

import (
	"encoding/json"
	"fmt"
	"strconv"

	"fermion/backend_core/controllers/inventory_orders/grn"
	"fermion/backend_core/controllers/inventory_orders/inventory_adjustments"
	"fermion/backend_core/controllers/orders/scrap_orders"
	"fermion/backend_core/internal/model/inventory_orders"
	"fermion/backend_core/internal/model/pagination"
	inventory_orders_repo "fermion/backend_core/internal/repository/inventory_orders"
	access_checker "fermion/backend_core/pkg/util/access"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type Service interface {
	CreateAsn(data *AsnRequest, token_id string, access_template_id string) (uint, error)
	BulkCreateAsn(data *[]AsnRequest, token_id string, access_template_id string) error
	UpdateAsn(data uint, userID *AsnRequest, token_id string, access_template_id string) error
	GetAsn(id uint, token_id string, access_template_id string) (interface{}, error)
	GetAllAsn(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_orders.ASN, error)
	DeleteAsn(id uint, user_id uint, token_id string, access_template_id string) error
	DeleteAsnLines(query interface{}, token_id string, access_template_id string) error

	SendMailAsn(q *SendMailAsn) error
	GetAsnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error)
}

type service struct {
	asnRepository      inventory_orders_repo.Asn
	grnService         grn.Service
	inventoryService   inventory_adjustments.Service
	scrapOrdersService scrap_orders.Service
}

func NewService() *service {
	AsnRepository := inventory_orders_repo.NewAsn()
	return &service{AsnRepository,
		grn.NewService(),
		inventory_adjustments.NewService(),
		scrap_orders.NewService()}
}

func (s *service) CreateAsn(data *AsnRequest, token_id string, access_template_id string) (uint, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "ASN", *token_user_id)
	if !access_module_flag {
		return 0, fmt.Errorf("you dont have access for create asn at view level")
	}
	if data_access == nil {
		return 0, fmt.Errorf("you dont have access for create asn at data level")
	}
	var AsnData inventory_orders.ASN
	dto, err := json.Marshal(*data)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &AsnData)
	if err != nil {
		return 0, res.BuildError(res.ErrUnprocessableEntity, err)
	}

	if data.AutoCreateAsnNumber {
		AsnData.AsnNumber = helpers.GenerateSequence("ASN", token_id, "asns")
	}
	if data.AutoGenerateReferenceNumber {
		AsnData.ReferenceNumber = helpers.GenerateSequence("REF", token_id, "asns")
	}

	defaultStatus, err := helpers.GetLookupcodeId("GRN_STATUS", "DRAFT")
	if err != nil {
		return 0, err
	}
	data.StatusID = defaultStatus
	id, err := s.asnRepository.CreateAsn(&AsnData, token_id)
	return id, err
}

func (s *service) BulkCreateAsn(data *[]AsnRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "ASN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create asn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create asn at data level")
	}
	bulk_data := []map[string]interface{}{}
	for index, value := range *data {
		v := map[string]interface{}{}
		count := helpers.GetCount("SELECT COUNT(*) FROM asns") + 1 + index
		if value.AutoCreateAsnNumber {
			value.AsnNumber = "ASN/" + token_id + "/000" + strconv.Itoa(count)
		}
		if value.AutoGenerateReferenceNumber {
			value.ReferenceNumber = "REF/" + token_id + "/000" + strconv.Itoa(count)
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

	err = s.asnRepository.BulkCreateAsn(&AsnData, token_id)
	return err
}

func (s *service) UpdateAsn(id uint, data *AsnRequest, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "ASN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update asn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update asn at data level")
	}
	var AsnData inventory_orders.ASN
	dto, err := json.Marshal(*data)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	err = json.Unmarshal(dto, &AsnData)
	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	old_data, er := s.asnRepository.GetAsn(id)
	if er != nil {
		return er
	}
	old_status := old_data.StatusID
	new_status := AsnData.StatusID
	if new_status != old_status && new_status != 0 {
		result, _ := helpers.UpdateStatusHistory(old_data.StatusHistory, data.StatusID)
		AsnData.StatusHistory = result
	}

	err = s.asnRepository.UpdateAsn(id, &AsnData)
	for _, order_line := range AsnData.AsnOrderLines {
		query := map[string]interface{}{
			"asn_id":     uint(id),
			"product_id": uint(order_line.ProductID),
		}
		count, er := s.asnRepository.UpdateAsnLines(query, order_line)
		if er != nil {
			return er
		} else if count == 0 {
			order_line.Asn_id = id
			e := s.asnRepository.CreateAsnLines(order_line)
			if e != nil {
				return e
			}
		}
	}
	return err
}

func (s *service) GetAsn(id uint, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "ASN", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view asn at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view asn at data level")
	}
	result, er := s.asnRepository.GetAsn(id)
	if er != nil {
		return result, er
	}
	query := map[string]interface{}{
		"asn_id": id,
	}
	result_order_lines, err := s.asnRepository.GetAsnLines(query)
	result.AsnOrderLines = result_order_lines
	if err != nil {
		return result, err
	}

	return result, nil
}

func (s *service) GetAllAsn(p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) ([]inventory_orders.ASN, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "ASN", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list asn at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list asn at data level")
	}
	result, err := s.asnRepository.GetAllAsn(p)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteAsn(id uint, user_id uint, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "ASN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete asn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete asn at data level")
	}
	_, er := s.asnRepository.GetAsn(id)
	if er != nil {
		return er
	}
	err := s.asnRepository.DeleteAsn(id, user_id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{"asn_id": id}
	err1 := s.asnRepository.DeleteAsnLines(query)
	return err1
}

func (s *service) DeleteAsnLines(query interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "ASN", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete asn at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete asn at data level")
	}
	data, err := s.asnRepository.GetAsnLines(query)
	if err != nil {
		return err
	}
	if len(data) <= 0 {
		return err
	}
	err = s.asnRepository.DeleteAsnLines(query)
	return err
}

func (s *service) SendMailAsn(q *SendMailAsn) error {

	id, _ := strconv.Atoi(q.ID)
	result, er := s.asnRepository.GetAsn(uint(id))
	if er != nil {
		return er
	}
	err := helpers.SendEmail("ASN", q.ReceiverEmail, "pkg/util/response/static/inventory_orders/asn_template.html", result)

	return err
}
func (s *service) GetAsnTab(id, tab string, page *pagination.Paginatevalue, access_template_id string, token_id string) (interface{}, error) {

	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "ASN", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list asn at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list asn at data level")
	}

	asnId, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	if tab == "grn" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("GRN_SOURCE_DOCUMENT_TYPES", "ASN")

		grnPage := page
		grnPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, asnId)
		data, err := s.grnService.GetAllGRN(grnPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	if tab == "inventory_adjustments" {

		sourceDocumentId, _ := helpers.GetLookupcodeId("INVENTORY_ADJUSTMENTS_SOURCE_DOCUMENT_TYPES", "ASN")

		inventoryPage := page
		inventoryPage.Filters = fmt.Sprintf("[[\"source_document_type_id\",\"=\",%v],[\"source_documents\",\"@>\",\"{\\\"id\\\":%v}\"]]", sourceDocumentId, asnId)
		data, err := s.inventoryService.GetAllInvAdj(inventoryPage, token_id, access_template_id, "LIST")
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	return nil, nil
}
