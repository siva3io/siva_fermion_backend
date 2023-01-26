package webstores

import (
	model_core "fermion/backend_core/internal/model/core"
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
type (
	Webstore struct {
		Name            string                 `json:"name"`
		Image_options   map[string]interface{} `json:"image_options"`
		Location_served map[string]interface{} `json:"location_served"`
	}
)
type (
	User_Webstore_Link struct {
		Webstore_id                    uint                           `json:"webstore_id"`
		Channel_details                Channel_details                `json:"channel_details"`
		Fullfilment_details            Fullfilment_details            `json:"fullfilment_details"`
		Inventory_details              Inventory_details              `json:"inventory_details"`
		Order_details                  Order_details                  `json:"order_details"`
		Advanced_channel_configuration Advanced_channel_configuration `json:"advanced_channel_configuration"`
		Payment_mapping                Payment_mapping                `json:"payment_mapping"`
		Pre_sync_configuration         Pre_sync_configuration         `json:"pre_sync_configuration"`
		Sync_configuration             Sync_configuration             `json:"sync_configuration"`
		Sync_fetch_data                Sync_fetch_data                `json:"sync_fetch_data"`
		Inventory_Automation_Details   Inventory_Automation_Details   `json:"inventory_Automation_Details"`
	}

	Channel_details struct {
		Price_list_id uint `json:"price_list_id"`
		Currency_id   uint `json:"currency_id"`
		Order_tags_id uint `json:"order_tags_id"`
	}
	Fullfilment_details struct {
		Enable_fba_fbf_to_fullfil_orders bool                   `json:"enable_fba_fbf_to_fullfil_orders"`
		Select_the_3pl_facility          map[string]interface{} `json:"select_the_3pl_facility"`
	}
	Inventory_details struct {
		Select_source_facility_id uint   `json:"select_source_facility_id"`
		Type_of_source_id         uint   `json:"type_of_source_id"`
		Fixed_inventory           string ` json:"fixed_inventory"`
		Inventory_range           string ` json:"inventory_range"`
	}
	Order_details struct {
		Assign_auto_fullfilment         bool `json:"assign_auto_fullfilment"`
		Order_select_source_facility_id uint `json:"order_select_source_facility_id"`
	}
	Advanced_channel_configuration struct {
		Prices_include_taxes_id      uint `json:"prices_include_taxes_id"`
		Partial_fullfilment_shipping bool `json:"partial_fulfilment_shipping"`
		Email_for_shipping_orders    bool `json:"email_for_shipping_orders"`
	}
	Payment_mapping struct {
		Payment_method             string `json:"payment_method"`
		Payment_method_name        string `json:"payment_method_name"`
		Payment_method_description string `json:"payment_method_description"`
	}
	Inventory_Automation_Details struct {
		Inventory_type_id                 uint   `json:"inventory_type_id"`
		Automation_select_source_facility string `json:"automation_select_source_facility"`
		Sync_inventory                    bool   `json:"sync_inventory"`
		Sync_price                        bool   `json:"sync_price"`
	}
	Pre_sync_configuration struct {
		Pause_inventory_update bool   `json:"pause_inventory_update"`
		Pause_price_update     bool   `json:"pause_price_update"`
		Pause_order_update     bool   `json:"pause_order_update"`
		Start_date             string `json:"start_date"`
	}
	Sync_configuration struct {
		Store_url       string `json:"store_url"`
		Consumer_key    string `json:"consumer_key"`
		Store_name      string `json:"store_name"`
		Consumer_secret string `json:"consumer_secret"`
	}
	Sync_fetch_data struct {
		Products  bool `json:"products"`
		Orders    bool `json:"orders"`
		Inventory bool `json:"inventory"`
	}
)

// Create Webstore response
type (
	WebstoredetailCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name WebstoredetailCreateResponse
)

// Update Webstore response
type (
	WebstoredetailUpdateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name WebstoredetailUpdateResponse
)

// Get  Webstore response
type (
	WebstoredetailGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name WebstoredetailGetResponse
)

type (
	WebstoredetailGetAllResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []User_Webstore_Link
		}
	} //@ name WebstoredetailGetAllResponse
)

type (
	AvailableWebstoresResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []Webstore
		}
	} //@ name AvailableWebstoresResponse
)

// Delete Webstore response
type (
	WebstoredetailDeleteResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data interface{}
		}
	} //@ name WebstoredetailDeleteResponse
)

type ThirdPartyWebstoreResponseDTO struct {
	model_core.ExternalMapper
}
