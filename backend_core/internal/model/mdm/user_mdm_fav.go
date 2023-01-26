package mdm

import (
	model_core "fermion/backend_core/internal/model/core"

	"github.com/lib/pq"
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
type UserMdmFav struct {
	model_core.Model
	UserID     uint                 `json:"user_id"`
	User       model_core.CoreUsers `json:"user" gorm:"foreignkey:UserID;references:ID"`
	ContactIds pq.Int64Array        `json:"contact_ids" gorm:"type:int[]"`
	//Contact    []*Partner        `json:"contact" gorm:"many2many:fav_contacts"`
	PricingIds pq.Int64Array `json:"pricing_ids" gorm:"type:int[]"`
	//Pricing    []*Pricing        `json:"pricing" gorm:"many2many:fav_pricings"`
	UomIds pq.Int64Array `json:"uom_ids" gorm:"type:int[]"`
	//Uom        []*Uom            `json:"uom" gorm:"many2many:fav_uom"`
	ProductIds pq.Int64Array `json:"product_ids" gorm:"type:int[]"`
	//Product    []*ProductVariant `json:"product" gorm:"many2many:fav_products"`
	VendorIds pq.Int64Array `json:"vendor_ids" gorm:"type:int[]"`
	//Vendor     []*Vendors        `json:"vendor" gorm:"many2many:fav_vendors"`
	LocationIds pq.Int64Array `json:"location_ids" gorm:"type:int[]"`
}
