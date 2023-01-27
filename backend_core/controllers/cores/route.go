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

	// cache api's
	g.POST("/local_cache/update", h.UpdateLocalCacheData)
	g.GET("/local_cache/get/:key", h.GetLocalCacheData, cmiddleware.Authorization)
	g.POST("/redis_cache/update", h.UpdateRedisCacheData)
	g.GET("/redis_cache/get/:key", h.GetRedisCacheData, cmiddleware.Authorization)

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
	g.POST("/apps/installed/update", h.UpdateInstalledApp, cmiddleware.Authorization)
	g.POST("/app/installed/renew", h.RenewSubscription, cmiddleware.Authorization)

	// meta data
	g.GET("/meta_data", h.ListMetaData)
	g.GET("/meta_data/:model_name", h.ViewMetaData)

	//company APIs
	g.GET("/company_preferences/:id", h.GetCompanyPreferences, cmiddleware.Authorization)
	g.GET("/company_preferences", h.ListCompanyPreferences, cmiddleware.Authorization)
	g.GET("/company/:id/list_users", h.ListCompanyUsers, cmiddleware.Authorization)
	g.POST("/company/:id/update", h.UpdateCompany, cmiddleware.Authorization)
	g.POST("/company/:id/register_ondc", h.OndcRegistration, cmiddleware.Authorization)
	g.POST("/ondc/update", h.UpdateOndcDetails, cmiddleware.Authorization)
	g.POST("/validate_aadhaar", h.AadhaarValidation, cmiddleware.Authorization)

	//channel status api's
	g.GET("/channel_status", h.ListChannelStatus, cmiddleware.Authorization)
	g.GET("/channel_status/:id", h.ViewChannelStatus, cmiddleware.Authorization)
	g.POST("/channel_status/create", h.CreateChannelStatus, cmiddleware.Authorization)
	g.POST("/channel_status/:id/update", h.UpdateChannelStatus, cmiddleware.Authorization)
	g.DELETE("/channel_status/:id/delete", h.DeleteChannelStatus, cmiddleware.Authorization)

	g.POST("/custom_solution", h.RequestSolution, cmiddleware.Authorization)

	//================================ADMIN==========================================
	g.GET("/company/admin/list", h.ListCompaniesAdmin, cmiddleware.Authorization)
}
