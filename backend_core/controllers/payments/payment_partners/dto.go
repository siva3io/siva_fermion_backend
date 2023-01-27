package paymentpartners
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
type PaymentPartnerDTO struct {
	PartnerName                   string                 `json:"partner_name"`
	PaymentPartneDTOPartnerTypeId uint                   `json:"shipping_partner_type_id"`
	ProfileOptions                map[string]interface{} `json:"profile_options"`
	AuthOptions                   map[string]interface{} `json:"auth_options"`
	SubscriptionOptions           map[string]interface{} `json:"subscription_options"`
	Id                            uint                   `json:"id"`
	AppId                         *uint                  `json:"app_id"`
}
