package middleware

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"fermion/backend_core/internal/model/core"
	res "fermion/backend_core/pkg/util/response"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
func Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	var (
		jwtKey = os.Getenv("JWT_KEY")
	)

	return func(c echo.Context) error {
		authToken := c.Request().Header.Get("Authorization")
		if authToken == "" {
			return res.RespError(c, &res.ErrUnauthorized)
		}

		splitToken := strings.Split(authToken, "Bearer ")

		if len(splitToken) != 2 {
			return res.RespError(c, &res.ErrUnauthorized)
		}

		token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(jwtKey), nil
		})

		if err != nil || !token.Valid {
			return res.RespError(c, &res.ErrUnauthorized)
		}

		requestId := uuid.New().String()
		userId := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["ID"])
		userName := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["Username"])
		accessTemplateId := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["access_template_id"])
		companyId := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["company_id"])

		c.Set("Username", userName)
		c.Set("TokenUserID", userId)
		c.Set("AccessTemplateId", accessTemplateId)
		c.Set("CompanyId", companyId)
		c.Set("RequestId", requestId)

		// temporarily above c.Set() will be there until ideal format complete
		tokenUserId, _ := strconv.Atoi(userId)
		access_template_id, _ := strconv.Atoi(accessTemplateId)
		company_id, _ := strconv.Atoi(companyId)
		var metaData = core.MetaData{
			RequestId:        requestId,
			Host:             c.Request().Host,
			Scheme:           c.Scheme(),
			TokenUserId:      uint(tokenUserId),
			AccessTemplateId: uint(access_template_id),
			CompanyId:        uint(company_id),
		}
		c.Set("MetaData", metaData)

		// fmt.Println("User Name :" + user_name+"\nUser ID :" + user_id+"\nAccess Template ID :" + access_template_id+"\nCompany ID :" + company_id)

		return next(c)
	}
}
