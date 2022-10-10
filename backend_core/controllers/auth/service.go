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
	mdm_service "fermion/backend_core/controllers/mdm/contacts"
	model_core "fermion/backend_core/internal/model/core"
	mdm_model "fermion/backend_core/internal/model/mdm"
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
	// UpdateUserProfile(query map[string]interface{}, data *model_core.CoreUsers) error
}

type service struct {
	userRepository    repository.User
	coreRepository    repository.Core
	contactRepository mdm_repo.Contacts
	contactService    mdm_service.Service
}

func NewService() *service {
	userRepository := repository.NewUser()
	// coreService := core_service.NewService()
	coreRepository := repository.NewCore()
	contactRepository := mdm_repo.NewContacts()
	contactService := mdm_service.NewService()
	return &service{userRepository, coreRepository, contactRepository, contactService}
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

		return map[string]interface{}{"id": resp.ID, "status": "user_created", "new_user": true, "otp_sent": true}, nil
	}

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
	var user UserObject

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
				access_id := fmt.Sprintf("%v", access_ids[0])
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
	contactResponse, _ := s.contactRepository.FindOneContact(query)
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

		err = s.contactRepository.SaveContact(&contact)
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
	response, err := s.userRepository.FindUser(map[string]interface{}{"id": uint(ID)})
	if err != nil {
		return UserResponse, err
	}
	dto, err := json.Marshal(&response)
	if err != nil {
		return UserResponse, err
	}
	json.Unmarshal(dto, &UserResponse)
	result, err := s.contactService.GetContact(map[string]interface{}{"primary_email": UserResponse.Email}, id, access_template_id)
	jsonresult, err := json.Marshal(&result)
	if err != nil {
		return UserResponse, err
	}
	var contactDetails map[string]interface{}
	json.Unmarshal(jsonresult, &contactDetails)
	defaultAddress, _ := json.Marshal(contactDetails["default_address"])
	UserResponse.AddressDetails = defaultAddress

	if contactDetails["parent"] != nil {
		parentDetails := contactDetails["parent"].(map[string]interface{})
		if parentDetails["first_name"] != nil {
			UserResponse.TeamHead = parentDetails["first_name"].(string)
		}
		if parentDetails["last_name"] != nil {
			UserResponse.TeamHead = UserResponse.TeamHead + " " + parentDetails["last_name"].(string)
		}
	}

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
