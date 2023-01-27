package products

import (
	"time"

	core_dto "fermion/backend_core/controllers/cores"
	mdm_UOM "fermion/backend_core/controllers/mdm/uom"
	"fermion/backend_core/internal/model/mdm"

	"github.com/lib/pq"
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

//============================Template and Variant DTO========================================================

type ProductSelectedAttributesAndValuesDTO struct {
	ID              uint                                    `json:"id"`
	Name            string                                  `json:"name"`
	AttributeValues []ProductSelectedAttributesAndValuesDTO `json:"attribute_values,omitempty"`
}
type Ids struct {
	ID uint `json:"id"`
}
type CreateProductVariantDTO struct {
	ID                       *uint                    `json:"id"`
	CreatedByID              *uint                    `json:"created_by"`
	UpdatedByID              *uint                    `json:"updated_by"`
	DeletedByID              *uint                    `json:"deleted_by"`
	CompanyId                uint                     `json:"company_id"`
	ParentSKUId              string                   `json:"parent_sku_id"`
	SkuId                    string                   `json:"sku_id,omitempty"`
	ProductName              string                   `json:"product_name"`
	AttributeKeyValuesId     pq.Int64Array            `json:"attribute_key_values_id"`
	ImageOptions             datatypes.JSON           `json:"image_options,omitempty"`
	VariantTypeId            *uint                    `json:"variant_type_id"`
	Barcode                  string                   `json:"barcode"`
	StandardProductTypes     []Ids                    `json:"standard_product_types,omitempty"`
	StandardProductTypeId    string                   `json:"standard_product_type_id"`
	ConditionID              *uint                    `json:"condition_id"`
	ProductDimensions        datatypes.JSON           `json:"product_dimensions,omitempty"`
	CategoryID               *uint                    `json:"category_id"`
	LeafCategoryID           *uint                    `json:"leaf_category_id"`
	Description              datatypes.JSON           `json:"description,omitempty"`
	KeywordIds               datatypes.JSON           `json:"keyword_ids,omitempty"`
	PackageDimensions        datatypes.JSON           `json:"package_dimensions,omitempty"`
	PackageMaterialOptions   datatypes.JSON           `json:"package_material_options,omitempty"`
	VendorPriceListIds       pq.Int64Array            `json:"vendor_price_list_ids"`
	PriceListDetails         datatypes.JSON           `json:"price_list_details,omitempty"`
	ShippingOptions          datatypes.JSON           `json:"shipping_options,omitempty"`
	CostDetails              datatypes.JSON           `json:"cost_details,omitempty"`
	ForecastingOptions       datatypes.JSON           `json:"forecasting_options,omitempty"`
	ProductPricingDetails    PricingDetailsDTO        `json:"product_pricing_details,omitempty"`
	ValidationInfo           string                   `json:"validation_info"`
	InventoryDetails         []Ids                    `json:"inventory_details"`
	PackageTemplateOptions   datatypes.JSON           `json:"package_template_options,omitempty"`
	StatusId                 *uint                    `json:"status_id"`
	ProductTemplateId        uint                     `json:"product_template_id"`
	IsEnabled                *bool                    `json:"is_enabled" gorm:"default:true"`
	IsActive                 *bool                    `json:"is_active" gorm:"default:true"`
	AppId                    *uint                    `json:"app_id"`
	TemplateSku              string                   `json:"template_sku"`
	ChannelCode              string                   `json:"channel_code"`
	ProductVariantSku        string                   `json:"product_variant_sku"`
	ExternalId               string                   `json:"external_id"`
	ExternalDetails          datatypes.JSON           `json:"external_details,omitempty"`
	Location                 Location                 `json:"location"`
	DeliverySlaDetails       DeliverySlaDetails       `json:"delivery_sla_details"`
	InventoryDetail          InventoryDetails         `json:"inventory_detail"`
	ProductCancellationTerms ProductCancellationTerms `json:"product_cancellation_terms"`
	ProductReturnTerms       ProductReturnTerms       `json:"product_return_terms"`
	ProductReplacementTerms  ProductReplacementTerms  `json:"product_replacement_terms"`
	FoodItemDetails          FoodItemDetails          `json:"food_item_details"`
	ManufacturerDetails      ManufacturerDetails      `json:"manufacturer_details"`
	ProductCriticalDetails   ProductCriticalDetails   `json:"product_critical_details"`
	AdditivesInformation     string                   `json:"additives_information"`
	ShortDescription         string                   `json:"short_description"`
	DomainId                 *uint                    `json:"domain_id"`
	Domain                   core_dto.LookupCodesDTO  `json:"domain" gorm:"foreignkey:DomainId;references:ID"`
	OfferDetails             Offers                   `json:"offer_details,omitempty"`

	RatingAverage float64 `json:"rating_average"`
	RatingCount   uint    `json:"rating_count"`

	Channel string `json:"channel"`
}
type PricingDetailsDTO struct {
	PaymentMethodId *uint                   `json:"payment_method_id"`
	PaymentMethod   core_dto.LookupCodesDTO `json:"payment_method"`
	DeclaredPrice   float64                 `json:"declared_price"`
	ID              uint                    `json:"id"`
	SalesPrice      float64                 `json:"sales_price"`
	CostPrice       float64                 `json:"cost_price"`
	MRP             float64                 `json:"mrp" `
	TaxOptions      float64                 `json:"tax_options"`
	CurrencyId      *uint                   `json:"currency_id"`
	Tax             bool                    `json:"tax"`
	Shipping        bool                    `json:"shipping"`
}
type PricingDetailsResponseDTO struct {
	ID              uint                    `json:"id"`
	SalesPrice      float64                 `json:"sales_price"`
	CostPrice       float64                 `json:"cost_price"`
	MRP             float64                 `json:"mrp" `
	TaxOptions      float64                 `json:"tax_options"`
	CurrencyId      *uint                   `json:"currency_id"`
	Currency        core_dto.CurrencyDTO    `json:"currency"`
	Tax             bool                    `json:"tax"`
	Shipping        bool                    `json:"shipping"`
	DeclaredPrice   float64                 `json:"declared_price"`
	PaymentMethodId *uint                   `json:"payment_method_id"`
	PaymentMethod   core_dto.LookupCodesDTO `json:"payment_method"`
}
type CreateProductTemplatePayload struct {
	ID                             *uint  `json:"id"`
	CreatedByID                    *uint  `json:"created_by"`
	UpdatedByID                    *uint  `json:"updated_by"`
	DeletedByID                    *uint  `json:"deleted_by"`
	CompanyId                      uint   `json:"company_id"`
	Name                           string `json:"product_name"`
	BrandID                        *uint  `json:"brand_id"`
	SKUCode                        string `json:"sku_code"`
	HSNCODE                        string `json:"hsn_code"`
	ProductConditionID             *uint  `json:"product_condition_id"`
	ProductTypeID                  *uint  `json:"product_type_id"`
	ProductProcurementTreatmentIds []Ids  `json:"product_procurement_treatment_ids" `
	//EstimatedDeliveryTimeID        *uint                                   `json:"estimated_delivery_time_id"`
	StockTreatmentIds      []Ids                                   `json:"stock_treatment_ids"`
	InventoryTrackingID    *uint                                   `json:"inventory_tracking_id"`
	UomID                  *uint                                   `json:"uom_id"`
	SecondaryUom           datatypes.JSON                          `json:"secondary_uom,omitempty"`
	ImageOptions           datatypes.JSON                          `json:"image_options,omitempty"`
	PrimaryCategoryID      *uint                                   `json:"primary_category_id"`
	SecondaryCategoryID    *uint                                   `json:"secondary_category_id"`
	Description            datatypes.JSON                          `json:"description,omitempty"`
	AttributeKeyValues     []ProductSelectedAttributesAndValuesDTO `json:"attribute_key_values"`
	ProductVariantIds      []CreateProductVariantDTO               `json:"product_variant_ids"`
	ProductPricingDetails  PricingDetailsDTO                       `json:"product_pricing_details"`
	IsEnabled              *bool                                   `json:"is_enabled" gorm:"default:true"`
	IsActive               *bool                                   `json:"is_active" gorm:"default:true"`
	AppId                  *uint                                   `json:"app_id"`
	VendorPriceListIds     pq.Int64Array                           `json:"vendor_price_list_ids"`
	PriceListDetails       datatypes.JSON                          `json:"price_list_details"`
	ShippingOptions        datatypes.JSON                          `json:"shipping_options,omitempty"`
	PackageMaterialOptions datatypes.JSON                          `json:"package_material_options,omitempty"`
	PackageDimensions      datatypes.JSON                          `json:"package_dimensions,omitempty"`
	AdditivesInformation   string                                  `json:"additives_information"`
	ShortDescription       string                                  `json:"short_description"`
	DomainId               *uint                                   `json:"domain_id"`
	Domain                 core_dto.LookupCodesDTO                 `json:"domain"`

	Location                     Location                 `json:"location"`
	DeliverySlaDetails           DeliverySlaDetails       `json:"delivery_sla_details"`
	InventoryDetail              InventoryDetails         `json:"inventory_detail"`
	ProductCancellationTermsFlag bool                     `json:"product_cancellation_terms_flag"`
	ProductCancellationTerms     ProductCancellationTerms `json:"product_cancellation_terms"`
	ProductReturnTermsFlag       bool                     `json:"product_return_terms_flag"`
	ProductReturnTerms           ProductReturnTerms       `json:"product_return_terms"`
	ProductReplacementTermsFlag  bool                     `json:"product_replacement_terms_flag"`
	ProductReplacementTerms      ProductReplacementTerms  `json:"product_replacement_terms"`
	FoodItemDetails              FoodItemDetails          `json:"food_item_details"`
	ManufacturerDetails          ManufacturerDetails      `json:"manufacturer_details"`
	ProductCriticalDetails       ProductCriticalDetails   `json:"product_critical_details"`

	TemplateOptions datatypes.JSON `json:"template_options,omitempty"`
	StatusId        *uint          `json:"status_id"`
	Related         datatypes.JSON `json:"related,omitempty"`
	Recommended     datatypes.JSON `json:"recommended,omitempty"`
	OtherStatuses   datatypes.JSON `json:"other_statuses,omitempty"`

	TemplateSku     string         `json:"template_sku"`
	ChannelCode     string         `json:"channel_code"`
	ExternalId      string         `json:"external_id"`
	ExternalDetails datatypes.JSON `json:"external_details,omitempty"`

	Channel string `json:"channel"`
}
type Location struct {
	Pincode      string             `json:"pincode"`
	City         string             `json:"city"`
	StateId      *uint              `json:"state_id"`
	State        *core_dto.StateDTO `json:"state"`
	AddressLine1 string             `json:"address_line1"`
	AddressLine2 string             `json:"address_line2"`
	AddressLine3 string             `json:"address_line3"`
}

type DeliverySlaDetails struct {
	EstimatedDeliveryTimeId *uint                   `json:"estimated_delivery_time_id"`
	EstimatedDeliveryTime   core_dto.LookupCodesDTO `json:"estimated_delivery_time"`
}

type InventoryDetails struct {
	MaximumQuantity   uint `json:"maximum_quantity"`
	AvailableQuantity uint `json:"available_quantity"`
}

type ProductCriticalDetails struct {
	UnitPerBox         uint                    `json:"unit_per_box"`
	TimeToShipId       *uint                   `json:"time_to_ship_id"`
	TimeToShip         core_dto.LookupCodesDTO `json:"time_to_ship"`
	CustomerCareNumber string                  `json:"customer_care_number"`
}

type ManufacturerDetails struct {
	ManufacturerName    string    `json:"manufacturer_name"`
	ManufacturerAddress string    `json:"manufacturer_address"`
	CommodityName       string    `json:"commodity_name"`
	NetQuantity         uint      `json:"net_quantity"`
	ManufacturerTime    time.Time `json:"manufacturer_time"`
	ManufacturerDate    time.Time `json:"manufacturer_date"`
}

type ProductCancellationTerms struct {
	RefundEligible  *bool                   `json:"refund_eligible"`
	CancelTimeId    *uint                   `json:"cancel_time_id"`
	CancelTime      core_dto.LookupCodesDTO `json:"cancel_time"`
	CancellationFee uint                    `json:"cancellation_fee"`
}
type ProductReturnTerms struct {
	RefundEligibleReturn   *bool                   `json:"refund_eligible_return"`
	FulfillmentManagedById *uint                   `json:"fulfillment_managed_by_id"`
	FulfillmentManagedBy   core_dto.LookupCodesDTO `json:"fulfillment_managed_by"`
	ReturnWithinId         *uint                   `json:"return_within_id"`
	ReturnWithin           core_dto.LookupCodesDTO `json:"return_within"`
}
type ProductReplacementTerms struct {
	ReplacementEligible *bool                   `json:"replacement_eligible"`
	ReplacementWithinId *uint                   `json:"replacement_within_id"`
	ReplacementWithin   core_dto.LookupCodesDTO `json:"replacement_within"`
}
type FoodItemDetails struct {
	FoodTypeId              *uint                   `json:"food_type_id"`
	FoodType                core_dto.LookupCodesDTO `json:"food_type"`
	FSSAILicenceNumber      string                  `json:"fssai_licence_number"`
	OtherFSSAILicenceNumber string                  `json:"other_fssai_licence_number"`
	ImportersFSSAINumber    string                  `json:"importers_fssai_number"`
	TimeToLife              string                  `json:"time_to_life"`
	IngredientsInfo         string                  `json:"ingredients_info"`
	NutritionalInfo         string                  `json:"nutritional_info"`
}
type TemplateSelectedAttributeValuesResponse struct {
	AttributeValueId uint   `json:"attribute_value_id"`
	Name             string `json:"name"`
	AppId            *uint  `json:"app_id"`
}
type TemplateSelectedAttributesResponse struct {
	AttributeId     uint                                      `json:"attribute_id"`
	Name            string                                    `json:"name"`
	AttributeValues []TemplateSelectedAttributeValuesResponse `json:"attribute_values,omitempty"`
	AppId           *uint                                     `json:"app_id"`
}
type TemplateReponseDTO struct {
	ID                 *uint                                `json:"id"`
	CreatedByID        *uint                                `json:"created_by"`
	UpdatedByID        *uint                                `json:"updated_by"`
	CreatedDate        time.Time                            `json:"created_date"`
	UpdatedDate        time.Time                            `json:"updated_date"`
	CompanyId          uint                                 `json:"company_id"`
	Name               string                               `json:"product_name"`
	BrandID            *uint                                `json:"brand_id"`
	Brand              BrandRequestAndResponseDTO           `json:"brand"`
	AttributeValues    []TemplateSelectedAttributesResponse `json:"attribute_values"`
	SkuCode            string                               `json:"sku_code"`
	HSNCODE            string                               `json:"hsn_code"`
	HSNCodesData       mdm.HSNCodesData                     `json:"hsn_codes_data"`
	ProductConditionID *uint                                `json:"product_condition_id"`
	ProductCondition   core_dto.LookupCodesDTO              `json:"product_condition"`
	ProductTypeID      *uint                                `json:"product_type_id"`
	ProductType        core_dto.LookupCodesDTO              `json:"product_type"`
	//EstimatedDeliveryTimeID        *uint                                `json:"estimated_delivery_time_id"`
	//EstimatedDeliveryTime          core_dto.LookupCodesDTO              `json:"estimated_delivery_time"`
	ProductProcurementTreatmentIds datatypes.JSON            `json:"product_procurement_treatment_ids"`
	StockTreatmentIds              datatypes.JSON            `json:"stock_treatment_ids"`
	InventoryTrackingID            *uint                     `json:"inventory_tracking_id"`
	InventoryTracking              core_dto.LookupCodesDTO   `json:"inventory_tracking"`
	UomID                          *uint                     `json:"uom_id"`
	Uom                            mdm_UOM.UomResponseDTO    `json:"uom"`
	SecondaryUom                   datatypes.JSON            `json:"secondary_uom"`
	ImageOptions                   datatypes.JSON            `json:"image_options"`
	PrimaryCategoryID              *uint                     `json:"primary_category_id"`
	PrimaryCategory                CategoryResponseDTO       `json:"primary_category"`
	SecondaryCategoryID            *uint                     `json:"secondary_category_id"`
	SecondaryCategory              CategoryResponseDTO       `json:"secondary_category"`
	Description                    datatypes.JSON            `json:"description"`
	ProductVariantIds              []VariantResponseDTO      `json:"product_variant_ids"`
	ProductPricingDetails          PricingDetailsResponseDTO `json:"product_pricing_details"`
	VendorPriceListIds             pq.Int64Array             `json:"vendor_price_list_ids"`
	PriceListDetails               []map[string]interface{}  `json:"price_list_details"`
	ShippingOptions                datatypes.JSON            `json:"shipping_options"`
	PackageMaterialOptions         datatypes.JSON            `json:"package_material_options"`
	PackageDimensions              datatypes.JSON            `json:"package_dimensions"`
	StatusId                       *uint                     `json:"status_id"`
	Status                         core_dto.LookupCodesDTO   `json:"status"`
	AdditivesInformation           string                    `json:"additives_information"`
	ShortDescription               string                    `json:"short_description"`
	DomainId                       *uint                     `json:"domain_id"`
	Domain                         core_dto.LookupCodesDTO   `json:"domain"`

	Location                 Location                 `json:"location"`
	DeliverySlaDetails       DeliverySlaDetails       `json:"delivery_sla_details"`
	InventoryDetail          InventoryDetails         `json:"inventory_detail"`
	ProductCancellationTerms ProductCancellationTerms `json:"product_cancellation_terms"`
	ProductReturnTerms       ProductReturnTerms       `json:"product_return_terms"`
	ProductReplacementTerms  ProductReplacementTerms  `json:"product_replacement_terms"`
	FoodItemDetails          FoodItemDetails          `json:"food_item_details"`
	ManufacturerDetails      ManufacturerDetails      `json:"manufacturer_details"`
	ProductCriticalDetails   ProductCriticalDetails   `json:"product_critical_details"`

	TemplateOptions datatypes.JSON `json:"template_options"`
	Related         datatypes.JSON `json:"related"`
	Recommended     datatypes.JSON `json:"recommended"`
	OtherStatuses   datatypes.JSON `json:"other_statuses"`
	Channel         string         `json:"channel"`
}
type VariantResponseDTO struct {
	ID                     *uint     `json:"id"`
	CreatedByID            *uint     `json:"created_by"`
	UpdatedByID            *uint     `json:"updated_by"`
	CreatedDate            time.Time `json:"created_date"`
	CreatedBy              core_dto.UserResponseDTO
	UpdatedDate            time.Time                                `json:"updated_date"`
	CompanyId              uint                                     `json:"company_id"`
	SerialNumber           string                                   `json:"serial_number"`
	ProductTemplateId      uint                                     `json:"product_template_id"`
	ParentSkuId            string                                   `json:"parent_sku_id"`
	SkuId                  string                                   `json:"sku_id"`
	ProductName            string                                   `json:"product_name"`
	AttributeKeyValuesId   pq.Int64Array                            `json:"attribute_key_values_id"`
	AttributeValues        []map[string]interface{}                 `json:"attribute_values"`
	ImageOptions           datatypes.JSON                           `json:"image_options"`
	VariantTypeId          *uint                                    `json:"variant_type_id"`
	VariantType            core_dto.LookupCodesDTO                  `json:"variant_type"`
	Barcode                string                                   `json:"barcode"`
	StandardProductTypes   datatypes.JSON                           `json:"standard_product_types"`
	StandardProductTypeId  string                                   `json:"standard_product_type_id"`
	ConditionID            *uint                                    `json:"condition_id"`
	Condition              core_dto.LookupCodesDTO                  `json:"condition"`
	CategoryID             *uint                                    `json:"category_id"`
	Category               CategoryResponseDTO                      `json:"category"`
	LeafCategoryID         *uint                                    `json:"leaf_category_id"`
	LeafCategory           CategoryResponseDTO                      `json:"leaf_category"`
	Description            datatypes.JSON                           `json:"description"`
	ProductDimensions      datatypes.JSON                           `json:"product_dimensions"`
	PackageDimensions      datatypes.JSON                           `json:"package_dimensions"`
	PackageMaterialOptions datatypes.JSON                           `json:"package_material_options"`
	VendorPriceListIds     pq.Int64Array                            `json:"vendor_price_list_ids"`
	PriceListDetails       []map[string]interface{}                 `json:"price_list_details"`
	ShippingOptions        datatypes.JSON                           `json:"shipping_options"`
	ProductPricingDetails  PricingDetailsResponseDTO                `json:"product_pricing_details"`
	InventoryDetails       []DecentralizedBasicInventoryResponseDTO `json:"inventory_details"`
	StatusId               *uint                                    `json:"status_id"`
	Status                 core_dto.LookupCodesDTO                  `json:"status"`
	KeywordIds             datatypes.JSON                           `json:"keyword_ids" gorm:"type:json;default:'{}'"`
	CostDetails            datatypes.JSON                           `json:"cost_details" gorm:"type:json;default:'{}'"`
	ForecastingOptions     datatypes.JSON                           `json:"forecasting_options" gorm:"type:json;default:'{}'"`
	ValidationInfo         string                                   `json:"validation_info"`
	PackageTemplateOptions datatypes.JSON                           `json:"package_template_options" gorm:"type:json;default:'{}'"`
	StoreDetails           core_dto.CompanyDetails                  `json:"store_details"`
	AdditivesInformation   string                                   `json:"additives_information"`
	ShortDescription       string                                   `json:"short_description"`
	DomainId               *uint                                    `json:"domain_id"`
	Domain                 core_dto.LookupCodesDTO                  `json:"domain" gorm:"foreignkey:DomainId;references:ID"`

	Location                 Location                 `json:"location"`
	DeliverySlaDetails       DeliverySlaDetails       `json:"delivery_sla_details"`
	InventoryDetail          InventoryDetails         `json:"inventory_detail"`
	ProductCancellationTerms ProductCancellationTerms `json:"product_cancellation_terms"`
	ProductReturnTerms       ProductReturnTerms       `json:"product_return_terms"`
	ProductReplacementTerms  ProductReplacementTerms  `json:"product_replacement_terms"`
	FoodItemDetails          FoodItemDetails          `json:"food_item_details"`
	ProductCriticalDetails   ProductCriticalDetails   `json:"product_critical_details"`
	ManufacturerDetails      ManufacturerDetails      `json:"manufacturer_details"`
	OfferDetails             Offers                   `json:"offer_details,omitempty"`

	BrandID      *uint                      `json:"brand_id"`
	Brand        BrandRequestAndResponseDTO `json:"brand"`
	HSNCODE      string                     `json:"hsn_code"`
	HSNCodesData mdm.HSNCodesData           `json:"hsn_codes_data"`

	Channel       string  `json:"channel"`
	RatingAverage float64 `json:"rating_average"`
	RatingCount   uint    `json:"rating_count"`
}

type DecentralizedBasicInventoryResponseDTO struct {
	Id                 uint           `json:"id"`
	ChannelCode        string         `json:"channel_code"`
	ProductVariantId   uint           `json:"product_variant_id"`
	ProductDetails     datatypes.JSON `json:"product_details"`
	PhysicalLocationId uint           `json:"physical_location_id"`
	PhysicalLocation   datatypes.JSON `json:"physical_location"`
	Quantity           int64          `json:"quantity"`
	OnHandQuantity     int64          `json:"on_hand_quantity"`
	PlannedInQuantity  int64          `json:"planned_in_quantity"`
	PlannedOutQuantity int64          `json:"planned_out_quantity"`
	BinTag             string         `json:"bin_tag"`
	AvailableStock     int64          `json:"available_stock"`
	StockExpected      int64          `json:"stock_expected"`
	CommitedStock      int64          `json:"commited_stock"`
	NoOfUnits          int64          `json:"no_of_units"`
	CreatedByID        *uint          `json:"created_by"`
	UpdatedByID        *uint          `json:"updated_by"`
	IsEnabled          *bool          `json:"is_enabled" gorm:"default:true"`
	IsActive           *bool          `json:"is_active" gorm:"default:true"`
}

// ================================Brands, Category, Base Attribute, Base Attribute Value , Selected Attribute======================================
// ====================================Selected Attribute Value, Bundles DTO's============================================================================
type BrandRequestAndResponseDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"brand_name"`
	Description string `json:"description"`
	IsEnabled   *bool  `json:"is_enabled"`
	IsActive    *bool  `json:"is_active"`
	CreatedByID *uint  `json:"created_by"`
	UpdatedByID *uint  `json:"updated_by"`
}
type CategoryResponseDTO struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	ParentCategoryID     *uint  `json:"parent_category_id"`
	ExternalID           *uint  `json:"external_id"`
	AdditivesInformation string `json:"additives_information"`
	ShortDescription     string `json:"short_description"`
	DomainId             *uint  `json:"domain_id"`
	CategoryCode         string `json:"category_code"`
}
type CategoryAndSubcategoryRequestAndResponseDTO struct {
	ID                   uint                                          `json:"id"`
	Name                 string                                        `json:"name"`
	ParentCategoryID     *uint                                         `json:"parent_category_id"`
	ChildCategoryIds     []CategoryAndSubcategoryRequestAndResponseDTO `json:"child_ids,omitempty"`
	AppId                *uint                                         `json:"app_id"`
	ExternalID           *uint                                         `json:"external_id"`
	IsEnabled            *bool                                         `json:"is_enabled"`
	IsActive             *bool                                         `json:"is_active"`
	CreatedByID          *uint                                         `json:"created_by"`
	UpdatedByID          *uint                                         `json:"updated_by"`
	AdditivesInformation string                                        `json:"additives_information"`
	ShortDescription     string                                        `json:"short_description"`
	DomainId             *uint                                         `json:"domain_id"`
	CategoryCode         string                                        `json:"category_code"`
}
type ProductBaseAttributesRequestAndResponseDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	IsEnabled   *bool  `json:"is_enabled"`
	IsActive    *bool  `json:"is_active"`
	CreatedByID *uint  `json:"created_by"`
	UpdatedByID *uint  `json:"updated_by"`
}
type ProductBaseAttributesValuesRequestAndResponseDTO struct {
	ID              uint   `json:"id"`
	BaseAttributeId uint   `json:"base_attribute_id"`
	Name            string `json:"name"`
	IsEnabled       *bool  `json:"is_enabled"`
	IsActive        *bool  `json:"is_active"`
	CreatedByID     *uint  `json:"created_by"`
	UpdatedByID     *uint  `json:"updated_by"`
}
type ProductSelectedAttributesRequestAndREsponseDTO struct {
	ID          uint   `json:"id"`
	TemplateId  uint   `json:"template_id"`
	AttributeId uint   `json:"attribute_id"`
	Name        string `json:"name"`
	IsEnabled   *bool  `json:"is_enabled"`
	IsActive    *bool  `json:"is_active"`
	CreatedByID *uint  `json:"created_by"`
	UpdatedByID *uint  `json:"updated_by"`
}
type ProductSelectedAttributesValuesRequestAndResponseeDTO struct {
	ID               uint   `json:"id"`
	TemplateId       uint   `json:"template_id"`
	AttributeId      uint   `json:"attribute_id"`
	AttributeValueId uint   `json:"attribute_value_id"`
	Name             string `json:"name"`
	IsEnabled        *bool  `json:"is_enabled"`
	IsActive         *bool  `json:"is_active"`
	CreatedByID      *uint  `json:"created_by"`
	UpdatedByID      *uint  `json:"updated_by"`
}
type CreateBundlePayload struct {
	CreatedByID    *uint             `json:"created_by"`
	UpdatedByID    *uint             `json:"updated_by"`
	DeletedByID    *uint             `json:"deleted_by"`
	BundleCode     string            `json:"bundle_id"`
	BundleName     string            `json:"bundle_name"`
	Instructions   string            `json:"instructions"`
	Description    datatypes.JSON    `json:"description"`
	ImageOptions   datatypes.JSON    `json:"image_options"`
	Products       []BundleLineItems `json:"products"`
	SalesPrice     float64           `json:"sales_price"`
	CostPrice      float64           `json:"cost_price"`
	MRP            float64           `json:"mrp" `
	TaxOptions     float64           `json:"tax_options"`
	CurrencyId     *uint             `json:"currency_id"`
	SellingPrice   float64           `json:"selling_price"`
	Tax            bool              `json:"tax"`
	Shipping       bool              `json:"shipping"`
	PackageDetails datatypes.JSON    `json:"package_details"`
	StatusId       *uint             `json:"status_id"`
	AppId          *uint             `json:"app_id"`
}
type BundleLineItems struct {
	BundleId         uint  `json:"bundle_id"`
	ProductVariantId uint  `json:"product_variant_id"`
	Quantity         int   `json:"quantity"`
	AppId            *uint `json:"app_id"`
}
type BudleResponseDTO struct {
	ID           *uint                     `json:"id"`
	CreatedByID  *uint                     `json:"created_by"`
	UpdatedByID  *uint                     `json:"updated_by"`
	CreatedDate  time.Time                 `json:"created_date"`
	UpdatedDate  time.Time                 `json:"updated_date"`
	BundleCode   string                    `json:"bundle_id"`
	BundleName   string                    `json:"bundle_name"`
	Instructions string                    `json:"instructions"`
	Description  datatypes.JSON            `json:"description"`
	ImageOptions datatypes.JSON            `json:"image_options"`
	Products     []BundleLineItemsResponse `json:"products"`
	SalesPrice   float64                   `json:"sales_price"`
	CostPrice    float64                   `json:"cost_price"`
	MRP          float64                   `json:"mrp" `
	TaxOptions   float64                   `json:"tax_options"`
	CurrencyId   *uint                     `json:"currency_id"`
	Currency     core_dto.CurrencyDTO      `json:"currency"`
	// SellingPrice   float64                   `json:"selling_price"`
	Tax            bool                    `json:"tax"`
	Shipping       bool                    `json:"shipping"`
	PackageDetails datatypes.JSON          `json:"package_details"`
	StatusId       *uint                   `json:"status_id"`
	Status         core_dto.LookupCodesDTO `json:"status"`
	AppId          *uint                   `json:"app_id"`
}
type BundleLineItemsResponse struct {
	BundleId         uint                       `json:"bundle_id"`
	ProductVariantId uint                       `json:"product_variant_id"`
	ProductVariant   VariantDetailsForBundleDTO `json:"product_variant"`
	Quantity         int                        `json:"quantity"`
	AppId            *uint                      `json:"app_id"`
}
type VariantDetailsForBundleDTO struct {
	ParentSkuId string `json:"parent_sku_id"`
	SkuId       string `json:"sku_id"`
	ProductName string `json:"product_name"`
	AppId       *uint  `json:"app_id"`
}

// ===========================Product to channel Map DTO's==================================================================
type QueryParamsDTO struct {
	Filters  string `json:"filters,omitempty" query:"filters"`
	Per_page int    `json:"per_page,omitempty" query:"per_page"`
	Page_no  int    `json:"page_no,omitempty" query:"page_no"`
	Sort     string `json:"sort,omitempty" query:"sort"`

	Inventory       bool   `json:"inventory,omitempty" query:"inventory"`
	Pricing         bool   `json:"pricing,omitempty" query:"pricing"`
	ProductTemplate bool   `json:"product_template,omitempty" query:"product_template"`
	ProductVariant  bool   `json:"product_variant,omitempty" query:"product_variant"`
	Ids             string `json:"ids,omitempty" query:"ids"`
	Condition       string `json:"condition,omitempty" query:"condition"`
	ColumnName      string `json:"column_name,omitempty" query:"column_name"`
}
type ChannelMapResponseDTO struct {
	ID                     *uint                   `json:"id"`
	ChannelCode            string                  `json:"channel_code"`
	ProductTemplateId      uint                    `json:"product_template_id"`
	TemplateSku            string                  `json:"template_sku"`
	ProductVariantId       uint                    `json:"product_variant_id"`
	ProductVariantSku      string                  `json:"product_variant_sku"`
	StatusId               *uint                   `json:"status_id"`
	ChannelStatus          core_dto.LookupCodesDTO `json:"status"`
	ExternalId             string                  `json:"external_id"`
	ExternalDetails        datatypes.JSON          `json:"external_details"`
	Inventory              interface{}             `json:"inventory"`
	Pricing                interface{}             `json:"pricing"`
	ProductTemplateDetails interface{}             `json:"product_template_details"`
	ProductVariantDetails  interface{}             `json:"product_variant_details"`
}
type ChannelInventoryResponseDTO struct {
	Id                 uint           `json:"id"`
	ChannelCode        string         `json:"channel_code"`
	PhysicalLocationId uint           `json:"physical_location_id"`
	PhysicalLocation   datatypes.JSON `json:"physical_location"`
	Quantity           int64          `json:"quantity"`
	OnHandQuantity     int64          `json:"on_hand_quantity"`
	PlannedInQuantity  int64          `json:"planned_in_quantity"`
	PlannedOutQuantity int64          `json:"planned_out_quantity"`
	BinTag             string         `json:"bin_tag"`
	AvailableStock     int64          `json:"available_stock"`
	StockExpected      int64          `json:"stock_expected"`
	CommitedStock      int64          `json:"commited_stock"`
	NoOfUnits          int64          `json:"no_of_units"`
	IsEnabled          *bool          `json:"is_enabled" gorm:"default:true"`
	IsActive           *bool          `json:"is_active" gorm:"default:true"`
}
type ChannelPricingResponseDTO struct {
	Id                 uint    `json:"id"`
	SPL_id             uint    `json:"spL_id"`
	CategoryCommission int     `json:"category_commission"`
	Uom_id             *uint   `json:"uom_id"`
	Mrp                float64 `json:"mrp"`
	SaleRate           float64 `json:"sale_rate"`
	Duties             float64 `json:"duties"` //product_taxes
	Price              float64 `json:"price"`
}

// ===========================Filter Tab Response ===================================
type PriceListTabResponse struct {
	SalesPriceList    []SalesPriceListTab    `json:"sales_price_list"`
	PurchasePriceList []PurchasePriceListTab `json:"purchase_price_list"`
}
type SalesPriceListTab struct {
	Id                     uint                    `json:"id"`
	SPL_id                 uint                    `json:"spL_id"`
	CategoryCommission     int                     `json:"category_commission"`
	Uom_id                 *uint                   `json:"uom_id"`
	UOM                    mdm_UOM.UomResponseDTO  `json:"uom_details"`
	QuantityValue          datatypes.JSON          `json:"quantity_value"`
	Quantity_value_type_id *uint                   `json:"quantity_value_type_id"`
	Quantity_value_type    core_dto.LookupCodesDTO `json:"quantity_value_type"`
	Mrp                    float64                 `json:"mrp"`
	SaleRate               float64                 `json:"sale_rate"`
	Duties                 float64                 `json:"duties"` //product_taxes
	Pricing_options_id     *uint                   `json:"pricing_options_id"`
	Pricing_options        core_dto.LookupCodesDTO `json:"pricing_options"`
	Price                  float64                 `json:"price"`
	Add_channel_of_sale_id *uint                   `json:"add_channel_of_sale_id"`
	Add_channel_of_sale    core_dto.LookupCodesDTO `json:"add_channel_of_sale"`
	Price_list_name        string                  `json:"price_list_name"`
}
type PurchasePriceListTab struct {
	PPL_id                 uint
	Vendor_rate            int                     `json:"vendor_rate"`
	LeadTime               string                  `json:"lead_time"`
	Price_quantity         float64                 `json:"price_quantity"`
	CreditPeriod           string                  `json:"credit_period"`
	ExpectedDeliveryTime   string                  `json:"expected_delivery_time"`
	SalesPeriod            string                  `json:"sales_period"`
	Mrp                    float64                 `json:"mrp"`
	QuantityValue          datatypes.JSON          `json:"quantity_value"`
	Quantity_value_type_id *uint                   `json:"quantity_value_type_id"`
	Quantity_value_type    core_dto.LookupCodesDTO `json:"quantity_value_type"`
	MinimumOrderQuantity   int                     `json:"minimum_order_quantity"`
	Price                  float64                 `json:"price"`
	VendorName             string                  `json:"vendor_name"`
	Price_list_name        string                  `json:"price_list_name"`
	VendorCreditPeriod     string                  `json:"vendor_credit_period"`
}

type ProductBundlesTabResponse struct {
	Id           uint                    `json:"id"`
	BundleName   string                  `json:"bundle_name"`
	BundleCode   string                  `json:"bundle_id"`
	SalesPrice   float64                 `json:"sales_price"`
	Instructions string                  `json:"instructions"`
	Description  datatypes.JSON          `json:"description"`
	ImageOptions datatypes.JSON          `json:"image_options"`
	StatusId     *uint                   `json:"status_id"`
	Status       core_dto.LookupCodesDTO `json:"status"`
	CostPrice    float64                 `json:"cost_price"`
	MRP          float64                 `json:"mrp" `
	Tax          bool                    `json:"tax"`
	Quantity     int                     `json:"quantity"`
}

type Offers struct {
	OfferId        uint      `json:"offer_id,omitempty"`
	Title          string    `json:"title,omitempty"`
	Value          float64   `json:"value,omitempty"`
	DiscountTypeId uint      `json:"discount_type_id,omitempty"`
	ValidFrom      time.Time `json:"valid_from,omitempty"`
	ValidTo        time.Time `json:"valid_to,omitempty"`
}
