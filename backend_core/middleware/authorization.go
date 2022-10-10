package middleware

import (
	"fmt"
	"os"
	"strings"

	res "fermion/backend_core/pkg/util/response"

	"github.com/golang-jwt/jwt/v4"
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
		user_id := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["ID"])
		user_name := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["Username"])
		access_template_id := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["access_template_id"])
		company_id := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["company_id"])

		c.Set("Username", user_name)
		c.Set("TokenUserID", user_id)
		c.Set("AccessTemplateId", access_template_id)
		c.Set("CompanyId", company_id)

		// fmt.Println("User Name :" + user_name+"\nUser ID :" + user_id+"\nAccess Template ID :" + access_template_id+"\nCompany ID :" + company_id)

		return next(c)
	}
}
