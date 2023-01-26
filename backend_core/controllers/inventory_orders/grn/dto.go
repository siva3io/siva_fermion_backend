package grn

import (
	"fermion/backend_core/internal/model/inventory_orders"
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
	CreateGRN struct {
		Created_GRN_Id int
	}
	GRNCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data CreateGRN
		}
	} // @name GRNCreateResponse
)
type (
	GRNRequest struct {
		CreatedByID                 *uint                  `json:"created_by"`
		UpdatedByID                 *uint                  `json:"updated_by"`
		DeletedByID                 *uint                  `json:"deleted_by"`
		GRNNumber                   string                 `json:"grn_number"`
		AutoGenerateGrnNumber       bool                   `json:"auto_generate_grn_number"`
		ReferenceNumber             string                 `json:"reference_number"`
		AutoGenerateReferenceNumber bool                   `json:"auto_generate_reference_number"`
		SourceDocumentTypeId        uint                   `json:"source_document_type_id"`
		SourceDocuments             map[string]interface{} `json:"source_document"`
		CreateScrapOrder            bool                   `json:"create_scrap_order"`
		GRNOrderLines               []GRNOrderLines        `json:"grn_order_lines"`
		WarehouseID                 *uint                  `json:"warehouse_id"`
		StatusId                    uint                   `json:"status_id"`
	}

	GRNOrderLines struct {
		ProductID          uint   `json:"product_id"`
		ProductTemplateId  uint   `json:"product_template_id"`
		LotNumber          string `json:"lot_number"`
		UOMId              uint   `json:"uom_id"`
		OrderedUnits       int    `json:"ordered_units"`
		ReceivedUnits      int    `json:"received_units"`
		PendingUnits       int    `json:"pending_units"`
		QualityCheck       bool   `json:"quality_check"`
		ShelfLocation      string `json:"shelf_location"`
		RejectedQuantities int    `json:"rejected_quantities"`
		ReasonOfRejection  string `json:"reason_of_rejection"`
		PutawayStatus      bool   `json:"putaway_status"`
		PutawayLocation    string `json:"putaway_location"`
	}
)

type (
	GetGRN struct {
	}
	GRNGetResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data GetGRN
		}
	} // @name GRNGetResponse
)

type (
	GetAllGRN struct {
		inventory_orders.GRN
	}
	GetAllGRNResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []GetAllGRN
		}
	} // @name GetAllGRNResponse
)

type (
	UpdateGRN struct {
		Updated_GRN_Id int
	}
	UpdateGRNResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data UpdateGRN
		}
	} // @name UpdateGRNResponse
)

type (
	SearchGRN struct {
		inventory_orders.GRN
	}
	SearchGRNResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data []SearchGRN
		}
	} // @name SearchGRNResponse
)

type SearchKeys struct {
	GRNNumber         string `json:"grn_number,omitempty" query:"grn_number"`
	ReferenceNumberId string `json:"reference_number_id,omitempty" query:"reference_number_id"`
	Status            string `json:"status,omitempty" query:"status"`
	Name              string `json:"name,omitempty" query:"name"`
}

type (
	DeleteGRN struct {
		Deleted_GRN_Id int
	}
	DeleteGRNResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteGRN
		}
	} // @name DeleteGRNResponse
)

type (
	DeleteGRNOrderLine struct {
		Deleted_GRN_Line_Item string
		Grn_Id                int
	}
	DeleteGRNOrderLineResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data DeleteGRNOrderLine
		}
	} // @name DeleteGRNOrderLineResponse
)

type (
	BulkGRNCreate struct {
		Created bool
	}

	BulkGRNCreateResponse struct {
		Body struct {
			Meta response.MetaResponse
			Data BulkGRNCreate
		}
	} // @ name BulkGRNCreateResponse
)
