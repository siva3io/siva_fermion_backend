package accounting

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
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
type CurrencyExchange struct {
	model_core.Model
	QuoteCurrencyID *uint                 `json:"quote_currency_id" gorm:"type:integer"`
	QuoteCurrency   *model_core.Currency  `json:"quote_currency" gorm:"foreignKey:QuoteCurrencyID; references:ID"`
	QuoteCountryID  *uint                 `json:"quote_country_id" gorm:"type:integer"`
	QuoteCountry    *model_core.Country   `json:"country" gorm:"foreignKey:QuoteCountryID; references:ID"`
	ExchangeId      *uint                 `json:"exchange_id"`
	Exchange        model_core.Lookupcode `gorm:"foreignKey:ExchangeId;references:ID" json:"exchange"`
	BreakDown       []BreakDownIntervals  `json:"breakdown_intervals" gorm:"foreignkey:BreakDownId; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type BreakDownIntervals struct {
	model_core.Model
	BreakDownId     uint
	ConversionRatio float64   `json:"conversion_ratio"`
	StartDate       string    `json:"start_date" gorm:"type:date"`
	StartTime       time.Time `json:"start_time"`
	EndDate         string    `json:"end_date" gorm:"type:date"`
	EndTime         time.Time `json:"end_time"`
}
