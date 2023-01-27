package cores

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"strings"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/repository"
	omnichannel_repo "fermion/backend_core/internal/repository/omnichannel"
	AppInit "fermion/backend_core/pkg/pkg_init/app_init"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/google/uuid"
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
	ListInstalledApps(page *pagination.Paginatevalue, registered bool) ([]model_core.InstalledApps, error)
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

	//company
	GetCompanyPreferences(query map[string]interface{}) (*Company, error)
	ListCompanyUsers(query map[string]interface{}, page *pagination.Paginatevalue) ([]UserResponseDTO, error)
	UpdateCompany(query map[string]interface{}, company Company, host string, scheme string) error

	//channel status
	ListChannelStatus(page *pagination.Paginatevalue) ([]model_core.ChannelLookupCodes, error)
	ViewChannelStatus(id uint) ([]model_core.ChannelLookupCodes, error)
	CreateChannelStatus(data *model_core.ChannelLookupCodes) error
	UpdateChannelStatus(data *model_core.ChannelLookupCodes, query map[string]interface{}) error
	DeleteChannelStatus(id uint) error

	//ondc details
	UpdateOndcDetails(query map[string]interface{}, data map[string]interface{}) error
	OndcRegistration(query map[string]interface{}, company_details Company) (*model_core.OndcDetails, error)
	RequestSolution(data model_core.CustomSolution) error

	ListCompanies(page *pagination.Paginatevalue) ([]Company, error)

	UpdateInstalledApp(data *model_core.InstalledApps, query map[string]interface{}) error

	RenewSubscription(query map[string]interface{}, user_id string) error
}

type service struct {
	coreRepository             repository.Core
	userRepository             repository.User
	omnichannelFieldRepository omnichannel_repo.OmnichannelField
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{
		repository.NewCore(),
		repository.NewUser(),
		omnichannel_repo.NewOmnichannelFieldRepo(),
	}
	return newServiceObj
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
	p.Per_page = -1
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
	if er != nil {
		return nil, er
	}
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
	if er != nil {
		return nil, er
	}
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
	if er != nil {
		return nil, er
	}
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

func (s *service) UpdateInstalledApp(data *model_core.InstalledApps, query map[string]interface{}) error {
	err := s.coreRepository.UpdateInstalledApp(query, data)
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
func (s *service) ListInstalledApps(page *pagination.Paginatevalue, registered bool) ([]model_core.InstalledApps, error) {

	installedApps, err := s.coreRepository.ListInstalledApps(page)
	if err != nil {
		return nil, err
	}

	if !registered {
		return installedApps, nil
	}

	var result []model_core.InstalledApps
	for _, app := range installedApps {

		appData := map[string]interface{}{
			"id":           app.ID,
			"access_token": app.AccessToken,
		}
		appField, _ := helpers.GetAppAuthFields(appData, []string{"register_status"})

		if appField["register_status"] != nil && appField["register_status"].(map[string]interface{})["value"] != nil {
			registerStatus := appField["register_status"].(map[string]interface{})["value"].(bool)
			if registerStatus {
				result = append(result, app)
			}

		}
	}
	return result, nil
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
		// install_app_code := install_app_data.Code
		// install_app_code = strings.ReplaceAll(install_app_code, " ", "")
		// cache.SetCacheVariable(install_app_code+"AUTHORIZATION", install_app_data.AccessToken)
		create_err := s.coreRepository.InstallApp(install_app_data)
		AppInit.AppInit(*install_app_data)
		return create_err
	}
	return nil
	// var install_app_data model.InstalledApps
}

func (s *service) RenewSubscription(query map[string]interface{}, user_id string) error {
	app_query := map[string]interface{}{
		"id": user_id,
	}
	user, _ := s.userRepository.FindUser(app_query)
	token, err := user.GenerateToken(2, TrialDuration)
	if err != nil {
		return err
	}
	p := new(pagination.Paginatevalue)
	code := query["app_code"].(string)
	p.Filters = "[[\"code\",\"=\",\"" + code + "\"],[\"is_active\",\"=\",\"true\"]]"
	data, err := s.coreRepository.ListInstalledApps(p)
	if err != nil {
		return err
	}
	var app_data model_core.InstalledApps
	if len(data) > 0 {
		app_data = data[0]
	}
	app_data.AccessToken = token
	app_id := map[string]interface{}{
		"id": app_data.ID,
	}
	update := s.coreRepository.UpdateInstalledApp(app_id, &app_data)
	if err != nil {
		return err
	}
	return update

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

func (s *service) GetCompanyPreferences(query map[string]interface{}) (*Company, error) {
	var company Company
	result, err := s.userRepository.GetCompany(query)
	dto, _ := json.Marshal(result)
	_ = json.Unmarshal(dto, &company)
	if err != nil {
		return nil, err
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

func (s *service) UpdateCompany(query map[string]interface{}, company Company, host string, scheme string) error {
	var company_details model_core.Company
	bt, _ := json.Marshal(company)
	_ = json.Unmarshal(bt, &company_details)

	_, err := s.userRepository.UpdateCompany(query, company_details)

	if err != nil {
		return res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return nil
}

func (s *service) OndcRegistration(query map[string]interface{}, company_details Company) (*model_core.OndcDetails, error) {
	result, _ := s.userRepository.GetCompany(query)
	var OndcDetails model_core.OndcDetails

	if len(result.OndcDetails.SigningPrivateKey) <= 0 && len(result.OndcDetails.SubscriberId) <= 0 {
		replacer := strings.NewReplacer("-", "_", " ", "_")
		OndcDetails.SubscriberURL = fmt.Sprintf("%s_BPP", replacer.Replace(result.CompanyDetails.BusinessName))
		OndcDetails.SubscriberId = fmt.Sprintf("SIVA_ONDC_%s_BPP", replacer.Replace(result.CompanyDetails.BusinessName))
		OndcDetails.SigningPrivateKey, OndcDetails.SigningPublicKey = s.GeneratePublicPrivateKeys()
		OndcDetails.EncryptionPrivateKey, OndcDetails.EncryptionPublicKey = s.GeneratePublicPrivateKeys()
		OndcDetails.UniqueId = uuid.New().String()
		OndcDetails.Type = company_details.Type

		OndcDetailsId, err := s.coreRepository.CreateOndcDetails(&OndcDetails)
		if err != nil {
			fmt.Println(err)
		}
		result.OndcDetailsId = OndcDetailsId
	} else {
		fmt.Println("Already in ONDC")
		OndcDetails = result.OndcDetails
	}

	result.Type = company_details.Type

	s.userRepository.UpdateCompany(query, result)

	// Need to hit the Sandbox API for registering with ONDC network (NEED TO BE IMPLEMENTED)

	return &OndcDetails, nil
}

func (s *service) GeneratePublicPrivateKeys() (string, string) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Generation error : %s", err)
	}

	b, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		fmt.Println(err)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	}
	privateKey := base64.StdEncoding.EncodeToString(pem.EncodeToMemory(block))
	// err = ioutil.WriteFile(fb, pem.EncodeToMemory(block), 0600)
	// if err != nil {
	// 	return err
	// }

	// public key
	b, err = x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		fmt.Println(err)
	}

	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: b,
	}
	publicKey := base64.StdEncoding.EncodeToString(pem.EncodeToMemory(block))
	// a, _ := pem.Decode(b)
	// fmt.Println(a)
	// fileName := fb + ".pub"
	// err = ioutil.WriteFile(fileName, pem.EncodeToMemory(block), 0644)
	// return
	return privateKey, publicKey
}

func (s *service) UpdateOndcDetails(query map[string]interface{}, data map[string]interface{}) error {

	var OndcDetails *model_core.OndcDetails
	dto, _ := json.Marshal(data)
	_ = json.Unmarshal(dto, &OndcDetails)

	err := s.coreRepository.UpdateOndcDetails(map[string]interface{}{"id": query["id"]}, OndcDetails)
	if err != nil {
		_, err := s.coreRepository.CreateOndcDetails(OndcDetails)
		return err
	}
	return err
}

// channel status
func (s *service) ListChannelStatus(page *pagination.Paginatevalue) ([]model_core.ChannelLookupCodes, error) {
	results, er := s.coreRepository.ListChannelStatus(page)
	return results, er
}
func (s *service) ViewChannelStatus(id uint) ([]model_core.ChannelLookupCodes, error) {
	response, err := s.coreRepository.ViewChannelStatus(id)
	return response, err
}
func (s *service) CreateChannelStatus(data *model_core.ChannelLookupCodes) error {
	err := s.coreRepository.CreateChannelStatus(data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) UpdateChannelStatus(data *model_core.ChannelLookupCodes, query map[string]interface{}) error {
	err := s.coreRepository.UpdateChannelStatus(query, data)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteChannelStatus(id uint) error {
	err := s.coreRepository.DeleteChannelStatus(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) RequestSolution(data model_core.CustomSolution) error {
	err := s.coreRepository.CreateCustomSolution(data)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *service) ListCompanies(page *pagination.Paginatevalue) ([]Company, error) {
	var users []Company
	result, err := s.userRepository.FindAllCompanies(page)
	dto, _ := json.Marshal(result)
	_ = json.Unmarshal(dto, &users)
	if err != nil {
		return nil, res.BuildError(res.ErrUnprocessableEntity, err)
	}
	return users, err
}
