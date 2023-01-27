package shipping_orders_wd

import (
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/shipping"
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
type WDRequest struct {
	model_core.Model
	ShippingOrderId        uint                   `json:"shipping_order_id"`
	CourierPartnerFileId   map[string]interface{} `json:"courier_partner_file_id"`
	TransactionId          string                 `json:"transaction_id"`
	InitialWeightTypeId    uint                   `json:"initial_weight_type_id"`
	FinalWeightTypeId      uint                   `json:"final_weight_type_id"`
	InitialWeightTaken     float64                `json:"initial_weight_taken"`
	FinalWeightTaken       float64                `json:"final_weight_taken"`
	FinalAmount            float64                `json:"final_amount"`
	InitialAmount          float64                `json:"initial_amount"`
	DiscrepancyAmount      float64                `json:"discrepancy_amount"`
	WeightDiscrepancyProof map[string]interface{} `json:"Weight_discrepancy_proof"`
}

type (
	CreateWD struct {
		Created_WD_Id int
	}
	WDCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CreateWD
		}
	} // @name WDCreateResponse
)

type (
	BulkWDCreate struct {
		Created bool
	}

	BulkWDCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkWDCreate
		}
	} // @ name BulkWDCreateResponse
)

type (
	GetWD struct {
		shipping.WD
	}
	WDGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data GetWD
		}
	} // @name WDGetResponse
)

type (
	GetAllWD struct {
		shipping.WD
	}
	GetAllWDResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []GetAllWD
		}
	} // @name GetAllWDResponse
)

type (
	UpdateWD struct {
		Updated_WD_Id int
	}
	UpdateWDResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data UpdateWD
		}
	} // @name UpdateWDResponse
)

type (
	DeleteWD struct {
		Deleted_WD_Id int
	}
	DeleteWDResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteWD
		}
	} // @name DeleteWDResponse
)
