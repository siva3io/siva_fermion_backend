package contacts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	core_model "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"
	core_user_repo "fermion/backend_core/internal/repository"
	access_template "fermion/backend_core/internal/repository/access"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
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
	CreateOrUpdateSubUser(data PartnerRequestDTO) error
	CreateContact(data *mdm.Partner, token_id string, access_template_id string) error
	UpdateContact(query map[string]interface{}, data *mdm.Partner, token_id string, access_template_id string) error
	DeleteContact(query map[string]interface{}, token_id string, access_template_id string) error
	GetContact(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error)
	GetContactList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error)
	SearchContact(query string, token_id string, access_template_id string) ([]IdNameDTO, error)
}

type service struct {
	contactsRepository       mdm_repo.Contacts
	coreuserRepository       core_user_repo.User
	coreRepository           core_user_repo.Core
	accessTemplateRepository access_template.Template
}

func NewService() *service {
	contactsRepository := mdm_repo.NewContacts()
	coreuserRepository := core_user_repo.NewUser()
	coreRepository := core_user_repo.NewCore()
	accessTemplateRepository := access_template.NewTemplate()
	return &service{contactsRepository, coreuserRepository, coreRepository, accessTemplateRepository}
}

func (s *service) CreateContact(data *mdm.Partner, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "CREATE", "CONTACTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for create contacts at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for create contacts at data level")
	}
	data.PrimaryEmail = strings.ToLower(data.PrimaryEmail)
	email_query := map[string]interface{}{
		"primary_email": data.PrimaryEmail,
	}
	_, err := s.contactsRepository.FindOneContact(email_query)
	if err == nil {
		return res.BuildError(res.ErrDuplicate, errors.New("oops! Email already Exists"))
	} else {
		if data.ParentId != nil && *data.ParentId == 0 {
			data.ParentId = nil
		}
		err := s.contactsRepository.SaveContact(data)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *service) UpdateContact(query map[string]interface{}, data *mdm.Partner, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "UPDATE", "CONTACTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for update contacts at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for update contacts at data level")
	}
	_, err := s.contactsRepository.FindOneContact(query)
	if err != nil {
		return err
	}
	err = s.contactsRepository.UpdateContact(query, data)
	if err != nil {
		return res.BuildError(res.ErrDataNotFound, err)
	}
	return nil
}
func (s *service) DeleteContact(query map[string]interface{}, token_id string, access_template_id string) error {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "DELETE", "CONTACTS", *token_user_id)
	if !access_module_flag {
		return fmt.Errorf("you dont have access for delete contacts at view level")
	}
	if data_access == nil {
		return fmt.Errorf("you dont have access for delete contacts at data level")
	}
	q := map[string]interface{}{
		"id": query["id"].(int),
	}
	_, er := s.contactsRepository.FindOneContact(q)
	if er != nil {
		return er
	}
	err := s.contactsRepository.DeleteContact(query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) GetContact(query map[string]interface{}, token_id string, access_template_id string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "READ", "CONTACTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for view contacts at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for view contacts at data level")
	}
	result, err := s.contactsRepository.FindOneContact(query)
	if err != nil {
		return result, err
	}

	var response PartnerResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &response)
	if err != nil {
		return nil, err
	}

	for _, address := range response.AddressDetails {
		if address.MarkDefaultAddress {
			response.DefaultAddress = address
			break
		}
	}

	return response, nil
}
func (s *service) GetContactList(query interface{}, p *pagination.Paginatevalue, token_id string, access_template_id string, access_action string) (interface{}, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, access_action, "CONTACTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list contacts at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list contacts at data level")
	}
	result, err := s.contactsRepository.FindAllContact(query, p)
	if err != nil {
		return result, err
	}
	var response []PartnerResponseDTO
	marshaldata, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaldata, &response)
	if err != nil {
		return nil, err
	}

	for index, ContactDetails := range response {
		for _, address := range ContactDetails.AddressDetails {
			if address.MarkDefaultAddress {
				response[index].DefaultAddress = address
				break
			}
		}
	}

	return response, nil
}
func (s *service) SearchContact(query string, token_id string, access_template_id string) ([]IdNameDTO, error) {
	token_user_id := helpers.ConvertStringToUint(token_id)
	access_module_flag, data_access := access_checker.ValidateUserAccess(access_template_id, "LIST", "CONTACTS", *token_user_id)
	if !access_module_flag {
		return nil, fmt.Errorf("you dont have access for list contacts at view level")
	}
	if data_access == nil {
		return nil, fmt.Errorf("you dont have access for list contacts at data level")
	}
	result, err := s.contactsRepository.SearchContact(query)
	var contactdetails []IdNameDTO
	if err != nil {
		return nil, err
	}
	for _, v := range result {
		var data IdNameDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		contactdetails = append(contactdetails, data)
	}
	return contactdetails, nil
}

func (s *service) CreateOrUpdateSubUser(data PartnerRequestDTO) error {
	var newUser core_model.CoreUsers
	newUser.Name = data.FirstName + " " + data.LastName
	newUser.FirstName = data.FirstName
	newUser.LastName = data.LastName
	newUser.Email = data.PrimaryEmail
	newUser.MobileNumber = data.PrimaryPhone
	newUser.Profile = data.ImageOptions
	newUser.CreatedByID = data.CreatedByID
	newUser.ExternalDetails = data.ExternalDetails

	query := map[string]interface{}{"id": *data.CreatedByID}
	SuperUserResponse, _ := s.coreuserRepository.FindUser(query)
	if SuperUserResponse.ID != 0 {
		newUser.CompanyId = SuperUserResponse.CompanyId
	}

	// fmt.Println("======>>>>new user :--", *newUser.CompanyId, *newUser.CreatedByID)

	query = map[string]interface{}{
		"id": *SuperUserResponse.CompanyId,
	}
	companyResponse, _ := s.coreRepository.FindCompany(query)
	if companyResponse.Name != "" {
		var AccessTemplateId uint
		companyDefaultSettings := map[string]interface{}{}
		AccessTemplateDetails := map[string]interface{}{}
		var ok bool
		json.Unmarshal(companyResponse.CompanyDefaults, &companyDefaultSettings)
		if companyDefaultSettings["access_template_details"] != nil {
			AccessTemplateDetails, ok = companyDefaultSettings["access_template_details"].(map[string]interface{})
			if !ok {
				return errors.New("access_template_details is not a map[string]interface{}")
			}
		}

		if AccessTemplateDetails == nil || AccessTemplateDetails["default_user_template_id"] == nil {
			defaultUserTemplate := new(core_model.AccessTemplate)
			defaultUserTemplate.Name = "DEFAULT USER TEMPLATE"
			defaultUserTemplate.CompanyId = companyResponse.ID
			// defaultUserTemplate.CreatedByID = &companyResponse.CreatedByID

			template, err := s.accessTemplateRepository.CreateTemplate(defaultUserTemplate)
			if err != nil {
				return err
			}
			AccessTemplateId = template
			AccessTemplateDetails["default_user_template_id"] = template
			companyDefaultSettings["access_template_details"] = AccessTemplateDetails

			updateQuery := map[string]interface{}{
				"name": companyResponse.Name,
			}
			jsondata, err := json.Marshal(companyDefaultSettings)
			if err != nil {
				return err
			}
			updateCompany := new(core_model.Company)
			updateCompany.CompanyDefaults = jsondata
			s.coreRepository.UpdateCompany(updateQuery, updateCompany)
		} else {
			fmt.Println("=========>>>> it takes already exists one")
			AccessTemplateId = uint(AccessTemplateDetails["default_user_template_id"].(float64))
		}

		newUser.AccessTemplateId = append(newUser.AccessTemplateId, &core_model.AccessTemplate{ID: AccessTemplateId})

		query = map[string]interface{}{"login": data.PrimaryEmail}
		resp, _ := s.coreuserRepository.FindLogin(query)
		if resp.ID != 0 {
			err := s.coreuserRepository.UpdateUser(query, &newUser)
			if err != nil {
				return err
			}
			fmt.Println("--------------------->>> Sub User Updated Successfully <<<-----------------------------")

		} else {
			err := s.coreuserRepository.Save(&newUser)
			if err != nil {
				return err
			}
			fmt.Println("--------------------->>> Sub User Created Successfully <<<-----------------------------")

		}
	}
	return nil
}
