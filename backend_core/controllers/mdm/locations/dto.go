package locations

import (
	"time"

	app_core "fermion/backend_core/controllers/cores"
	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/pkg/util/response"

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
type LocationsDTO struct {
	Name               string                 `json:"name"`
	LocationTypeId     uint                   `json:"location_type_id"`
	LocationCode       string                 `json:"location_id"`
	ParentLocationId   *uint                  `json:"parent_id"`
	Address            map[string]interface{} `json:"address"`
	LocationDocs       map[string]interface{} `json:"location_docs"`
	Latitude           float32                `json:"latitude"`
	Longitude          float32                `json:"longitude"`
	ServiceableAreaIds map[string]interface{} `json:"serviceable_area_ids"`
	RelatedLocationID  uint                   `json:"related_location_id"`
	LocationDetails    struct {
		ExternalDetails            map[string]interface{}    `json:"external_details"`
		PaymentMapping             map[string]interface{}    `json:"payment_mapping"`
		LocationInchargerDetails   map[string]interface{}    `json:"location_incharge_details"`
		Storename                  string                    `json:"store_name"`
		CurrencyId                 uint                      `json:"currency_id"`
		PriceListId                *uint                     `json:"price_list_id"`
		LinkedFacilityId           *uint                     `json:"linked_facility_id"`
		AllowBackOrder             bool                      `json:"allow_back_order"`
		OrderTagIds                []app_core.LookupCodesDTO `json:"order_tag_ids"`
		PriceIncludesTaxId         uint                      `json:"price_includes_tax_id"`
		EmailNotification          bool                      `json:"email_notification"`
		PartialFulfilment          bool                      `json:"partial_fulfilment"`
		InventoryTypeId            uint                      `json:"inventory_type_id"`
		SourceFacilityId           uint                      `json:"source_facility_id"`
		ContactDetails             map[string]interface{}    `json:"contact_details"`
		IsScrapLocation            bool                      `json:"is_scrap_location"`
		IsReturnLocation           bool                      `json:"is_return_location"`
		ShippingPartners           map[string]interface{}    `json:"shipping_partners"`
		IntegratedChannels         map[string]interface{}    `json:"integrated_channels"`
		WarehouseStorageManagement map[string]interface{}    `json:"warehouse_storage_management"`
		Zone                       []StorageLocationDTO      `json:"zone"`
		// RacksCapacity              uint                      `json:"racks_capacity"`
		// ShelvesCapacity            uint                      `json:"shelves_capacity"`
		// BinsCapacity               uint                      `json:"bins_capacity"`
		// RacksUomId                 *uint                     `json:"racks_uom_id"`
		// ShelvesUomId               *uint                     `json:"shelves_uom_id"`
		// BinsUomId                  *uint                     `json:"bins_uom_id"`
	} `json:"location_details"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
	Notes        string `json:"notes"`
}

type LocationViewDTO struct {
	Meta response.MetaResponse
	Data struct {
		Id                 uint                    `json:"id"`
		Name               string                  `json:"name"`
		LocationType       app_core.LookupCodesDTO `json:"LocationType"`
		LocationCode       string                  `json:"location_id"`
		Parent             *LocationResponseDTO    `json:"parent_location"`
		ChildLocationIds   []LocationResponseDTO   `json:"child_locations"`
		Address            map[string]interface{}  `json:"address"`
		LocationDocs       map[string]interface{}  `json:"location_docs"`
		Latitude           float32                 `json:"latitude"`
		Longitude          float32                 `json:"longitude"`
		ServiceableAreaIds map[string]interface{}  `json:"serviceable_area_ids"`
		RelatedLocationID  uint                    `json:"related_location_id"`
		LocationDetails    struct {
			LocationDetails struct {
				ExternalDetails            map[string]interface{}    `json:"external_details"`
				PaymentMapping             map[string]interface{}    `json:"payment_mapping"`
				LocationInchargerDetails   map[string]interface{}    `json:"location_incharge_details"`
				Storename                  string                    `json:"store_name"`
				Currency                   Currency                  `json:"currency"`
				PriceList                  PriceListResponseDTO      `json:"price_list"`
				LinkedFacility             IdNameDTO                 `json:"linked_facility"`
				AllowBackOrder             bool                      `json:"allow_back_order"`
				OrderTagIds                []app_core.LookupCodesDTO `json:"order_tag_ids"`
				PriceIncludesTax           app_core.LookupCodesDTO   `json:"price_includes_tax"`
				EmailNotification          bool                      `json:"email_notification"`
				PartialFulfilment          bool                      `json:"partial_fulfilment"`
				InventoryType              app_core.LookupCodesDTO   `json:"inventory_type"`
				SourceFacility             IdNameDTO                 `json:"source_facility"`
				ContactDetails             map[string]interface{}    `json:"contact_details"`
				IsScrapLocation            bool                      `json:"is_scrap_location"`
				IsReturnLocation           bool                      `json:"is_return_location"`
				ShippingPartners           map[string]interface{}    `json:"shipping_partners"`
				IntegratedChannels         map[string]interface{}    `json:"integrated_channels"`
				WarehouseStorageManagement map[string]interface{}    `json:"warehouse_storage_management"`
				Zone                       []StorageLocationDTO      `json:"zone"`
				// RacksCapacity              uint                      `json:"racks_capacity"`
				// ShelvesCapacity            uint                      `json:"shelves_capacity"`
				// BinsCapacity               uint                      `json:"bins_capacity"`
				// RacksUom                   IdNameDTO                 `json:"racks_uom"`
				// ShelvesUom                 IdNameDTO                 `json:"shelves_uom"`
				// BinsUom                    IdNameDTO                 `json:"bins_uom"`
			} `json:"location_details"`
		} `json:"location_details"`
		Email        string `json:"email"`
		MobileNumber string `json:"mobile_number"`
		Notes        string `json:"notes"`
	}
}
type LocationListDTO struct {
	Meta response.MetaResponse
	Data []LocationListResponseDTO
}

// -----------------LocationDTO for model struct-------------------------------
type ModelDto struct {
	ID          uint      `json:"id"`
	CompanyId   uint      `json:"company_id" gorm:"column:company_id"`
	IsEnabled   bool      `json:"is_enabled" gorm:"default:true"`
	IsActive    *bool     `json:"is_active" gorm:"default:true"`
	UpdatedDate time.Time `json:"updated_date" gorm:"autoUpdateTime"`
	CreatedDate time.Time `json:"created_date" gorm:"autoCreateTime"`
}
type Currency struct {
	Id             uint   `json:"id"`
	Name           string `json:"name"`
	CurrencySymbol string `json:"currency_symbol"`
	CurrencyCode   string `json:"currency_code"`
}
type IdNameDTO struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
type PriceListResponseDTO struct {
	Id              uint   `json:"id"`
	Price_list_name string `json:"price_list_name"`
}
type LocationRequestDTO struct {
	CreatedByID        *uint          `json:"created_by"`
	UpdatedByID        *uint          `json:"updated_by"`
	DeletedByID        *uint          `json:"deleted_by"`
	CompanyId          *uint          `json:"company_id"`
	Name               string         `json:"name"`
	LocationTypeId     uint           `json:"location_type_id"`
	LocationCode       string         `json:"location_id"`
	ParentLocationId   *uint          `json:"parent_id"`
	Address            datatypes.JSON `json:"address,omitempty"`
	LocationDocs       datatypes.JSON `json:"location_docs,omitempty"`
	Latitude           float32        `json:"latitude"`
	Longitude          float32        `json:"longitude"`
	ServiceableAreaIds datatypes.JSON `json:"serviceable_area_ids,omitempty"`
	RelatedLocationID  uint           `json:"related_location_id"`
	LocationDetails    datatypes.JSON `json:"location_details,omitempty"`
	Email              string         `json:"email"`
	MobileNumber       string         `json:"mobile_number"`
	Notes              string         `json:"notes"`
	IsEnabled          *bool          `json:"is_enabled" gorm:"default:true"`
	IsActive           *bool          `json:"is_active" gorm:"default:true"`
}
type LocationResponseDTO struct {
	ModelDto
	Name                string                  `json:"name"`
	LocationType        app_core.LookupCodesDTO `json:"location_type"`
	LocationTypeDetails app_core.LookupCodesDTO `json:"LocationType"`
	LocationCode        string                  `json:"location_id"`
	Parent              *LocationResponseDTO    `json:"parent_location"`
	ChildLocationIds    []LocationResponseDTO   `json:"child_locations"`
	Address             datatypes.JSON          `json:"address"`
	LocationDocs        datatypes.JSON          `json:"location_docs"`
	Latitude            float32                 `json:"latitude"`
	Longitude           float32                 `json:"longitude"`
	ServiceableAreaIds  datatypes.JSON          `json:"serviceable_area_ids"`
	RelatedLocationID   uint                    `json:"related_location_id"`
	LocationDetails     datatypes.JSON          `json:"location_details"`
	Email               string                  `json:"email"`
	MobileNumber        string                  `json:"mobile_number"`
	Notes               string                  `json:"notes"`
	IsEnabled           *bool                   `json:"is_enabled" gorm:"default:true"`
	IsActive            *bool                   `json:"is_active" gorm:"default:true"`
}
type LocationListResponseDTO struct {
	ModelDto
	Name               string                  `json:"name"`
	LocationType       app_core.LookupCodesDTO `json:"LocationType"`
	LocationCode       string                  `json:"location_id"`
	LocationDocs       datatypes.JSON          `json:"location_docs"`
	Address            datatypes.JSON          `json:"address"`
	Email              string                  `json:"email"`
	MobileNumber       string                  `json:"mobile_number"`
	Notes              string                  `json:"notes"`
	ServiceableAreaIds datatypes.JSON          `json:"serviceable_area_ids"`
}
type VirtualLocationDetailsResponseDTO struct {
	ExternalDetails          datatypes.JSON `json:"external_details"`
	PaymentMapping           datatypes.JSON `json:"payment_mapping"`
	LocationInchargerDetails datatypes.JSON `json:"location_incharge_details"`
}
type LocalWarehouseDetailsResponseDTO struct {
	ContactDetails             datatypes.JSON       `json:"contact_details"`
	IsScrapLocation            bool                 `json:"is_scrap_location"`
	IsReturnLocation           bool                 `json:"is_return_location"`
	ShippingPartners           datatypes.JSON       `json:"shipping_partners"`
	IntegratedChannels         datatypes.JSON       `json:"integrated_channels"`
	WarehouseStorageManagement datatypes.JSON       `json:"warehouse_storage_management"`
	Zone                       []StorageLocationDTO `json:"zone"`
	// RacksCapacity              uint           `json:"racks_capacity"`
	// ShelvesCapacity            uint           `json:"shelves_capacity"`
	// BinsCapacity               uint           `json:"bins_capacity"`
	// RacksUom                   IdNameDTO      `json:"racks_uom"`
	// ShelvesUom                 IdNameDTO      `json:"shelves_uom"`
	// BinsUom                    IdNameDTO      `json:"bins_uom"`
}

type StorageLocationDTO struct {
	ID                      uint                              `json:"id"`
	CreatedByID             *uint                             `json:"created_by_id"`
	UpdatedByID             *uint                             `json:"updated_by_id"`
	DeletedByID             *uint                             `json:"deleted_by_id"`
	ParentStorageLocationId *uint                             `json:"parent_storage_location_id"`
	ParentStorageLocation   *StorageLocationDTO               `json:"parent_storage_location"`
	LocalWarehouseId        *uint                             `json:"local_warehouse_id"`
	LocalWarehouse          *LocalWarehouseDetailsResponseDTO `json:"local_warehouse"`
	ZoneId                  *uint                             `json:"zone_id"`
	RowId                   *uint                             `json:"row_id"`
	RackId                  *uint                             `json:"rack_id"`
	ShelfId                 *uint                             `json:"shelf_id"`
	Capacity                int64                             `json:"capacity"`
	UomId                   *uint                             `json:"uom_id"`
	Uom                     datatypes.JSON                    `json:"uom"`
	Name                    string                            `json:"name"`
	TypeId                  uint                              `json:"type_id"`
	Type                    model_core.Lookupcode             `json:"type"`
	ChildStorageLocations   []*StorageLocationDTO             `json:"child_storage_locations"`
	StorageQuantity         []*StorageQuantityDTO             `json:"storage_quantity"`
	//Zone                    *StorageLocationDTO               `json:"zone"`
	//Row                     *StorageLocationDTO               `json:"row"`
	//Rack                    *StorageLocationDTO               `json:"rack"`
	//Shelf                   *StorageLocationDTO               `json:"shelf"`
}

type StorageLocationListDTO struct {
	Meta response.MetaResponse
	Data []StorageLocationDTO
}

type StorageQuantityDTO struct {
	ID                uint           `json:"id"`
	StorageLocationId uint           `json:"storage_location_id"`
	ProductVariantId  *uint          `json:"product_variant_id"`
	ProductVariant    datatypes.JSON `json:"product_variant"`
	StorageCount      int64          `json:"stock_count"`
}

type StorageQuantityListDTO struct {
	Meta response.MetaResponse
	Data []StorageQuantityDTO
}

type RetailResponseDTO struct {
	Storename                string                    `json:"store_name"`
	Currency                 Currency                  `json:"currency"`
	PriceList                PriceListResponseDTO      `json:"price_list"`
	LinkedFacility           IdNameDTO                 `json:"linked_facility"`
	AllowBackOrder           bool                      `json:"allow_back_order"`
	OrderTagIds              []app_core.LookupCodesDTO `json:"order_tag_ids"`
	PriceIncludesTax         app_core.LookupCodesDTO   `json:"price_includes_tax"`
	EmailNotification        bool                      `json:"email_notification"`
	PartialFulfilment        bool                      `json:"partial_fulfilment"`
	InventoryType            app_core.LookupCodesDTO   `json:"inventory_type"`
	SourceFacility           IdNameDTO                 `json:"source_facility"`
	PaymentMapping           datatypes.JSON            `json:"payment_mapping"`
	LocationInchargerDetails datatypes.JSON            `json:"location_incharge_details"`
}
type OfficeOrFactoryResponseDTO struct {
	LocationInchargerDetails datatypes.JSON `json:"location_incharge_details"`
}

type ThirdPartyVirtualWarehouseResponseDTO struct {
	model_core.ExternalMapper
}
