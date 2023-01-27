package omnichannel

import (
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm/shared_pricing_and_location"

	"gorm.io/datatypes"
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
type (
	Marketplace struct {
		model_core.Model
		Code               string         `json:"code" gorm:"unique"`
		Name               string         `json:"name" gorm:"varchar(50)"`
		Image_options      datatypes.JSON `json:"image_options" gorm:"json"`
		Location_served    datatypes.JSON `json:"location_served" gorm:"json"`
		MarketPlaceDetails datatypes.JSON `json:"market_place_details" gorm:"json"`
	}
)

// register for marketplace
type (
	User_Marketplace_Registration struct {
		model_core.Model
		MarketplaceCode string      `json:"marketplace_code" gorm:"type:text"`
		Marketplace     Marketplace `json:"marketplace" gorm:"foreignKey:MarketplaceCode; references:Code"`
		Full_name       string      `json:"full_name" gorm:"varchar"`
		Email           string      `json:"email" gorm:"varchar"`
		Mobile_no       int64       `json:"mobile_no" `
		Date_of_birth   string      `json:"date_of_birth" gorm:"date"`
		Bank_info       Bank_info   `gorm:"embedded"`
		KYC_info        Kyc_info    `gorm:"embedded"`
		Status          string      `json:"status"`
		Status_history  string      `json:"status_history" gorm:"status"`
	}

	Bank_info struct {
		Bank_name        string         `json:"bank_name" gorm:"varchar"`
		Holder_name      string         `json:"holder_name" gorm:"varchar"`
		Account_number   int64          `json:"account_number"`
		Ifsc_code        string         `json:"ifsc_code" gorm:"varchar"`
		Bank_statement   datatypes.JSON `json:"bank_statement" gorm:"json"`
		Cancelled_cheque datatypes.JSON `json:"cancelled_cheque" gorm:"json"`
		// Bank_statement_id   int64       `json:"bank_statement" gorm:"int"`
		// Attachment          Attachments `json:"-" gorm:"foreignKey:bank_statement_id; references:ID"`
		// Cancelled_cheque_id int64       `json:"cancelled_cheque" gorm:"int"`
		// Attachments         Attachments `json:"-" gorm:"foreignKey:cancelled_cheque_id; references:ID"`
	}

	Kyc_info struct {
		Iec_id         datatypes.JSON `json:"iec" gorm:"json"`
		Passport_id    datatypes.JSON `json:"passport" gorm:"json"`
		Adhaar_card_id datatypes.JSON `json:"adhaar_card" gorm:"json"`
		Voter_id       datatypes.JSON `json:"voter_id" gorm:"json"`
		Pan_card_id    datatypes.JSON `json:"pan_card" gorm:"json"`
		Gstin_id       datatypes.JSON `json:"gstin" gorm:"json"`
		// Iec_id         int64       `json:"iec" gorm:"int"`
		// IEC            Attachments `json:"-" gorm:"foreignKey:Iec_id; references:ID" `
		// Passport_id    int64       `json:"passport" gorm:"int"`
		// Passport       Attachments `json:"-" gorm:"foreignKey:Passport_id; references:ID" `
		// Adhaar_card_id int64       `json:"adhaar_card" gorm:"int"`
		// Adhaar         Attachments `json:"-" gorm:"foreignKey:Adhaar_card_id; references:ID" `
		// Voter_id       int64       `json:"voter_id" gorm:"int"`
		// Voter          Attachments `json:"-" gorm:"foreignKey:Voter_id; references:ID" `
		// Pan_card_id    int64       `json:"pan_card" gorm:"int"`
		// Pan_card       Attachments `json:"-" gorm:"foreignKey:Pan_card_id; references:ID" `
		// Gstin_id       int64       `json:"gstin" gorm:"int"`
		// Gst            Attachments `json:"-" gorm:"foreignKey:Gstin_id; references:ID" `
	}
)

// link marketplace
type (
	User_Marketplace_Link struct {
		model_core.Model
		MarketPlaceCode                string                         `json:"marketplace_code" gorm:"unique"`
		MarketPlace                    Marketplace                    `json:"marketPlace" gorm:"foreignKey:MarketPlaceCode; references:Code"`
		Auth                           datatypes.JSON                 `gorm:"type:json;default:'[]';not null" json:"auth"`
		Channel_details                Channel_details                `gorm:"embedded"`
		Fullfilment_details            Fullfilment_details            `gorm:"embedded"`
		Inventory_details              Inventory_details              `gorm:"embedded"`
		Order_details                  Order_details                  `gorm:"embedded"`
		Advanced_channel_configuration Advanced_channel_configuration `gorm:"embedded"`
		Payment_mapping                Payment_mapping                `gorm:"embedded"`
		Pre_sync_configuration         Pre_sync_configuration         `gorm:"embedded"`
		Sync_configuration             Sync_configuration             `gorm:"embedded"`
		Sync_fetch_data                Sync_fetch_data                `gorm:"embedded"`
		Inventory_Automation_Details   Inventory_Automation_Details   `gorm:"embedded"`
	}
	Channel_details struct {
		Price_list_id *uint                                `json:"price_list_id" gorm:"type:integer"`
		Price_list    *shared_pricing_and_location.Pricing `json:"price_list" gorm:"foreignKey:Price_list_id; references:ID"`
		Currency_id   *uint                                `json:"currency_id" gorm:"type:integer"`
		Currency      *model_core.Currency                 `json:"currency" gorm:"foreignKey:Currency_id; references:ID"`
		// OrderTagIds   []*model.Lookupcode `json:"order_tag_ids" gorm:"many2many:retail_order_Tags;foreignKey:ID;references:ID"`
		Order_tags_id *uint                 `json:"order_tags_id" gorm:"type:integer"`
		Order_tags    model_core.Lookupcode `json:"order_tags" gorm:"foreignKey:Order_tags_id; references:ID"`
	}
	Fullfilment_details struct {
		Enable_fba_fbf_to_fullfil_orders *bool          `gorm:"type:boolean"`
		Select_the_3pl_facility          datatypes.JSON `gorm:"type:varchar"`
	}
	Inventory_details struct {
		Select_source_facility_id *uint                                  `json:"select_source_facility_id" gorm:"type:integer"`
		Select_source_facility    *shared_pricing_and_location.Locations `json:"select_source_facility" gorm:"foreignKey:Select_source_facility_id; references:ID"`
		Type_of_source_id         *uint                                  `json:"type_of_source_id" gorm:"type:integer"`
		Type_of_source            model_core.Lookupcode                  `json:"type_of_source" gorm:"foreignKey:Type_of_source_id; references:ID"`
		Fixed_inventory           string                                 `json:",omitempty" gorm:"type:varchar"`
		Inventory_range           string                                 `json:",omitempty" gorm:"type:varchar"`
	}
	Order_details struct {
		Assign_auto_fullfilment         *bool                                  `gorm:"type:boolean"`
		Order_select_source_facility_id *uint                                  `json:"order_select_source_facility_id" gorm:"type:integer"`
		Order_select_source_facility    *shared_pricing_and_location.Locations `json:"order_select_source_facility" gorm:"foreignKey:Order_select_source_facility_id; references:ID"`
	}
	Advanced_channel_configuration struct {
		Prices_include_taxes_id      *uint                 `json:"prices_include_taxes_id" gorm:"type:integer"`
		Prices_include_taxes         model_core.Lookupcode `json:"prices_include_taxes" gorm:"foreignKey:Prices_include_taxes_id; references:ID"`
		Partial_fullfilment_shipping *bool                 `json:"partial_fulfilment_shipping" gorm:"type:boolean"`
		Email_for_shipping_orders    *bool                 `json:"email_for_shipping_orders" gorm:"type:boolean"`
	}
	Payment_mapping struct {
		Payment_method             string `gorm:"type:varchar"`
		Payment_method_name        string `gorm:"type:varchar"`
		Payment_method_description string `gorm:"type:varchar"`
	}
	Inventory_Automation_Details struct {
		Inventory_type_id                 *uint                 `json:"inventory_type_id" gorm:"type:integer"`
		Inventory_type                    model_core.Lookupcode `json:"inventory_type" gorm:"foreignKey:Inventory_type_id; references:ID"`
		Automation_select_source_facility string                `gorm:"type:varchar"`
		Sync_inventory                    *bool                 `gorm:"type:boolean"`
		Sync_price                        *bool                 `gorm:"type:boolean"`
	}
	Pre_sync_configuration struct {
		Pause_inventory_update *bool  `gorm:"type:boolean"`
		Pause_price_update     *bool  `gorm:"type:boolean"`
		Pause_order_update     *bool  `gorm:"type:boolean"`
		Start_date             string `gorm:"type:varchar"`
	}
	Sync_configuration struct {
		Store_url       string `gorm:"type:varchar"`
		Consumer_key    string `gorm:"type:varchar"`
		Store_name      string `gorm:"type:varchar"`
		Consumer_secret string `gorm:"type:varchar"`
	}
	Sync_fetch_data struct {
		Products  *bool `gorm:"type:boolean"`
		Orders    *bool `gorm:"type:boolean"`
		Inventory *bool `gorm:"type:boolean"`
	}
)
