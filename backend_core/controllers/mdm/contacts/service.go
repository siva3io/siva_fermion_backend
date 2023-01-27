package contacts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/internal/model/core"
	core_model "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"
	core_user_repo "fermion/backend_core/internal/repository"
	access_template "fermion/backend_core/internal/repository/access"
	inventory_orders_repo "fermion/backend_core/internal/repository/inventory_orders"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	orders_repo "fermion/backend_core/internal/repository/orders"
	returns_repo "fermion/backend_core/internal/repository/returns"
	access_checker "fermion/backend_core/pkg/util/access"

	// ds "fermion/backend_core/pkg/util/dynamic_struct"
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
	Create(metaData core.MetaData, data *mdm.Partner) error
	Update(metaData core.MetaData, data *mdm.Partner) error
	Delete(metaData core.MetaData) error
	FindOne(metaData core.MetaData) (interface{}, error)
	FindAll(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	Upsert(metaData core.MetaData, data []interface{}) (interface{}, error)

	GetRelatedTabs(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error)
	CreateOrUpdateSubUser(metaData core.MetaData, data PartnerRequestDTO) (uint, error)
}

type service struct {
	ContactsRepository       mdm_repo.Contacts
	CoreuserRepository       core_user_repo.User
	CoreRepository           core_user_repo.Core
	AccessTemplateRepository access_template.Template
	SalesorderRepository     orders_repo.Sales
	SalesReturnRepository    returns_repo.SalesReturn
	DeliveryOrdersRepository orders_repo.DeliveryOrders
	PurchaseOrderRepository  orders_repo.Purchase
	AsnRepository            inventory_orders_repo.Asn
	GrnRepository            inventory_orders_repo.GRN
	PurchaseReturnRepository returns_repo.PurchaseReturn
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		mdm_repo.NewContacts(),
		core_user_repo.NewUser(),
		core_user_repo.NewCore(),
		access_template.NewTemplate(),
		orders_repo.NewSalesOrder(),
		returns_repo.NewSalesReturn(),
		orders_repo.NewDo(),
		orders_repo.NewPurchaseOrder(),
		inventory_orders_repo.NewAsn(),
		inventory_orders_repo.NewGRN(),
		returns_repo.NewPurchaseReturn()}

	return newServiceObj
}

// ============================ CRUD services =================================================================
func (s *service) Create(metaData core.MetaData, data *mdm.Partner) error {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "CREATE", "CONTACTS", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for create contacts at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for create contacts at data level")
	}

	data.PrimaryEmail = strings.ToLower(data.PrimaryEmail)
	data.CreatedByID = &metaData.TokenUserId
	data.CompanyId = metaData.CompanyId

	searchQuery := map[string]interface{}{
		"primary_email": data.PrimaryEmail,
		"company_id":    data.CompanyId,
	}
	_, err := s.ContactsRepository.FindOne(searchQuery)
	if err == nil {
		return errors.New("oops! Email already Exists")
	}
	if data.ParentId != nil && *data.ParentId == 0 {
		data.ParentId = nil
	}
	err = s.ContactsRepository.Create(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Update(metaData core.MetaData, data *mdm.Partner) error {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "UPDATE", "CONTACTS", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for update contacts at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for update contacts at data level")
	}

	data.UpdatedByID = &metaData.TokenUserId
	err := s.ContactsRepository.Update(metaData.Query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) Delete(metaData core.MetaData) error {
	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "DELETE", "CONTACTS", metaData.TokenUserId)
	if !accessModuleFlag {
		return fmt.Errorf("you dont have access for delete contacts at view level")
	}
	if dataAccess == nil {
		return fmt.Errorf("you dont have access for delete contacts at data level")
	}

	err := s.ContactsRepository.Delete(metaData.Query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) FindOne(metaData core.MetaData) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "CONTACTS", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for view contacts at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for view contacts at data level")
	}

	contactDetails, err := s.ContactsRepository.FindOne(metaData.Query)
	if err != nil {
		return contactDetails, err
	}

	contactDetails.Properties = helpers.GetLookupCodesFromArrOfObjectIdsToArrOfObjects(contactDetails.Properties)

	return contactDetails, nil
}
func (s *service) FindAll(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, metaData.ModuleAccessAction, "CONTACTS", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list contacts at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list contacts at data level")
	}

	arrayOfContactDetails, err := s.ContactsRepository.FindAll(metaData.Query, p)
	if err != nil {
		return arrayOfContactDetails, err
	}

	return arrayOfContactDetails, nil
}
func (s *service) Upsert(metaData core.MetaData, data []interface{}) (interface{}, error) {

	var success []interface{}
	var failures []interface{}

	for index, payload := range data {
		var requestdto PartnerRequestDTO
		contact := new(mdm.Partner)
		payloadResponse := map[string]interface{}{}

		helpers.JsonMarshaller(payload, &requestdto)
		helpers.JsonMarshaller(payload, contact)

		searchQuery := map[string]interface{}{
			"primary_email": requestdto.PrimaryEmail,
			"company_id":    metaData.CompanyId,
		}
		contactResponse, _ := s.ContactsRepository.FindOne(searchQuery)
		if contactResponse.PrimaryEmail != "" {
			contact.UpdatedByID = &metaData.TokenUserId
			err := s.ContactsRepository.Update(searchQuery, contact)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			}
			payloadResponse["msg"] = "updated"
		} else {
			contact.CompanyId = metaData.CompanyId
			contact.CreatedByID = &metaData.TokenUserId
			err := s.ContactsRepository.Create(contact)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
			}
			payloadResponse["msg"] = "created"
		}
		contactResponse, _ = s.ContactsRepository.FindOne(searchQuery)
		//================= create or update sub_user =================
		if requestdto.IsAllowedLogin != nil && *requestdto.IsAllowedLogin {
			helpers.JsonMarshaller(contactResponse, &requestdto)
			userId, err := s.CreateOrUpdateSubUser(metaData, requestdto)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			contact.UserId = userId
			err = s.ContactsRepository.Update(searchQuery, contact)
			if err != nil {
				failures = append(failures, map[string]interface{}{"status": false, "serial_number": index + 1, "msg": err})
				continue
			}
			payloadResponse["sub_user"] = "upserted"
		}

		//================= create or update vendor ===================
		isVendorCreate := false
		var contactDetails PartnerResponseDTO
		helpers.JsonMarshaller(contactResponse, &contactDetails)
		for _, object := range contactDetails.Properties {
			LookupCode, _ := helpers.GetLookupcodeName(object.Id)
			if LookupCode == "VENDOR" {
				isVendorCreate = true
				break
			}
		}
		if isVendorCreate {
			vendor := new(mdm.Vendors)
			if contactDetails.FirstName != "" || contactDetails.LastName != "" {
				vendor.Name = contactDetails.FirstName + " " + contactDetails.LastName
			} else {
				if contactDetails.CompanyName != "" {
					vendor.Name = contactDetails.CompanyName
				}
			}
			vendor.ContactId = &contactDetails.Id
			vendor.PrimaryContactId = contactDetails.ParentId
			vendorPayload := []interface{}{}
			vendorPayload = append(vendorPayload, vendor)
			request_payload := map[string]interface{}{
				"meta_data": metaData,
				"data":      vendorPayload,
			}
			eda.Produce(eda.UPSERT_VENDOR, request_payload)
			payloadResponse["vendor"] = "upserted"
		}
		payloadResponse["status"] = true
		payloadResponse["serial_number"] = index + 1
		success = append(success, payloadResponse)
	}
	response := map[string]interface{}{
		"success":  success,
		"failures": failures,
	}
	return response, nil
}

// ============================= other services ================================================================
func (s *service) GetRelatedTabs(metaData core.MetaData, p *pagination.Paginatevalue) (interface{}, error) {

	accessTemplateId := strconv.FormatUint(uint64(metaData.AccessTemplateId), 10)
	contactId := metaData.AdditionalFields["contact_id"]

	accessModuleFlag, dataAccess := access_checker.ValidateUserAccess(accessTemplateId, "READ", "CONTACTS", metaData.TokenUserId)
	if !accessModuleFlag {
		return nil, fmt.Errorf("you dont have access for list contacts at view level")
	}
	if dataAccess == nil {
		return nil, fmt.Errorf("you dont have access for list contacts at data level")
	}

	tabName := metaData.AdditionalFields["tab_name"].(string)

	if tabName == "sales_orders" {
		salesOrderPage := p
		salesOrderPage.Filters = fmt.Sprintf("[[\"customer_billing_address\",\"@>\",\"{\\\"id\\\":%v}\"]]", contactId)
		relatedSalesOrders, err := s.SalesorderRepository.FindAll(nil, salesOrderPage)
		if err != nil {
			return nil, err
		}
		return relatedSalesOrders, nil
	}
	if tabName == "sales_returns" {
		salesReturnsPage := p
		salesReturnsPage.Filters = fmt.Sprintf("[[\"customer_billing_address\",\"@>\",\"{\\\"id\\\":%v}\"]]", contactId)
		relatedSalesReturns, err := s.SalesReturnRepository.FindAll(nil, salesReturnsPage)
		if err != nil {
			return nil, err
		}
		return relatedSalesReturns, nil
	}
	if tabName == "delivery_orders" {
		deliveryOrdersPage := p
		deliveryOrdersPage.Filters = fmt.Sprintf("[[\"billing_address_details\",\"@>\",\"{\\\"id\\\":%v}\"]]", contactId)
		relatedDeliveryOrders, err := s.DeliveryOrdersRepository.AllDeliveryOrders(nil, deliveryOrdersPage)
		if err != nil {
			return nil, err
		}
		return relatedDeliveryOrders, nil
	}
	if tabName == "purchase_orders" {
		purchaseOrdersPage := p
		purchaseOrdersPage.Filters = fmt.Sprintf("[[\"billing_address\",\"@>\",\"{\\\"id\\\":%v}\"]]", contactId)
		relatedPurchaseOrders, err := s.PurchaseOrderRepository.FindAll(nil, purchaseOrdersPage)
		if err != nil {
			return nil, err
		}
		return relatedPurchaseOrders, nil
	}
	if tabName == "asn" {
		asnPage := p
		asnPage.Filters = fmt.Sprintf("[[\"destination_location_details\",\"@>\",\"{\\\"id\\\":%v}\"]]", contactId)
		relatedAsns, err := s.AsnRepository.GetAllAsn(metaData.Query, asnPage)
		if err != nil {
			return nil, err
		}
		return relatedAsns, nil
	}
	if tabName == "purchase_returns" {
		purchaseReturnsPage := p
		purchaseReturnsPage.Filters = fmt.Sprintf("[[\"vendor_details\",\"@>\",\"{\\\"contact_id\\\":%v}\"]]", contactId)
		relatedPurchaseReturns, err := s.PurchaseReturnRepository.FindAll(nil, purchaseReturnsPage)
		if err != nil {
			return nil, err
		}
		return relatedPurchaseReturns, nil
	}
	return nil, nil
}

func (s *service) CreateOrUpdateSubUser(metaData core.MetaData, data PartnerRequestDTO) (uint, error) {

	var newUser core_model.CoreUsers
	newUser.Name = data.FirstName + " " + data.LastName
	newUser.FirstName = data.FirstName
	newUser.LastName = data.LastName
	newUser.Email = data.PrimaryEmail
	newUser.Username = data.PrimaryEmail
	newUser.MobileNumber = data.PrimaryPhone
	newUser.Profile = data.ImageOptions
	newUser.ExternalDetails = data.ExternalDetails
	newUser.RoleId = data.RoleId
	newUser.UserTypes = data.UserTypes
	newUser.AccessIds = data.AccessIds

	newUser.AccessTemplateId = data.AccessTemplateId
	newUser.CreatedByID = &metaData.TokenUserId
	newUser.CompanyId = &metaData.CompanyId

	if len(newUser.AccessTemplateId) == 0 {
		query := map[string]interface{}{
			"id": metaData.CompanyId,
		}
		companyResponse, _ := s.CoreRepository.FindCompany(query)
		if companyResponse.Name != "" {
			var AccessTemplateId uint
			companyDefaultSettings := map[string]interface{}{}
			json.Unmarshal(companyResponse.CompanyDefaults, &companyDefaultSettings)
			if companyDefaultSettings["access_template_details"] != nil {
				AccessTemplateDetails, ok := companyDefaultSettings["access_template_details"].(map[string]interface{})
				if !ok {
					return 0, errors.New("access_template_details not exists in company default settings")
				}
				if AccessTemplateDetails != nil {
					AccessTemplateId = uint(AccessTemplateDetails["default_user_template_id"].(float64))
				}
			}
			newUser.AccessTemplateId = append(newUser.AccessTemplateId, &core_model.AccessTemplate{ID: AccessTemplateId})
		}
	}

	query := map[string]interface{}{"email": data.PrimaryEmail}
	resp, _ := s.CoreuserRepository.FindUser(query)
	if resp.ID != 0 {
		err := s.CoreuserRepository.UpdateUser(query, &newUser)
		if err != nil {
			return resp.ID, err
		}
		fmt.Println("--------------------->>> Sub User Updated Successfully <<<-----------------------------")

	} else {
		err := s.CoreuserRepository.Save(&newUser)
		if err != nil {
			return 0, err
		}
		fmt.Println("--------------------->>> Sub User Created Successfully <<<-----------------------------")

	}
	return newUser.ID, nil
}
