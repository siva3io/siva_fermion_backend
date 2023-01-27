package basic_inventory

import (
	cmiddleware "fermion/backend_core/middleware"

	"github.com/labstack/echo/v4"
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
func (h *handler) Route(g *echo.Group) {
	g.POST("/centralized/create", h.CreateCentrailizedInventoryEvent, cmiddleware.Authorization, BasicInventoryCentralCreateValidate)
	g.POST("/centralized/:id/update", h.UpdateCentrailizedInventoryEvent, cmiddleware.Authorization, BasicInventoryCentralUpdateValidate)
	g.DELETE("/centralized/:id/delete", h.DeleteCentrailizedInventory, cmiddleware.Authorization)
	g.GET("/centralized/:id", h.GetCentrailizedInventory, cmiddleware.Authorization)
	g.GET("/centralized", h.GetCentrailizedInventoryList, cmiddleware.Authorization)
	g.GET("/centralized/dropdown", h.GetCentrailizedInventoryListDropdown, cmiddleware.Authorization)

	g.POST("/decentralized/create", h.CreateDecentralizedInventoryEvent, cmiddleware.Authorization, BasicInventoryDeCentralCreateValidate)
	g.POST("/decentralized/:id/update", h.UpdateDecentralizedInventoryEvent, cmiddleware.Authorization, BasicInventoryDeCentralUpdateValidate)
	g.DELETE("/decentralized/:id/delete", h.DeleteDecentralizedInventory, cmiddleware.Authorization)
	g.GET("/decentralized/:id", h.GetDecentralizedInventory, cmiddleware.Authorization)
	g.GET("/decentralized", h.GetDecentralizedInventoryList, cmiddleware.Authorization)
	g.GET("/decentralized/dropdown", h.GetDecentralizedInventoryListDropdown, cmiddleware.Authorization)

	//-------------Channel API's----------------------------------------------------
	g.POST("/channels/upsert", h.DeCentralisedInventoryUpsertEvent, cmiddleware.Authorization)
	//-------------Inventory Transaction Log API's----------------------------------------------------
	g.POST("/logs/create", h.InventoryTransactionCreate, cmiddleware.Authorization)
}
