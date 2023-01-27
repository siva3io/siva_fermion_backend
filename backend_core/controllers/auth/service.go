package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	// "strconv"

	// "errors"
	"regexp"
	"strconv"
	"strings"

	// core_service "fermion/backend_core/controllers/cores"
	template_service "fermion/backend_core/controllers/access/template"
	"fermion/backend_core/controllers/eda"
	contact_service "fermion/backend_core/controllers/mdm/contacts"
	location_service "fermion/backend_core/controllers/mdm/locations"
	"fermion/backend_core/internal/model/core"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
	mdm_model "fermion/backend_core/internal/model/mdm"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/repository"
	mdm_repo "fermion/backend_core/internal/repository/mdm"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"golang.org/x/crypto/bcrypt"
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
	//-------------------Basic Auth----------------------------------------------------------------
	Register(dto *RegisterDTO) (interface{}, error)
	Login(dto *LoginDTO) (LoginResponseDTO, error)

	//-----------------------OTP--------------------------------------------------------------------
	UserLogin(login map[string]interface{}) (map[string]interface{}, error)
	SendOtp(user model_core.CoreUsers) error
	VerifyOtp(dto VerifyOtpDTO) (interface{}, error)
	UpdateAssignTemplate(id uint, data model_core.CoreUsers) error
	// UpdateUserProfile(query map[string]interface{}, data UserProfileUpdateDTO, scheme string, host string) error

	UpdateProfile(data UpdateProfileDTO) error
	GetUser(id string, access_template_id string) (UserResponseDTO, error)
	GetUserById(id uint) (UserVisitDTO, error)
	UpdateUser(id uint, data model_core.CoreUsers) error
	FindAllUsers(query map[string]interface{}, p *pagination.Paginatevalue) ([]model_core.CoreUsers, error)
	// UpdateUserProfile(query map[string]interface{}, data *model_core.CoreUsers) error
}

type service struct {
	userRepository        repository.User
	coreRepository        repository.Core
	contactRepository     mdm_repo.Contacts
	contactService        contact_service.Service
	locationService       location_service.Service
	accessTemplateService template_service.Services
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	userRepository := repository.NewUser()
	coreRepository := repository.NewCore()
	contactRepository := mdm_repo.NewContacts()
	contactService := contact_service.NewService()
	locationService := location_service.NewService()
	accessTemplateService := template_service.NewService()
	newServiceObj = &service{userRepository, coreRepository, contactRepository, contactService, locationService, accessTemplateService}
	return newServiceObj
}

// -------------------Basic Auth----------------------------------------------------------------
func (s *service) Register(dto *RegisterDTO) (interface{}, error) {
	query := map[string]interface{}{"username": dto.Username}
	user, err := s.userRepository.FindUser(query)
	if err == nil {
		return nil, err
	}

	user.FirstName = dto.FirstName
	user.LastName = dto.LastName
	user.Username = dto.Username
	user.Password = dto.Password
	user.Name = user.FirstName + " " + user.LastName

	user.HashPassword()

	err = s.userRepository.Save(&user)
	if err != nil {
		return nil, err
	}

	value, _ := json.Marshal(user)
	var response RegisterRespData
	json.Unmarshal(value, &response)

	return response, nil
}
func (s *service) Login(dto *LoginDTO) (LoginResponseDTO, error) {
	var resDto LoginResponseDTO

	query := map[string]interface{}{
		"username": dto.Username,
	}
	user, err := s.userRepository.FindUser(query)
	if err != nil {
		return resDto, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)) != nil {
		return resDto, errors.New("incorrect password")
	}

	token, err := user.GenerateToken(0, 72) //change the passing value based on the login access template selected and why is it needed in the login
	if err != nil {
		return resDto, res.BuildError(res.ErrServerError, err)
	}

	resDto.Token = token
	return resDto, nil
}

// -------------------OTP----------------------------------------------------------------
func (s *service) UserLogin(login map[string]interface{}) (map[string]interface{}, error) {
	resp, err := s.userRepository.FindLogin(login)

	//---------------create new user----------------------------------------------------
	if err != nil && err.Error() == "record not found" {
		if regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(login["login"].(string)) {
			resp.Email = strings.ToLower(login["login"].(string))
		} else if regexp.MustCompile(`^[0-9]+$`).MatchString(login["login"].(string)) {
			resp.MobileNumber = login["login"].(string)
		} else {
			return map[string]interface{}{
				"status":  false,
				"message": "USER_NOT_FOUND",
			}, nil
		}
		//------------create company----------------------------------------------------
		var company model_core.Company
		if resp.Email != "" {
			company.Name = resp.Email
		}
		if resp.MobileNumber != "" {
			company.Name = resp.MobileNumber
		}
		query := map[string]interface{}{
			"name": company.Name,
		}
		companyResponse, _ := s.coreRepository.FindCompany(query)
		if companyResponse.Name == "" {
			data, err := s.coreRepository.FindCompany(map[string]interface{}{"name": "Eunimart Ltd"})
			if err != nil {
				return nil, err
			}
			json_data, _ := json.Marshal(data.CompanyDefaults)
			json.Unmarshal(json_data, &company.CompanyDefaults)
			company.FilePreferenceID = data.FilePreferenceID
			err = s.coreRepository.CreateCompany(&company)
			if err != nil {
				return nil, err
			}
			fmt.Println("--------------------->>> Company Created Successfully <<<-----------------------------")
		}
		resp.AccessTemplateId = append(resp.AccessTemplateId, &model_core.AccessTemplate{ID: 2})
		resp.CompanyId = &company.ID
		err := s.userRepository.Save(&resp)
		if err != nil {
			return nil, err
		}
		fmt.Println("--------------------->>> user Created Successfully <<<-----------------------------")

		//-----send OTP------------------

		otp_err := s.SendOtp(resp)
		if otp_err != nil {
			return nil, otp_err
		}
		//------------------create default location-------------------------------------------------
		var virtual_location location_service.LocationRequestDTO
		virtual_location.LocationTypeId, _ = helpers.GetLookupcodeId("LOCATION_TYPE", "VIRTUAL_LOCATION")
		virtual_location.Name = "Default Location"
		virtual_location.CompanyId = resp.CompanyId

		//-----fetching company admin access id-------------------
		p := new(pagination.Paginatevalue)
		p.Filters = "[[\"name\",\"=\",\"COMPANY_ADMIN\"]]"
		access_template_list, _ := s.accessTemplateService.GetAllTemplateList(p)

		//creating location
		var metaData core.MetaData
		if len(access_template_list) > 0 {
			metaData.AccessTemplateId = access_template_list[0].ID
		}
		metaData.TokenUserId = resp.ID
		metaData.CompanyId = company.ID

		request_payload := map[string]interface{}{
			"meta_data": metaData,
			"data":      virtual_location,
		}
		eda.Produce(eda.CREATE_LOCATION, request_payload)
		fmt.Println("--------------------->>> Default Location Creation InProgress  <<<-----------------------------")

		return map[string]interface{}{"id": resp.ID, "status": "user_created", "new_user": true, "otp_sent": true}, nil
	}
	// location_id := location_data.(shared_pricing_and_location.Locations).ID

	//---------If user is already existing then send OTP directly---------------------------------
	otp_err := s.SendOtp(resp)
	if otp_err != nil {
		return nil, otp_err
	}

	return map[string]interface{}{"id": resp.ID, "new_user": false, "otp_sent": true}, nil
}
func (s *service) SendOtp(user model_core.CoreUsers) error {
	otp_token, err := helpers.GenerateOTP()
	if err != nil {
		return err
	}
	var auth model_core.Auth
	var query = make(map[string]interface{})
	auth.OtpToken = otp_token
	value, _ := json.Marshal(auth)
	user.Auth = value
	query["id"] = user.ID

	if user.Email != "" {
		query["email"] = user.Email
	}
	if user.MobileNumber != "" {
		query["mobile_number"] = user.MobileNumber
	}
	er := s.userRepository.UpdateUser(query, &user)
	if er != nil {
		return er
	}

	//TODO: implement function to send otp as email and sms

	return nil
}
func (s *service) VerifyOtp(dto VerifyOtpDTO) (interface{}, error) {
	response, err := s.userRepository.FindUser(map[string]interface{}{"id": dto.ID})
	if err != nil {
		return nil, err
	}
	var auth model_core.Auth
	var resp VerifyOtpResponse
	var user UserResponseDTO

	value, _ := json.Marshal(response)
	_ = json.Unmarshal(response.Auth, &auth)
	resp.Status = false
	resp.Message = "INVALID_OTP"
	if auth.OtpToken != "" {
		if auth.OtpToken != dto.Otp {
			resp.Status = false
			resp.Message = "INVALID_OTP"
		}
		if auth.OtpToken == dto.Otp {
			resp.Status = true
			resp.Message = "OTP_VALID"
			if dto.AccessTemplateId == 0 && len(response.AccessIds) < 3 {
				dto.AccessTemplateId = 2
			} else if len(response.AccessIds) > 2 {

				var access_ids []interface{}
				value, _ := json.Marshal(response.AccessIds)
				json.Unmarshal(value, &access_ids)
				var access_id string
				if len(access_ids) > 0 {
					access_id = fmt.Sprintf("%v", access_ids[0])
				}
				dto.AccessTemplateId = *helpers.ConvertStringToUint(access_id)
			}
			resp.Token, _ = response.GenerateToken(dto.AccessTemplateId, 72)
			json.Unmarshal(value, &user)
			resp.User = user
			s.UpdateAssignTemplate(dto.AccessTemplateId, response)
			resp.AccessTemplateId = dto.AccessTemplateId
		}
	}

	return resp, nil
}
func (s *service) UpdateAssignTemplate(id uint, data model_core.CoreUsers) error {

	query := map[string]interface{}{"id": id}
	_, er := s.userRepository.FindUser(query)
	if er != nil {
		return er
	}

	count, err := s.userRepository.UpdateAssignTemplate(id, data)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("no data gets updated")
	}

	return nil
}

// func (s *service) UpdateUserProfile(query map[string]interface{}, data UserProfileUpdateDTO, scheme string, host string) error {
// 	var user model_core.CoreUsers
// 	result, _ := s.userRepository.FindUser(query)
// 	var user_companyId = make(map[string]interface{})
// 	user_companyId["id"] = result.CompanyId
// 	company, _ := s.coreService.GetCompanyPreferences(user_companyId)
// 	serviceProvider := company.FilePreference.LookupCode
// 	dto, err := json.Marshal(&data)
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Unmarshal(dto, &user)
// 	if err != nil {
// 		return err
// 	}
// 	if data.ProfileDpImage != "" && serviceProvider != "LOCAL" {
// 		var imageLinkArray = make([]map[string]interface{}, 0)
// 		imageLinkArray = append(imageLinkArray, map[string]interface{}{"data": data.ProfileDpImage})
// 		response, _ := helpers.UploadFile(imageLinkArray, "profileimage", strconv.Itoa(int(data.ID)), data.ID, scheme, host, helpers.AWS, "Profile", "Profile", ".png", helpers.BASE64)
// 		link, _ := json.Marshal(&response)
// 		json.Unmarshal(link, &user.Profile)
// 	} else {
// 		link, _ := json.Marshal(map[string]interface{}{"data": data.ProfileDpImage})
// 		json.Unmarshal(link, &user.Profile)
// 	}
// 	err = s.userRepository.UpdateUser(query, user)
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }

func (s *service) UpdateProfile(data UpdateProfileDTO) error {
	query := map[string]interface{}{
		"id": data.ID,
	}
	user, err := s.userRepository.FindUser(query)
	if err != nil || user.ID == 0 {
		return errors.New("user not found")
	}

	//====================update company===================================
	if user.MobileNumber != "" {
		query = map[string]interface{}{
			"name": user.MobileNumber,
		}
	}
	if user.Email != "" {
		query = map[string]interface{}{
			"name": user.Email,
		}
	}

	companyResponse, _ := s.coreRepository.FindCompany(query)
	if companyResponse.Name != "" {
		var company model_core.Company
		company.Name = data.CompanyName
		company.Email = data.CompanyEmail

		err := s.coreRepository.UpdateCompany(query, &company)
		if err != nil {
			return err
		}
		fmt.Println("--------------------->>> Company Updated Successfully <<<-----------------------------")
	}

	//=======================update user profile details==========================================
	var userDetails model_core.CoreUsers
	jsondata, _ := json.Marshal(data)
	json.Unmarshal(jsondata, &userDetails)
	userDetails.Name = userDetails.FirstName + " " + userDetails.LastName
	userDetails.HashPassword()

	query = map[string]interface{}{
		"id": data.ID,
	}
	err = s.userRepository.UpdateUser(query, &userDetails)
	if err != nil {
		return err
	}
	fmt.Println("--------------------->>> Profile Updated Successfully <<<-----------------------------")

	//=================create contact===============================================
	var contact mdm_model.Partner
	query = map[string]interface{}{"primary_email": data.Email}
	contactResponse, _ := s.contactRepository.FindOne(query)
	if contactResponse.ID == 0 {

		defaultContactTypeID, err := helpers.GetLookupcodeId("CONTACT_TYPE", "INDIVIDUAL")
		if err != nil {
			return err
		}
		if contactResponse.ContactTypeId == nil {
			contact.ContactTypeId = &defaultContactTypeID
		}
		contact.CreatedByID = &data.ID
		contact.PrimaryEmail = data.Email
		contact.PrimaryPhone = data.MobileNumber
		contact.FirstName = data.FirstName
		contact.LastName = data.LastName
		contact.ImageOptions = data.Profile
		contact.UserId = data.ID

		err = s.contactRepository.Create(&contact)
		if err != nil {
			return err
		}
		fmt.Println("--------------------->>> Contact Created Successfully <<<-----------------------------")
	}
	return nil
}

func (s *service) GetUser(id string, access_template_id string) (UserResponseDTO, error) {
	var UserResponse UserResponseDTO
	ID, _ := strconv.Atoi(id)
	accessTemplateId, _ := strconv.Atoi(access_template_id)

	response, err := s.userRepository.FindUser(map[string]interface{}{"id": uint(ID)})
	if err != nil {
		return UserResponse, err
	}
	helpers.JsonMarshaller(response, &UserResponse)

	//============ fetch address details from contacts ========================
	var metaData model_core.MetaData
	metaData.TokenUserId = uint(ID)
	metaData.AccessTemplateId = uint(accessTemplateId)
	metaData.Query = make(map[string]interface{})
	metaData.Query["user_id"] = UserResponse.ID

	result, _ := s.contactService.FindOne(metaData)

	var contactDetails mdm.Partner
	helpers.JsonMarshaller(result, &contactDetails)
	var addressDetails []map[string]interface{}
	helpers.JsonMarshaller(contactDetails.AddressDetails, &addressDetails)
	for _, address := range addressDetails {
		if address["mark_default_address"] != nil && address["mark_default_address"] == true {
			UserResponse.AddressDetails, _ = json.Marshal(address)
			break
		}
	}
	if contactDetails.ParentId != nil && *contactDetails.ParentId != 0 {
		UserResponse.TeamHead = contactDetails.Parent.FirstName + " " + contactDetails.Parent.LastName
	}

	UserResponse.Token, _ = response.GenerateToken(uint(UserResponse.AccessIds[0]), 72)
	return UserResponse, nil
}
func (s *service) GetUserById(id uint) (UserVisitDTO, error) {
	var UserResponse UserVisitDTO
	response, err := s.userRepository.FindUser(map[string]interface{}{"id": id})
	if err != nil {
		return UserResponse, err
	}
	dto, err := json.Marshal(&response)
	if err != nil {
		return UserResponse, err
	}
	json.Unmarshal(dto, &UserResponse)
	return UserResponse, nil
}

func (s *service) UpdateUser(id uint, data model_core.CoreUsers) error {
	// _, er := s.templateRepository.FindByidRoles(id)

	// if er != nil {
	// 	return er
	// }

	count, err := s.userRepository.UpdateUserDetails(id, data)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("no data gets updated")
	}

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *service) FindAllUsers(query map[string]interface{}, p *pagination.Paginatevalue) ([]model_core.CoreUsers, error) {
	data, err := s.userRepository.FindAllCoreUser(query, p)
	if err != nil {
		return data, err
	}
	return data, nil

}
