package webapp

import (
	"net/http"
	"strconv"

	"fermion/backend_core/controllers/omnichannel/catalogue"

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
type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) GetCatalogueFields(c echo.Context) (err error) {

	//var requestObj []map[string]interface{}
	ext_category_id, _ := strconv.Atoi(c.QueryParam("ext_category_id"))
	// if(ext_category_id == nil || ext_category_id == "undefined") {
	// ext_category_id := 126
	// }
	ext_channel_id, _ := strconv.Atoi(c.QueryParam("ext_channel_id"))
	variant_id, _ := strconv.Atoi(c.QueryParam("variant_id"))
	company_id, _ := strconv.Atoi(c.QueryParam("company_id"))
	query := map[string]interface{}{
		"category_id": ext_category_id,
		"channel_id":  ext_channel_id,
		"variant_id":  variant_id,
		"company_id":  company_id,
	}
	responseObj, _ := catalogue.NewService().GetCatalogue(query)
	requestObj := responseObj.(map[string]interface{})["channel_attributes"]

	return c.Render(http.StatusOK, "catalogues.html", map[string]interface{}{
		"name":        "HOME",
		"msg":         "Hello, Eunima",
		"request_obj": requestObj,
		"category_id": ext_category_id,
		"channel_id":  ext_channel_id,
		"variant_id":  variant_id,
		"company_id":  company_id,
	})
}
