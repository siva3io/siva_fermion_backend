package notifications

type NotificationTitle string
type NotificationSuccessMessage string
type NotificationScopeType string


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


const (
	PRODUCT_CREATED_TITLE     NotificationTitle = "Product Created"
	SALES_ORDER_CREATED_TITLE NotificationTitle = "Sales Order Created"
	CONTACT_CREATED_TITLE     NotificationTitle = "Contact Created"
)

const (
	PRODUCT_CREATED_MESSAGE     NotificationSuccessMessage = "product created successfully"
	SALES_ORDER_CREATED_MESSAGE NotificationSuccessMessage = "sales order created successfully"
	CONTACT_CREATED_MESSAGE     NotificationSuccessMessage = "contact created successfully"
)

const (
	PRODUCT_LIST     NotificationScopeType = "PRODUCT LIST"
)
