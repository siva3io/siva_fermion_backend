package mdm

import (
	"time"

	model_core "fermion/backend_core/internal/model/core"

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
type ProductBrand struct {
	model_core.Model
	Name        string `json:"brand_name" gorm:"type:varchar(50)"`
	Description string `json:"description"`
}
type ProductCategory struct {
	model_core.Model
	Name                   string                `json:"name"`
	ParentCategoryID       *uint                 `json:"parent_category_id"`
	ExternalID             *uint                 `json:"external_id"`
	ChildCategoryIds       []*ProductCategory    `json:"child_ids" gorm:"foreignkey:ParentCategoryID;references:ID"`
	RelatedCategoryIds     datatypes.JSON        `json:"related_categorie_ids" gorm:"type:json"`
	RecommendedCategoryIds []*ProductCategory    `json:"recommended_category_ids" gorm:"many2many:recommended_category_ids"`
	AdditivesInformation   string                `json:"additives_information"`
	ShortDescription       string                `json:"short_description"`
	DomainId               *uint                 `json:"domain_id"`
	Domain                 model_core.Lookupcode `json:"domain" gorm:"foreignkey:DomainId;references:ID"`
	CategoryCode           string                `json:"category_code"`
}
type ProductBaseAttributes struct {
	model_core.Model
	Name string `json:"name" `
}
type ProductBaseAttributesValues struct {
	model_core.Model
	BaseAttributeId uint   `json:"base_attribute_id"`
	Name            string `json:"name" `
}
type ProductSelectedAttributes struct {
	model_core.Model
	TemplateId  uint   `json:"template_id"`
	AttributeId uint   `json:"attribute_id"`
	Name        string `json:"name"`
}
type ProductSelectedAttributesValues struct {
	model_core.Model
	TemplateId       uint   `json:"template_id"`
	AttributeId      uint   `json:"attribute_id"`
	AttributeValueId uint   `json:"attribute_value_id"`
	Name             string `json:"name" `
}
type ProductTemplate struct {
	model_core.Model
	Name                           string                `json:"product_name" gorm:"text;not null"`
	BrandID                        *uint                 `json:"brand_id"`
	Brand                          ProductBrand          `json:"brand" gorm:"foreignkey:BrandID;references:ID"`
	SkuCode                        string                `json:"sku_code" gorm:"type:varchar(100);UNIQUE"`
	HSNCODE                        string                `json:"hsn_code"`
	HSNCodesData                   HSNCodesData          `json:"hsn_codes_data" gorm:"foreignkey:HSNCODE; references:HsnCode"`
	ProductConditionID             *uint                 `json:"product_condition_id"`
	ProductCondition               model_core.Lookupcode `json:"product_condition" gorm:"foreignkey:ProductConditionID;references:ID"`
	ProductTypeID                  *uint                 `json:"product_type_id"`
	ProductType                    model_core.Lookupcode `json:"product_type" gorm:"foreignkey:ProductTypeID;references:ID"`
	ProductProcurementTreatmentIds datatypes.JSON        `json:"product_procurement_treatment_ids" gorm:"type:json"`
	//EstimatedDeliveryTimeID        *uint                 `json:"estimated_delivery_time_id"`
	//EstimatedDeliveryTime          model_core.Lookupcode `json:"estimated_delivery_time" gorm:"foreignkey:EstimatedDeliveryTimeID;references:ID"`
	StockTreatmentIds      datatypes.JSON        `json:"stock_treatment_ids" gorm:"type:json"`
	InventoryTrackingID    *uint                 `json:"inventory_tracking_id"`
	InventoryTracking      model_core.Lookupcode `json:"inventory_tracking" gorm:"foreignkey:InventoryTrackingID;references:ID"`
	UomID                  *uint                 `json:"uom_id"`
	Uom                    Uom                   `json:"uom" gorm:"foreignkey:UomID; references:ID"`
	SecondaryUom           datatypes.JSON        `json:"secondary_uom" gorm:"type:json;default:'{}'"`
	ImageOptions           datatypes.JSON        `json:"image_options" gorm:"type:json;default:'[]'"`
	PrimaryCategoryID      *uint                 `json:"primary_category_id"`
	PrimaryCategory        ProductCategory       `json:"primary_category" gorm:"foreignkey:PrimaryCategoryID;references:ID"`
	SecondaryCategoryID    *uint                 `json:"secondary_category_id"`
	SecondaryCategory      ProductCategory       `json:"secondary_category" gorm:"foreignkey:SecondaryCategoryID;references:ID"`
	Description            datatypes.JSON        `json:"description" gorm:"type:json;default:'{}'"`
	ProductVariantIds      []ProductVariant      `json:"product_variant_ids" gorm:"foreignkey:ProductTemplateId;references:ID"`
	ProductPricingDetails  ProductPricingDetails `json:"product_pricing_details" gorm:"embedded"`
	VendorPriceListIds     pq.Int64Array         `json:"vendor_price_list_ids" gorm:"type:int[]"`
	PriceListDetails       datatypes.JSON        `json:"price_list_details" gorm:"type:json"`
	ShippingOptions        datatypes.JSON        `json:"shipping_options" gorm:"type:json"`
	PackageMaterialOptions datatypes.JSON        `json:"package_material_options" gorm:"type:json;default:'[]'"`
	PackageDimensions      datatypes.JSON        `json:"package_dimensions" gorm:"type:json;default:'{}'"`
	StatusId               *uint                 `json:"status_id"`
	Status                 model_core.Lookupcode `json:"status" gorm:"foreignkey:StatusId;references:ID"`
	AdditivesInformation   string                `json:"additives_information"`
	ShortDescription       string                `json:"short_description"`
	DomainId               *uint                 `json:"domain_id"`
	Domain                 model_core.Lookupcode `json:"domain" gorm:"foreignkey:DomainId;references:ID"`

	ManufacturerDetails          ManufacturerDetails      `json:"manufacturer_details" gorm:"embedded"`
	ProductCriticalDetails       ProductCriticalDetails   `json:"product_critical_details" gorm:"embedded"`
	Location                     Location                 `json:"location" gorm:"embedded"`
	DeliverySlaDetails           DeliverySlaDetails       `json:"delivery_sla_details" gorm:"embedded"`
	InventoryDetail              InventoryDetails         `json:"inventory_detail" gorm:"embedded"`
	ProductCancellationTermsFlag *bool                    `json:"product_cancellation_terms_flag" gorm:"type:boolean;default:false"`
	ProductCancellationTerms     ProductCancellationTerms `json:"product_cancellation_terms"  gorm:"embedded"`
	ProductReturnTermsFlag       *bool                    `json:"product_return_terms_flag" gorm:"type:boolean;default:false"`
	ProductReturnTerms           ProductReturnTerms       `json:"product_return_terms" gorm:"embedded"`
	ProductReplacementTermsFlag  *bool                    `json:"product_replacement_terms_flag" gorm:"type:boolean;default:false"`
	ProductReplacementTerms      ProductReplacementTerms  `json:"product_replacement_terms" gorm:"embedded"`
	FoodItemDetails              FoodItemDetails          `json:"food_item_details" gorm:"embedded"`

	TemplateOptions datatypes.JSON `json:"template_options" gorm:"type:json"`
	Related         datatypes.JSON `json:"related" gorm:"type:json"`
	Recommended     datatypes.JSON `json:"recommended" gorm:"type:json"`
	OtherStatuses   datatypes.JSON `json:"other_statuses" gorm:"type:json"`

	Channel string `json:"channel" gorm:"default:'EUNIMART'"`
}
type ProductVariant struct {
	model_core.Model
	SerialNumber           string                         `json:"serial_number"`
	ProductTemplateId      uint                           `json:"product_template_id"`
	ParentSkuId            string                         `json:"parent_sku_id" gorm:"type:varchar(50)"`
	SkuId                  string                         `json:"sku_id" gorm:"type:varchar(100);UNIQUE"`
	ProductName            string                         `json:"product_name" gorm:"text"`
	AttributeKeyValuesId   pq.Int64Array                  `json:"attribute_key_values_id" gorm:"type:int[]"`
	ImageOptions           datatypes.JSON                 `json:"image_options" gorm:"type:json;default:'[]'"`
	VariantTypeId          *uint                          `json:"variant_type_id"`
	VariantType            model_core.Lookupcode          `json:"variant_type" gorm:"foreignkey:VariantTypeId;references:ID" `
	Barcode                string                         `json:"barcode"`
	StandardProductTypes   datatypes.JSON                 `json:"standard_product_types" gorm:"type:json"`
	StandardProductTypeId  string                         `json:"standard_product_type_id" gorm:""`
	ConditionID            *uint                          `json:"condition_id"`
	Condition              model_core.Lookupcode          `json:"condition" gorm:"foreignkey:ConditionID; references:ID"`
	CategoryID             *uint                          `json:"category_id" gorm:""`
	Category               ProductCategory                `json:"category" gorm:"foreignkey:CategoryID;references:ID"`
	LeafCategoryID         *uint                          `json:"leaf_category_id" gorm:"column:secondary_category_id"`
	LeafCategory           ProductCategory                `json:"leaf_category" gorm:"foreignkey:LeafCategoryID;references:ID"`
	Description            datatypes.JSON                 `json:"description" gorm:"type:json;default:'{}'"`
	ProductDimensions      datatypes.JSON                 `json:"product_dimensions" gorm:"type:json"`
	PackageDimensions      datatypes.JSON                 `json:"package_dimensions" gorm:"type:json"`
	PackageMaterialOptions datatypes.JSON                 `json:"package_material_options" gorm:"type:json"`
	VendorPriceListIds     pq.Int64Array                  `json:"vendor_price_list_ids" gorm:"type:int[]"`
	PriceListDetails       datatypes.JSON                 `json:"price_list_details" gorm:"type:json"`
	ShippingOptions        datatypes.JSON                 `json:"shipping_options" gorm:"type:json"`
	ProductPricingDetails  ProductPricingDetails          `json:"product_pricing_details" gorm:"embedded"`
	InventoryDetails       []*DecentralizedBasicInventory `json:"inventory_details" gorm:"foreignkey:ProductVariantId; references:ID"`
	StatusId               *uint                          `json:"status_id"`
	Status                 model_core.Lookupcode          `json:"status" gorm:"foreignkey:StatusId;references:ID"`
	AdditivesInformation   string                         `json:"additives_information"`
	ShortDescription       string                         `json:"short_description"`
	DomainId               *uint                          `json:"domain_id"`
	Domain                 model_core.Lookupcode          `json:"domain" gorm:"foreignkey:DomainId;references:ID"`

	KeywordIds               datatypes.JSON           `json:"keyword_ids" gorm:"type:json;default:'{}'"`
	CostDetails              datatypes.JSON           `json:"cost_details" gorm:"type:json;default:'{}'"`
	ForecastingOptions       datatypes.JSON           `json:"forecasting_options" gorm:"type:json;default:'{}'"`
	ValidationInfo           string                   `json:"validation_info" gorm:"default:'Valid'"`
	PackageTemplateOptions   datatypes.JSON           `json:"package_template_options" gorm:"type:json;default:'{}'"`
	ManufacturerDetails      ManufacturerDetails      `json:"manufacturer_details" gorm:"embedded"`
	ProductCriticalDetails   ProductCriticalDetails   `json:"product_critical_details" gorm:"embedded"`
	Location                 Location                 `json:"location" gorm:"embedded"`
	DeliverySlaDetails       DeliverySlaDetails       `json:"delivery_sla_details" gorm:"embedded"`
	InventoryDetail          InventoryDetails         `json:"inventory_detail" gorm:"embedded"`
	ProductCancellationTerms ProductCancellationTerms `json:"product_cancellation_terms" gorm:"embedded"`
	ProductReturnTerms       ProductReturnTerms       `json:"product_return_terms" gorm:"embedded"`
	ProductReplacementTerms  ProductReplacementTerms  `json:"product_replacement_terms" gorm:"embedded"`
	FoodItemDetails          FoodItemDetails          `json:"food_item_details" gorm:"embedded"`
	OfferDetails             Offers                   `json:"offer_details" gorm:"embedded"`

	Channel       string  `json:"channel" gorm:"default:'EUNIMART'"`
	RatingAverage float64 `json:"rating_average"`
	RatingCount   uint    `json:"rating_count"`
}
type ProductPricingDetails struct {
	PaymentMethodId *uint                 `json:"payment_method_id"`
	PaymentMethod   model_core.Lookupcode `json:"payment_method" gorm:"foreignkey:PaymentMethodId; references:ID"`
	SalesPrice      float64               `json:"sales_price"`
	DeclaredPrice   float64               `json:"declared_price"`
	CostPrice       float64               `json:"cost_price"`
	MRP             float64               `json:"mrp"`
	TaxOptions      float64               `json:"tax_options"`
	CurrencyId      *uint                 `json:"currency_id"`
	Currency        model_core.Currency   `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	Tax             bool                  `json:"tax" gorm:"default:false"`
	Shipping        bool                  `json:"shipping" gorm:"default:false"`
}
type ProductBundles struct {
	model_core.Model
	BundleCode     string                `json:"bundle_id" gorm:"unique;not null"`
	BundleName     string                `json:"bundle_name"`
	Instructions   string                `json:"instructions"`
	Description    datatypes.JSON        `json:"description"`
	ImageOptions   datatypes.JSON        `json:"image_options"`
	Products       []BundleLineItems     `json:"products" gorm:"foreignkey:BundleId; references:ID"`
	SalesPrice     float64               `json:"sales_price"`
	CostPrice      float64               `json:"cost_price"`
	MRP            float64               `json:"mrp" `
	TaxOptions     float64               `json:"tax_options"`
	CurrencyId     *uint                 `json:"currency_id"`
	Currency       model_core.Currency   `json:"currency" gorm:"foreignkey:CurrencyId; references:ID"`
	SellingPrice   float64               `json:"selling_price"`
	Tax            bool                  `json:"tax"`
	Shipping       bool                  `json:"shipping"`
	PackageDetails datatypes.JSON        `json:"package_details"`
	StatusId       *uint                 `json:"status_id"`
	Status         model_core.Lookupcode `json:"status" gorm:"foreignkey:StatusId;references:ID"`
}
type BundleLineItems struct {
	BundleId         uint           `json:"bundle_id"`
	ProductVariantId uint           `json:"product_variant_id"`
	ProductVariant   ProductVariant `json:"product_variant"`
	Quantity         int            `json:"quantity"`
	model_core.Model
}
type ProductToChannelMap struct {
	model_core.Model
	ChannelCode       string                `json:"channel_code"`
	ProductTemplateId uint                  `json:"product_template_id"`
	TemplateSku       string                `json:"template_sku"`
	ProductVariantId  uint                  `json:"product_variant_id"`
	ProductVariantSku string                `json:"product_variant_sku"`
	StatusId          *uint                 `json:"status_id"`
	ChannelStatus     model_core.Lookupcode `json:"status" gorm:"foreignkey:StatusId;references:ID"`
	ExternalId        string                `json:"external_id"`
	ExternalDetails   datatypes.JSON        `json:"external_details" gorm:"column:external_details;type:json;default:'[]';not null"`
}

type Location struct {
	Pincode      string            `json:"pincode"`
	City         string            `json:"city"`
	StateId      *uint             `json:"state_id"`
	State        *model_core.State `json:"state" gorm:"foreignkey:StateId;references:ID"`
	AddressLine1 string            `json:"address_line1"`
	AddressLine2 string            `json:"address_line2"`
	AddressLine3 string            `json:"address_line3"`
}

type DeliverySlaDetails struct {
	EstimatedDeliveryTimeId *uint                 `json:"estimated_delivery_time_id"`
	EstimatedDeliveryTime   model_core.Lookupcode `json:"estimated_delivery_time" gorm:"foreignkey:EstimatedDeliveryTimeId;references:ID"`
}

type InventoryDetails struct {
	MaximumQuantity   uint `json:"maximum_quantity"`
	AvailableQuantity uint `json:"available_quantity"`
}

type ProductCriticalDetails struct {
	UnitPerBox         uint                  `json:"unit_per_box"`
	TimeToShipId       *uint                 `json:"time_to_ship_id"`
	TimeToShip         model_core.Lookupcode `json:"time_to_ship" gorm:"foreignkey:TimeToShipId;references:ID"`
	CustomerCareNumber string                `json:"customer_care_number"`
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
	RefundEligible  *bool                 `json:"refund_eligible" gorm:"default:false"`
	CancelTimeId    *uint                 `json:"cancel_time_id"`
	CancelTime      model_core.Lookupcode `json:"cancel_time" gorm:"foreignkey:CancelTimeId;references:ID"`
	CancellationFee uint                  `json:"cancellation_fee"`
}
type ProductReturnTerms struct {
	RefundEligibleReturn   *bool                 `json:"refund_eligible_return" gorm:"default:false"`
	FulfillmentManagedById *uint                 `json:"fulfillment_managed_by_id"`
	FulfillmentManagedBy   model_core.Lookupcode `json:"fulfillment_managed_by" gorm:"foreignkey:FulfillmentManagedById;references:ID"`
	ReturnWithinId         *uint                 `json:"return_within_id"`
	ReturnWithin           model_core.Lookupcode `json:"return_within" gorm:"foreignkey:ReturnWithinId;references:ID"`
}
type ProductReplacementTerms struct {
	ReplacementEligible *bool                 `json:"replacement_eligible" gorm:"default:false"`
	ReplacementWithinId *uint                 `json:"replacement_within_id"`
	ReplacementWithin   model_core.Lookupcode `json:"replacement_within" gorm:"foreignkey:ReplacementWithinId;references:ID"`
}
type FoodItemDetails struct {
	FoodTypeId              *uint                 `json:"food_type_id"`
	FoodType                model_core.Lookupcode `json:"food_type" gorm:"foreignkey:FoodTypeId ;references:ID"`
	FSSAILicenceNumber      string                `json:"fssai_licence_number"`
	OtherFSSAILicenceNumber string                `json:"other_fssai_licence_number"`
	ImportersFSSAINumber    string                `json:"importers_fssai_number"`
	TimeToLife              string                `json:"time_to_life"`
	IngredientsInfo         string                `json:"ingredients_info"`
	NutritionalInfo         string                `json:"nutritional_info"`
}

type HSNCodesData struct {
	model_core.Model
	SNo         string    `json:"sno"`
	HsnCode     string    `json:"hsn_code" gorm:"unique;not null"`
	Description string    `json:"description"`
	CGSTRate    float64   `json:"cgst_rate"`
	IGSTRate    float64   `json:"igst_rate"`
	SGSTRate    float64   `json:"sgst_rate"`
	CreatedDate time.Time `json:"created_date" `
	UpdatedDate time.Time `json:"updated_date"`
	Status      string    `json:"status"`
}
type Offers struct {
	OfferId        *uint                 `json:"offer_id"`
	Title          string                `json:"title"`
	Value          float64               `json:"value"`
	DiscountTypeId *uint                 `json:"discount_type_id"`
	DiscountType   model_core.Lookupcode `json:"discount_type" gorm:"foreignkey:DiscountTypeId ;references:ID"`
	ValidFrom      time.Time             `json:"valid_from"`
	ValidTo        time.Time             `json:"valid_to"`
}
