package wallets

import (
	cmiddleware "fermion/backend_core/middleware"

	"github.com/labstack/echo/v4"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
func (h *handler) Route(g *echo.Group) {
	g.GET("", h.ListWallets, cmiddleware.Authorization)
	g.GET("/:id", h.ViewWallet, cmiddleware.Authorization)
	g.POST("/create", h.CreateWallet, cmiddleware.Authorization, WalletCreateValidate)
	g.POST("/:id/update", h.UpdateWallet, cmiddleware.Authorization, WalletUpdateValidate)
	g.DELETE("/:id/delete", h.DeleteWallet, cmiddleware.Authorization)
	g.GET("/get_wallet", h.GetWallet, cmiddleware.Authorization)
	g.POST("/add_money", h.AddMoney, cmiddleware.Authorization, WalletUpdateValidate)
	g.POST("/deduct_money", h.DeductMoney, cmiddleware.Authorization, WalletUpdateValidate)
}
