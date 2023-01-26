package payment_terms_and_record_payment

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
type PaymentTermsDTO struct {
	PaymentTermName      string               `gorm:"unique;" json:"payment_term_name"`
	DescriptionOnInvoice string               `json:"description_on_invoice"`
	AccountTypeId        uint                 `json:"account_type_id"`
	PaymentTermDetails   []PaymentTermDetails `gorm:"foreignkey:PaymentTermsId;references:ID" json:"term_details"`
}

type PaymentTermDetails struct {
	PaymentTermsId  uint   `json:"-"`
	PaymentTermNo   uint   `json:"payment_term_no" gorm:"unique"`
	TermTypeId      uint   `json:"term_type_id"`
	Value           uint   `gorm:"not null" json:"value"`
	Due             uint   `json:"due"`
	EventType       string `json:"event_type"`
	OnTheDayOfMonth string `json:"on_the_day_of_month"`
	EventId         *uint  `json:"event_id"`
}
