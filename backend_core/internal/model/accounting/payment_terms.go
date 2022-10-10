package accounting

import model_core "fermion/backend_core/internal/model/core"

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
type PaymentTerms struct {
	model_core.Model
	PaymentTermName      string                `gorm:"unique; not null" json:"payment_term_name"`
	DescriptionOnInvoice string                `json:"description_on_invoice"`
	AccountTypeId        uint                  `json:"account_type_id"`
	AccountType          model_core.Lookupcode `gorm:"foreignkey:AccountTypeId;references:ID" json:"account_type"`
	PaymentTermDetails   []PaymentTermDetails  `gorm:"foreignkey:PaymentTermsId;references:ID" json:"term_details"`
}

type PaymentTermDetails struct {
	model_core.Model
	PaymentTermsId  uint                  `json:"-"`
	PaymentTermNo   uint                  `json:"payment_term_no" gorm:"unique"`
	TermTypeId      uint                  `json:"term_type_id"`
	TermType        model_core.Lookupcode `gorm:"foreignkey:TermTypeId;references:ID" json:"term_type"`
	Value           uint                  `gorm:"not null" json:"value"`
	Due             uint                  `json:"due"`
	EventType       string                `json:"event_type"`
	OnTheDayOfMonth string                `json:"on_the_day_of_month"`
}
