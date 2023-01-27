package shipping_orders_ndr

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
type (
	NDRRequest struct {
		model_core.Model
		ShippingOrderId         uint       `json:"shipping_order_id"`
		Amount                  float64    `json:"amount"`
		DeliveryAttemptLeft     uint       `json:"delivery_attempt_left"`
		NDRLines                []NDRLines `json:"ndr_lines"`
		LastDeliveryAttemptDate time.Time  `json:"last_delivery_attempt_date"`
		FailureReason           string     `json:"failure_reason"`
	}

	NDRLines struct {
		NdrsId                  string    `json:"ndrs_id"`
		DeliveryAttemptMade     uint      `json:"delivery_attempt_made"`
		FailureReason           string    `json:"failure_reason"`
		NdrsRaisedDate          time.Time `json:"ndrs_raised_date"`
		LastDeliveryAttemptDate time.Time `json:"last_delivery_attempt_date"`
	}
)

type (
	CreateNDR struct {
		Created_NDR_Id int
	}
	NDRCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CreateNDR
		}
	} // @name NDRCreateResponse
)

type (
	BulkNDRCreate struct {
		Created bool
	}

	BulkNDRCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkNDRCreate
		}
	} // @ name BulkNDRCreateResponse
)

type (
	GetNDR struct {
		shipping.NDR
	}
	NDRGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data GetNDR
		}
	} // @name NDRGetResponse
)

type (
	GetAllNDR struct {
		shipping.NDR
	}
	GetAllNDRResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []GetAllNDR
		}
	} // @name GetAllNDRResponse
)

type (
	UpdateNDR struct {
		Updated_NDR_Id int
	}
	UpdateNDRResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data UpdateNDR
		}
	} // @name UpdateNDRResponse
)

type (
	DeleteNDR struct {
		Deleted_NDR_Id int
	}
	DeleteNDRResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteNDR
		}
	} // @name DeleteNDRResponse
)

type (
	DeleteNDRLine struct {
		Deleted_NDR_Line_Id int
	}
	DeleteNDRLineResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteNDRLine
		}
	} // @name DeleteNDRLineResponse
)
