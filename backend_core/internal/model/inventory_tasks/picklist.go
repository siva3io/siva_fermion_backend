package inventory_tasks

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/mdm"
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
type PickList struct {
	PickListNumber       string                                `json:"pick_list_number" gorm:"type:varchar(50)"`
	ReferenceNumber      string                                `json:"reference_number" gorm:"type:varchar(50)"`
	SourceDocumentTypeID *uint                                 `json:"source_document_type_id"`
	SourceDocumentType   model_core.Lookupcode                 `json:"source_doc_type" gorm:"foreignkey:SourceDocumentTypeID; references:ID"`
	AssigneeToID         *uint                                 `json:"assignee_to_id" gorm:"type:int"`
	AssigneeTo           mdm.Partner                           `json:"assignee_to" gorm:"foreignkey:AssigneeToID; references:ID"`
	SelectCustomerId     *uint                                 `json:"select_customer_id" gorm:"type:int"`
	SelectCustomer       mdm.Partner                           `json:"select_customer" gorm:"foreignkey:SelectCustomerId; references:ID"`
	PartnerID            *uint                                 `json:"partner_id"`
	Partner              mdm.Partner                           `json:"partner" gorm:"foreignkey:PartnerID; references:ID"`
	WarehouseID          *uint                                 `json:"warehouse_id"`
	Warehouse            shared_pricing_and_location.Locations `json:"warehouse" gorm:"foreignkey:WarehouseID; references:ID"`
	StatusID             uint                                  `json:"status_id"`
	Status               model_core.Lookupcode                 `json:"status"  gorm:"foreignkey:StatusID; references:ID"`
	StatusHistory        datatypes.JSON                        `json:"status_history" gorm:"type:json"`
	PicklistLines        []PickListLines                       `json:"picklist_lines" gorm:"foreignkey:PickList_Id; references:ID; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	InternalNotes        string                                `gorm:"type:varchar(250)" json:"internal_notes"`
	ExternalNotes        string                                `gorm:"type:varchar(250)" json:"external_notes"`
	AttachmentFiles      datatypes.JSON                        `json:"attachment_files" gorm:"type:json"`
	StartDateTime        time.Time                             `json:"start_date_time" gorm:"type:time"`
	EndDateTime          time.Time                             `json:"end_date_time" gorm:"type:time"`
	TotalItemsToPick     uint                                  `gorm:"type:integer" json:"items_to_pick"`
	TotalPickedItems     uint                                  `gorm:"type:integer" json:"total_picked_items"`
	PriceListID          *uint                                 `json:"price_list_id"`
	PriceList            shared_pricing_and_location.Pricing   `json:"price_list" gorm:"foreignkey:PriceListID; references:ID"`
	SourceDocuments      datatypes.JSON                        `json:"source_documents" gorm:"type:json; default:'[]'; not null"`
	model_core.Model
}

type PickListLines struct {
	PickList_Id       uint                 `json:"-"`
	ProductID         *uint                `json:"product_id"`
	Product           *mdm.ProductTemplate `json:"product" gorm:"foreignkey:ProductID; references:ID"`
	ProductVariantID  *uint                `json:"product_variant_id"`
	ProductVariant    *mdm.ProductVariant  `json:"product_variant" gorm:"foreignKey:ProductVariantID; references:ID"`
	SalesDocumentID   datatypes.JSON       `json:"sales_document_id" gorm:"type:json"`
	PartnerID         *uint                `json:"partner_id"`
	Partner           mdm.Partner          `json:"partner" gorm:"foreignkey:PartnerID; references:ID"`
	QuantityOrdered   uint                 `json:"quantity_ordered" gorm:"type:integer"`
	QuantityToPick    uint                 `json:"quantity_to_pick" gorm:"type:integer"`
	QuantityPicked    uint                 `json:"quantity_picked" gorm:"type:integer"`
	RemainingQuantity uint                 `json:"remaining_quantity" gorm:"type:integer"`
	CustomerName      string               `json:"customer_name" gorm:"type:varchar(250)"`
	model_core.Model
}
