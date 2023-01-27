package paymentpartners

import (
	res "fermion/backend_core/pkg/util/response"

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
	service Service
}

var CustomersHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if CustomersHandler != nil {
		return CustomersHandler
	}
	service := NewService()
	CustomersHandler = &handler{service}
	return CustomersHandler
}

func (h *handler) ListPaymentPartners(c echo.Context) error {
	data := h.service.ListPaymentPartners()
	return res.RespSuccess(c, "Payment Partners Listed Successfully", data)

}
