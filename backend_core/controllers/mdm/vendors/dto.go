package vendors

import (
	"time"

	"fermion/backend_core/pkg/util/response"
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
type VendorsListDTO struct {
	Meta response.MetaResponse
	Data []VendorsRequestDTO
}
type VendorDTO struct {
	Meta response.SuccessResponse
	Data VendorsRequestDTO
}
type VendorsObjDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type VendorsRequestDTO struct {
	Name             string                `json:"name"`
	ContactId        uint                  `json:"contact_id"`
	VendorDetails    Report                `json:"vendor_details"`
	PrimaryContactId uint                  `json:"primary_contact_id"`
	VendorDocs       Report                `json:"vendor_docs"`
	VendorPriceLists []VendorPriceListsDTO `json:"vendor_price_lists"`
}

type VendorPriceListsDTO struct {
	Name             string    `json:"name"`
	ProductVariantId uint      `json:"product_variant_id"`
	MinQty           int64     `json:"min_qty"`
	SupplyLeadTime   string    `json:"supply_lead_time"`
	PriceListId      uint      `json:"price_list_id"`
	ValidFrom        time.Time `json:"valid_from"`
	ValidTo          time.Time `json:"valid_to"`
	CreditPeriod     string    `json:"credit_period"`
	PaymentTypeId    string    `json:"payment_type_id"`
}

type Report struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
