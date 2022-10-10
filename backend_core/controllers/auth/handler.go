package auth

import (
	"strconv"

	// "fermion/backend_core/pkg/util/helpers"
	model_core "fermion/backend_core/internal/model/core"

	// "fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/labstack/echo/v4"
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
type handler struct {
	service Service
}

func NewHandler() *handler {

	service := NewService()
	return &handler{service}

	// coreRepository := repository.NewCore()

	// return &handler{service, coreRepository}
}

// Register godoc
// @Summary Register accounts
// @Description Register accounts
// @Tags auth
// @Accept  json
// @Produce  json
// @param register body RegisterDTO true "request body register"
// @Success 200 {object} RegisterResponseDOC
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /auth/register [post]
func (h *handler) Register(c echo.Context) (err error) {

	dto := c.Get("register_data").(*RegisterDTO)
	user, err := h.service.Register(dto)

	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccess(c, "Register success", user)
}

// Login godoc
// @Summary Login accounts
// @Description Login accounts
// @Tags auth
// @Accept  json
// @Produce  json
// @param register body LoginDTO true "request body login"
// @Success 200 {object} LoginResponseDOC
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /auth/login [post]
func (h *handler) Login(c echo.Context) (err error) {

	dto := c.Get("login_data").(*LoginDTO)
	data, err := h.service.Login(dto)

	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Login success", data)
}

// UserLogin godoc
// @Summary UserLogin
// @Description UserLogin
// @Tags auth
// @Accept  json
// @Produce  json
// @param login body UserLoginDTO true "request body login"
// @Success 200 {object} UserLoginResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /auth/user_login [post]
func (h *handler) UserLogin(c echo.Context) error {
	var data = make(map[string]interface{})
	c.Bind(&data)
	response, err := h.service.UserLogin(data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "success", response)
}

// VerifyOTP godoc
// @Summary VerifyOTP
// @Description VerifyOTP
// @Tags auth
// @Accept  json
// @Produce  json
// @param verifyotp body VerifyOtpDTO true "verify otp"
// @Success 200 {object} VerifyOtpResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /auth/verify_otp [post]
func (h *handler) VerifyOTP(c echo.Context) error {
	var data VerifyOtpDTO
	c.Bind(&data)
	response, err := h.service.VerifyOtp(data)
	if err != nil {
		return res.RespError(c, err)
	}
	return res.RespSuccess(c, "success", response)
}
func (h *handler) AssignTemplate(c echo.Context) (err error) {

	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	var data model_core.CoreUsers
	c.Bind(&data)
	// s := c.Get("TokenUserID").(string)
	// data.UpdatedByID = helpers.ConvertStringToUint(s)
	err = h.service.UpdateAssignTemplate(uint(id), data)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Module details updates successfully", map[string]interface{}{"updated_id": id})
}

// func (h *handler) UpdateUserProfile(c echo.Context) (err error) {
// 	ID := c.Param("id")
// 	id, _ := strconv.Atoi(ID)
// 	query := map[string]interface{}{
// 		"id": uint(id),
// 	}
// 	host := c.Request().Host
// 	scheme := c.Scheme()
// 	fmt.Println(query)
// 	var data UserProfileUpdateDTO
// 	if err := c.Bind(&data); err != nil {
// 		return err
// 	}
// 	err = h.service.UpdateUserProfile(query, data, scheme, host)
// 	return res.RespSuccess(c, "User details updates successfully", map[string]interface{}{"updated_id": id})
// }

func (h *handler) UpdateProfile(c echo.Context) error {
	profile_data := c.Get("profile_data").(*UpdateProfileDTO)

	err := h.service.UpdateProfile(*profile_data)

	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "profile updated successfully", profile_data)
}
func (h *handler) GetUserProfile(c echo.Context) (err error) {
	id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.GetUser(id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "user profile fetched successfully", result)
}
func (h *handler) GetUserById(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	result, err := h.service.GetUserById(uint(ID))
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "user details fetched successfully", result)
}
