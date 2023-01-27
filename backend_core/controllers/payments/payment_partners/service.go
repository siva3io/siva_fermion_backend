package paymentpartners

type Service interface {
	ListPaymentPartners() []map[string]interface{}
}
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
type service struct {
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	newServiceObj = &service{}
	return newServiceObj
}

func (s *service) ListPaymentPartners() []map[string]interface{}{
	var paymentpartner1 = map[string]interface{}{
		"partner_name":"Stripe India",
		"type":"Domestic",
		"logo":"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQGluJhW7I1NYU7jF77E-9K9I46_ib_DUNHw&usqp=CAU",
	}
	var paymentpartner2 = map[string]interface{}{
		"partner_name":"Stripe US",
		"type":"International",
		"logo":"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQQGluJhW7I1NYU7jF77E-9K9I46_ib_DUNHw&usqp=CAU",
	}
	return []map[string]interface{}{
		paymentpartner1,
		paymentpartner2}
}
