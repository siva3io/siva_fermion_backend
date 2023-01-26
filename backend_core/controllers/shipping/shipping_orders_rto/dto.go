package shipping_orders_rto

import (
	"time"

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
type RTORequest struct {
	model_core.Model
	ShippingOrderId       uint                     `json:"shipping_order_id"`
	BookingDate           time.Time                `json:"booking_date"`
	ReturnDate            time.Time                `json:"return_date"`
	TotalAmount           float64                  `json:"total_amount"`
	RTOCost               float64                  ` json:"rto_cost"`
	ActualDeliveryDate    time.Time                `json:"actual_delivery_date"`
	EstimatedDeliveryDate time.Time                `json:"estimate_delivery_date"`
	ShippingPartnerId     uint                     `json:"shipping_partner_id"`
	Quantity              int32                    `json:"quantity"`
	RTOStatusId           uint                     `json:"rto_status_id"`
	StatusHistory         []map[string]interface{} `json:"status_history"`
}

type (
	CreateRTO struct {
		Created_RTO_Id int
	}
	RTOCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CreateRTO
		}
	} // @name RTOCreateResponse
)

type (
	BulkRTOCreate struct {
		Created bool
	}

	BulkRTOCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkRTOCreate
		}
	} // @ name BulkRTOCreateResponse
)

type (
	GetRTO struct {
		shipping.RTO
	}
	RTOGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data GetRTO
		}
	} // @name RTOGetResponse
)

type (
	GetAllRTO struct {
		shipping.RTO
	}
	GetAllRTOResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []GetAllRTO
		}
	} // @name GetAllRTOResponse
)

type (
	UpdateRTO struct {
		Updated_RTO_Id int
	}
	UpdateRTOResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data UpdateRTO
		}
	} // @name UpdateRTOResponse
)

type (
	DeleteRTO struct {
		Deleted_RTO_Id int
	}
	DeleteRTOResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteRTO
		}
	} // @name DeleteRTOResponse
)
