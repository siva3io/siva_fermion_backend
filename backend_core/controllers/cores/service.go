package cores

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	cache "fermion/backend_core/controllers/cache"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/repository"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"
)

const (
	TrialDuration model_core.Hours = 720
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
	GetLookupTypes(interface{}, *pagination.Paginatevalue) ([]LookupTypesDTO, error)
	GetLookupCodes(interface{}, *pagination.Paginatevalue) ([]LookupCodesDTO, error)
	GetCountries(interface{}, *pagination.Paginatevalue) ([]CountryDTO, error)
	GetStates(interface{}, *pagination.Paginatevalue) ([]StateDTO, error)
	GetCurrencies(interface{}, *pagination.Paginatevalue) ([]CurrencyDTO, error)
	SearchLookupCodes(query string) ([]LookupCodesDTO, error)
	CreateLookupType(data *model_core.Lookuptype) error
	CreateLookupCodes(data *[]model_core.Lookupcode) error
	UpdateLookupType(data *model_core.Lookuptype, query interface{}) error
	UpdateLookupCode(data *model_core.Lookupcode, query interface{}) error
	DeleteLookupType(id uint) error
	DeleteLookupCode(id uint) error
	GetLookupTypesList(query interface{}, p *pagination.Paginatevalue) ([]LookupTypesDTO, error)
	GetLookupCodesList(query interface{}, p *pagination.Paginatevalue) ([]LookupCodeDTO, error)
	//apps
	CreateAppStore(data *model_core.AppStore) error
	UpdateAppStore(data *model_core.AppStore, query map[string]interface{}) error
	ListStoreApps(query interface{}, page *pagination.Paginatevalue) ([]model_core.AppStore, error)
	ListStoreAppsInstalled(page *pagination.Paginatevalue) ([]model_core.AppStore, error)
	ListInstalledApps(page *pagination.Paginatevalue) ([]model_core.InstalledApps, error)
	GetApp(query map[string]interface{}) (model_core.AppStore, error)
	SearchApps(data SearchAppsDTO) (model_core.AppStore, error)
	CheckState(code string) bool
	//install app
	InstallApp(app_code string, user_token string) error
	UninstallApp(app_code string) error

	AllChannelLookupCodes(interface{}, *pagination.Paginatevalue) ([]model_core.ChannelLookupCodes, error)
	//meta data
	ListMetaData() ([]model_core.IRModel, error)
	ViewMetaData(model_name string) ([]model_core.IRModelFields, error)
	UpdateCompanyPreferences(query map[string]interface{}, company Company, host string, schema string) (*model_core.Company, error)
	GetCompanyPreferences(query map[string]interface{}) (*Company, error)
	ListCompanyUsers(query map[string]interface{}, page *pagination.Paginatevalue) ([]UserResponseDTO, error)
	UpdateCompany(query map[string]interface{}, company model_core.Company) error
}

type service struct {
	coreRepository repository.Core
	userRepository repository.User
}

func NewService() *service {
	CoreRepository := repository.NewCore()
	userRepository := repository.NewUser()
	return &service{CoreRepository, userRepository}
}

func (s *service) GetLookupTypes(query interface{}, p *pagination.Paginatevalue) ([]LookupTypesDTO, error) {
	var data []LookupTypesDTO
	lookuptype_query := map[string]interface{}{
		"is_enabled": true,
	}
	results, er := s.coreRepository.GetLookupTypes(lookuptype_query, p)
	json_data, err := json.Marshal(results)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(json_data, &data)
	if err != nil {
		return data, err
	}
	return data, er
}

func (s *service) GetLookupTypesList(query interface{}, p *pagination.Paginatevalue) ([]LookupTypesDTO, error) {
	var data []LookupTypesDTO
	result, _ := s.coreRepository.FindAllLookupTypes(query, p)
	json_data, err := json.Marshal(result)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(json_data, &data)
	if err != nil {
		return data, err
	}
	for i := 0; i < len(result); i++ {
		data[i].NumberOfLookUpCodes = uint(len(result[i].Lookupcodes))
	}
	return data, err
}

func (s *service) GetLookupCodes(query interface{}, p *pagination.Paginatevalue) ([]LookupCodesDTO, error) {
	var lookup_codes_dto []LookupCodesDTO
	p.Per_page = 100
	query.(map[string]interface{})["is_enabled"] = true
	nilObj := new(pagination.Paginatevalue)
	lookuptype, err := s.coreRepository.GetLookupTypes(query, nilObj)
	if err != nil {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}
	if len(lookuptype) == 0 {
		return nil, res.BuildError(res.ErrDataNotFound, errors.New("data not found"))
	}
	lookupcode_query := map[string]interface{}{
		"lookup_type_id": lookuptype[0].ID,
		"is_enabled":     true,
	}
	lookupcodes, err := s.coreRepository.GetLookupCodes(lookupcode_query, p)
	if err != nil {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}
	jsondata, _ := json.Marshal(lookupcodes)
	err = json.Unmarshal(jsondata, &lookup_codes_dto)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return lookup_codes_dto, nil
}

func (s *service) GetLookupCodesList(query interface{}, p *pagination.Paginatevalue) ([]LookupCodeDTO, error) {
	var lookup_codes_dto []LookupCodeDTO
	query.(map[string]interface{})["is_enabled"] = true
	lookupcodes, err := s.coreRepository.FindAllLookupCodes(query, p)
	if err != nil {
		return nil, res.BuildError(res.ErrDataNotFound, err)
	}
	jsondata, _ := json.Marshal(lookupcodes)
	err = json.Unmarshal(jsondata, &lookup_codes_dto)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	for index := 0; index < len(lookup_codes_dto); index++ {
		lookup_codes_dto[index].LookupType = lookupcodes[index].LookupType.LookupType
	}
	return lookup_codes_dto, nil
}

func (s *service) GetCountries(query interface{}, p *pagination.Paginatevalue) ([]CountryDTO, error) {
	p.Per_page = 250
	results, er := s.coreRepository.GetCountries(query, p)
	var data []CountryDTO
	json_data, er := json.Marshal(results)
	if er != nil {
		return data, er
	}
	er = json.Unmarshal(json_data, &data)
	if er != nil {
		return data, er
	}
	return data, er
}
func (s *service) GetStates(query interface{}, p *pagination.Paginatevalue) ([]StateDTO, error) {
	p.Per_page = 150
	results, er := s.coreRepository.GetStates(query, p)
	var data []StateDTO
	json_data, er := json.Marshal(results)
	if er != nil {
		return data, er
	}
	er = json.Unmarshal(json_data, &data)
	if er != nil {
		return data, er
	}
	return data, er
}
func (s *service) GetCurrencies(query interface{}, p *pagination.Paginatevalue) ([]CurrencyDTO, error) {
	p.Per_page = 250
	results, er := s.coreRepository.GetCurrencies(query, p)
	var data []CurrencyDTO
	json_data, er := json.Marshal(results)
	if er != nil {
		return data, er
	}
	er = json.Unmarshal(json_data, &data)
	if er != nil {
		return data, er
	}
	return data, er
}
func (s *service) SearchLookupCodes(query string) ([]LookupCodesDTO, error) {

	result, err := s.coreRepository.SearchLookupCodes(query)
	var lookupcodes []LookupCodesDTO
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	for _, v := range result {
		var data LookupCodesDTO
		value, _ := json.Marshal(v)
		err := json.Unmarshal(value, &data)
		if err != nil {
			return nil, err
		}
		lookupcodes = append(lookupcodes, data)
	}
	return lookupcodes, nil
}
func (s *service) CreateLookupType(data *model_core.Lookuptype) error {
	er := s.coreRepository.CreateTypes(data)
	return er
}
func (s *service) CreateLookupCodes(data *[]model_core.Lookupcode) error {
	er := s.coreRepository.CreateCodes(data)
	return er
}
func (s *service) UpdateLookupType(data *model_core.Lookuptype, query interface{}) error {
	er := s.coreRepository.UpdateType(data, query)
	return er
}
func (s *service) UpdateLookupCode(data *model_core.Lookupcode, query interface{}) error {
	er := s.coreRepository.UpdateCode(data, query)
	return er
}
func (s *service) DeleteLookupType(id uint) error {
	err := s.coreRepository.DeleteType(id)
	if err != nil {
		return err
	}
	query := map[string]interface{}{
		"lookup_type_id": id,
	}
	err1 := s.coreRepository.DeleteCode(query)
	if err1 != nil {
		return err
	}
	return nil
}
func (s *service) DeleteLookupCode(id uint) error {
	query := map[string]interface{}{
		"id": id,
	}
	err := s.coreRepository.DeleteCode(query)
	return err
}

// apps

func (s *service) CreateAppStore(data *model_core.AppStore) error {
	err := s.coreRepository.CreateAppStore(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateAppStore(data *model_core.AppStore, query map[string]interface{}) error {
	err := s.coreRepository.UpdateAppStore(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) ListStoreApps(query interface{}, page *pagination.Paginatevalue) ([]model_core.AppStore, error) {
	results, er := s.coreRepository.ListStoreApps(query, page)
	return results, er
}
func (s *service) ListStoreAppsInstalled(page *pagination.Paginatevalue) ([]model_core.AppStore, error) {
	results, er := s.coreRepository.ListStoreAppsInstalled(page)
	return results, er
}
func (s *service) ListInstalledApps(page *pagination.Paginatevalue) ([]model_core.InstalledApps, error) {
	results, er := s.coreRepository.ListInstalledApps(page)
	return results, er
}

func (s *service) GetApp(query map[string]interface{}) (model_core.AppStore, error) {
	results, er := s.coreRepository.GetApp(query)
	return results, er
}
func (s *service) SearchApps(data SearchAppsDTO) (model_core.AppStore, error) {
	var query = make(map[string]interface{}, 0)
	if data.Name != "" {
		query["name"] = data.Name
	}
	if data.Code != "" {
		query["code"] = data.Code
	}
	result, err := s.coreRepository.SearchApps(query)
	return result, err
}
func (s *service) AllChannelLookupCodes(query interface{}, p *pagination.Paginatevalue) ([]model_core.ChannelLookupCodes, error) {
	results, er := s.coreRepository.GetChannelLookupCodes(query, p)
	return results, er
}
func (s *service) ListMetaData() ([]model_core.IRModel, error) {
	response, err := s.coreRepository.ListMetaData()
	return response, err
}
func (s *service) ViewMetaData(model_name string) ([]model_core.IRModelFields, error) {
	response, err := s.coreRepository.ViewMetaData(model_name)
	return response, err
}

func (s *service) InstallApp(app_code string, user_id string) error {
	query := map[string]interface{}{
		"code": app_code,
	}
	app_query := map[string]interface{}{
		"id": helpers.ConvertStringToUint(user_id),
	}
	get_response, err := s.coreRepository.GetApp(query)
	if err != nil {
		return err
	}
	check_code := s.coreRepository.CheckState(get_response.Code)
	if !check_code {
		var install_app_data = new(model_core.InstalledApps)
		value, _ := json.Marshal(get_response)
		json.Unmarshal(value, &install_app_data)
		install_app_data.Type = "Installed"
		install_app_data.CurrentVersion = "1.0.0"
		install_app_data.ID = 0
		if get_response.CategoryName == "core_app" {
			install_app_data.IsCore = true
		}
		if get_response.CategoryName != "core_app" {
			install_app_data.IsCore = false
		}

		user, _ := s.userRepository.FindUser(app_query)
		install_app_data.AccessToken, _ = user.GenerateToken(2, TrialDuration)
		install_app_code := install_app_data.Code
		install_app_code = strings.ReplaceAll(install_app_code, " ", "")
		cache.SetCacheVariable(install_app_code+"AUTHORIZATION", install_app_data.AccessToken)
		create_err := s.coreRepository.InstallApp(install_app_data)
		return create_err
	}
	return nil
	// var install_app_data model.InstalledApps
}
func (s *service) CheckState(code string) bool {
	check_code := s.coreRepository.CheckState(code)
	return check_code
}

func (s *service) UninstallApp(app_code string) error {
	query := map[string]interface{}{
		"code": app_code,
	}
	check_code := s.coreRepository.UninstallApp(query)
	return check_code
}

func (s *service) UpdateCompanyPreferences(query map[string]interface{}, company Company, host string, schema string) (*model_core.Company, error) {
	// user_token := query["user_token"].(string)
	delete(query, "user_token")
	company_details, _ := s.GetCompanyPreferences(query)
	if company.OrganizationDetails == nil {
		company.OrganizationDetails = company_details.OrganizationDetails
	}

	if company.KYCDocuments == nil {
		company.KYCDocuments = company_details.KYCDocuments
	} else {

		var country_updated []KYCDocuments
		flag := false

		for _, country_to_be_updated := range company.KYCDocuments {
			if company_details.KYCDocuments != nil {
				for _, country := range company_details.KYCDocuments {
					if country_to_be_updated.CountryId == country.CountryId {
						for k, v := range country_to_be_updated.Documents {
							country.Documents[k] = v
							flag = true
						}
					}
					country_updated = append(country_updated, country)
				}

			}
		}
		if !flag {
			company.KYCDocuments = append(company.KYCDocuments, company_details.KYCDocuments...)
		} else {
			company.KYCDocuments = country_updated

		}
		count := 1
		for _, value := range company.KYCDocuments {
			// country_id := value.CountryId
			documents := value.Documents
			for key, value := range documents {
				var imageLinkArray = make([]map[string]interface{}, 0)
				imageLinkArray = append(imageLinkArray, map[string]interface{}{"data": value})
				//------TODO : condition to check company file preferences and upload_file accordingly ---
				// response, _ := ipaas_utils.UploadFile(imageLinkArray, key, strconv.Itoa(int(*company.UpdatedByID)), *company.UpdatedByID, schema, host, ipaas_utils.AWS, "Core", "Company", ".png", ipaas_utils.BASE64, user_token)
				// documents[key] = response
				documents[key] = imageLinkArray
			}
			count++
		}
	}

	if company.FilePreferenceID == nil {
		company.FilePreferenceID = company_details.FilePreferenceID
	}

	if company.InvoiceGenerationID == nil {
		company.InvoiceGenerationID = company_details.InvoiceGenerationID
	}

	if company.BusinessTypeID == nil {
		company.BusinessTypeID = company_details.BusinessTypeID
	}

	var company_model model_core.Company
	dto, _ := json.Marshal(company)
	err := json.Unmarshal(dto, &company_model)

	result, err := s.userRepository.UpdateCompany(query, company_model)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return &result, err
}

func (s *service) GetCompanyPreferences(query map[string]interface{}) (*Company, error) {
	var company Company
	result, err := s.userRepository.GetCompany(query)
	dto, _ := json.Marshal(result)
	_ = json.Unmarshal(dto, &company)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return &company, err
}

func (s *service) ListCompanyUsers(query map[string]interface{}, page *pagination.Paginatevalue) ([]UserResponseDTO, error) {
	var users []UserResponseDTO
	result, err := s.userRepository.FindAllCompanyUsers(query, page)
	dto, _ := json.Marshal(result)
	_ = json.Unmarshal(dto, &users)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return users, err
}

func (s *service) UpdateCompany(query map[string]interface{}, company model_core.Company) error {
	result, _ := s.userRepository.GetCompany(query)
	company_address := make([]map[string]interface{}, 0)
	result_address := make([]map[string]interface{}, 0)
	_ = json.Unmarshal(company.Addresses, &company_address)
	_ = json.Unmarshal(result.Addresses, &result_address)

	result_address = append(result_address, company_address...)

	bt, _ := json.Marshal(result_address)
	json.Unmarshal(bt, &company.Addresses)

	_, err := s.userRepository.UpdateCompany(query, company)
	if err != nil {
		fmt.Println("error", err)
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}

	return err
}
