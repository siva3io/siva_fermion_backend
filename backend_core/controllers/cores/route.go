package cores

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
	g.GET("/lookup_types", h.GetLookupTypes)
	g.GET("/lookup_codes/:type", h.GetLookupCodes)
	g.GET("/countries", h.GetCountries)
	g.GET("/states/:id", h.GetStates)
	g.GET("/currencies", h.GetCurrencies)
	g.GET("/lookup_codes/search", h.SearchLookupCodes)
	g.POST("/lookup_types/create", h.CreateLookupType, LookupTypeCreateValidate, cmiddleware.Authorization)
	g.POST("/lookup_codes/create", h.CreateLookupCodes, LookupCodeCreateValidate, cmiddleware.Authorization)
	g.POST("/lookup_type/:id/update", h.UpdateLookupType, LookupTypeUpdateValidate, cmiddleware.Authorization)
	g.POST("/lookup_code/:id/update", h.UpdateLookupCode, LookupCodeUpdateValidate, cmiddleware.Authorization)
	g.DELETE("/lookup_type/:id/delete", h.DeleteLookupType, cmiddleware.Authorization)
	g.DELETE("/lookup_code/:id/delete", h.DeleteLookupCode, cmiddleware.Authorization)
	g.GET("/channels/lookup_codes", h.GetChannelLookupCodes)

	g.POST("/cache/update", h.UpdateCacheData)
	g.GET("/cache/get/:key", h.GetCacheData, cmiddleware.Authorization)
	//not sure where "/lookup_types" are used so using "/lookup_types/list" apis are developed
	g.GET("/lookup_types/list", h.GetAllLookupTypes, cmiddleware.Authorization)
	g.GET("/lookup_codes/list", h.GetAllLookupCodes, cmiddleware.Authorization)

	// apps
	g.GET("/store/apps", h.ListStoreApps)
	g.GET("/store/apps/:id", h.GetApp)
	g.GET("/store/search", h.SearchApps)
	g.GET("/apps/installed", h.ListInstalledApps, cmiddleware.Authorization)
	g.POST("/store/apps/:app_code/install", h.InstallApp, cmiddleware.Authorization)
	g.POST("/store/apps/:app_code/uninstall", h.UninstallApp, cmiddleware.Authorization)

	// meta data
	g.GET("/meta_data", h.ListMetaData)
	g.GET("/meta_data/:model_name", h.ViewMetaData)
	g.POST("/company_preferences/:id/update", h.UpdateCompanyPreferences, cmiddleware.Authorization)
	g.GET("/company_preferences/:id", h.GetCompanyPreferences, cmiddleware.Authorization)
	g.GET("/company/:id/list_users", h.ListCompanyUsers, cmiddleware.Authorization)
	g.POST("/company/:id/update", h.UpdateCompany, cmiddleware.Authorization)

}
