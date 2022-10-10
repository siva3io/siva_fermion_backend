package helpers

var SourceDocumentTypeWithModelName = map[string]string{
	"SALES_ORDERS":    "sales_orders",
	"PURCHASE_RETURN": "purchase_returns",
	"ASN":             "asns",
	"IST":             "internal_transfers",
}

var SourceDocumentTypeWithColumnName = map[string]string{
	"SALES_ORDERS":    "status_id",
	"PURCHASE_RETURN": "status_id",
	"ASN":             "status_id",
	"IST":             "status_id",
}