package rating

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
	g.POST("/create_rating", h.CreateRatingEvent, cmiddleware.Authorization)
	g.POST("/:id/update_rating", h.UpdateRatingEvent, cmiddleware.Authorization, RatingUpdateValidate)
	g.GET("/:id/rating", h.GetRating, cmiddleware.Authorization)
	g.GET("/rating_list", h.ListRating, cmiddleware.Authorization)
	g.DELETE("/:id/rating_delete", h.DeleteRating, cmiddleware.Authorization)

	g.POST("/create_feedbackform", h.CreateFeedBackFormEvent, cmiddleware.Authorization, FeedBackFormCreateValidate)
	g.GET("/rating_category_list", h.GetRatingCategoryList, cmiddleware.Authorization)
	g.GET("/feedback_category_list", h.GetFeedBackCategoryList, cmiddleware.Authorization)
	g.GET("/list_feedback_form", h.ListFeedBackForm, cmiddleware.Authorization)

	g.POST("/create_feedback", h.CreateFeedbackEvent, cmiddleware.Authorization, FeedbackCreateValidate)
	g.POST("/:id/update_feedback", h.UpdateFeedbackEvent, cmiddleware.Authorization, FeedbackUpdateValidate)
	g.GET("/:id/get_feedback", h.GetFeedback, cmiddleware.Authorization)
	g.GET("/feedback_list", h.ListFeedBack, cmiddleware.Authorization)
	g.DELETE("/:id/feedback_delete", h.DeleteFeedback, cmiddleware.Authorization)

	g.POST("/create_report", h.CreateReportEvent, cmiddleware.Authorization, ReportCreateValidate)
	g.POST("/:id/update_report", h.UpdateReportEvent, cmiddleware.Authorization, ReportUpdateValidate)
	g.GET("/:id/get_report", h.GetReport, cmiddleware.Authorization)
	g.GET("/report_list", h.ListReport, cmiddleware.Authorization)
	g.DELETE("/:id/report_delete", h.DeleteReport, cmiddleware.Authorization)

	g.GET("/products_rating_sync", h.ProductsRatingSync, cmiddleware.Authorization)
	g.GET("/company_rating_sync", h.CompanyRatingSync, cmiddleware.Authorization)
}
