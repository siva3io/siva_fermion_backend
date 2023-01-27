package rating

import (
	"errors"
	"fmt"
	"time"

	"fermion/backend_core/controllers/cores"
	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/controllers/mdm/products"
	scheduler "fermion/backend_core/controllers/scheduler/app"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	model "fermion/backend_core/internal/model/rating"
	"fermion/backend_core/ipaas_core/utils"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

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
type handler struct {
	service          Service
	coreService      cores.Service
	productservice   products.Service
	schedulerservice scheduler.Service
}

var RatingHandler *handler

func NewHandler() *handler {
	if RatingHandler != nil {
		return RatingHandler
	}
	return &handler{NewService(), cores.NewService(), products.NewService(), scheduler.NewService()}
}

func (h *handler) CreateRatingEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	//    var input interface{}
	input := utils.GetBodyData(c).(map[string]interface{})
	fmt.Println("input data-->", input)
	//    c.Request().GetBody(&input)
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      input,
	}
	eda.Produce(eda.CREATE_RATING, request_payload)
	return res.RespSuccess(c, "Rating Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateRating(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(model.Rating)

	p := new(pagination.Paginatevalue)
	category_name := data["rating_category"].(string)
	p.Filters = "[[\"category_name\",\"=\",\"" + category_name + "\"]]"
	rating_category_interface, err := h.service.RatingCategoriesList(edaMetaData, p)
	check := rating_category_interface.([]model.RatingCategory)
	if len(check) > 0 {
		createPayload.RatingCategoryId = rating_category_interface.([]model.RatingCategory)[0].ID
	} else {
		fmt.Println("rating_category_name", p.Filters)
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	data2 := createPayload.RelatedId
	helpers.JsonMarshaller(data["id"], &data2)

	data3 := createPayload.RatingValue
	helpers.JsonMarshaller(data["value"], &data3)

	createPayload.RelatedId = uint(data["id"].(float64))
	createPayload.RatingValue = int64(data["value"].(float64))
	feedback := map[string]interface{}{
		"feedback_form": data["feedback_form"],
		"feedback_id":   data["feedback_id"],
	}
	data1 := createPayload.FeedBack
	helpers.JsonMarshaller(feedback, &data1)
	createPayload.FeedBack = data1

	// data.value.(map[string]interface{})=createPayload.RatingValue
	err = h.service.SaveRating(edaMetaData, createPayload)
	if err != nil {
		return
	}
	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_RATING_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_RATING_ACK, responseMessage)
}

func (h *handler) UpdateRatingEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	RatingId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         RatingId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("rating"),
	}
	eda.Produce(eda.UPDATE_RATING, request_payload)
	return res.RespSuccess(c, "Rating Update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateRating(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(model.Rating)
	helpers.JsonMarshaller(data, updatePayLoad)
	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdateRating(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_RATING_ACK, responsemessage)
		return
	}
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_RATING_ACK, responsemessage)
}

func (h *handler) ListRating(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	result, err := h.service.GetRatingList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	var dtoData []RatingResponseDTO
	err = helpers.JsonMarshaller(result, &dtoData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", dtoData, p)
}

func (h *handler) GetRating(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	RatingId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         RatingId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetRating(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Rating list Retrieved Successfully", result)
}

func (h *handler) DeleteRating(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	RatingId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         RatingId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
		//"deleted_by": metaData.TokenUserId,
	}

	err = h.service.DeleteRating(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", RatingId)
}

func (h *handler) CreateFeedBackFormEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("rating"),
	}
	eda.Produce(eda.CREATE_FEEDBACKFORM, request_payload)
	return res.RespSuccess(c, "Feedbackform Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateFeedBackForm(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(model.FeedbackForm)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.SaveFeedBackForm(edaMetaData, createPayload)
	createPayload.CreatedByID = &edaMetaData.TokenUserId

	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_FEEDBACKFORM_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_FEEDBACKFORM_ACK, responseMessage)

}

func (h *handler) GetRatingCategoryList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	result, err := h.service.RatingCategoriesList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	var dtoData []RatingCategoryDto
	err = helpers.JsonMarshaller(result, &dtoData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, " List Retrieved successfully", dtoData, p)
}

func (h *handler) GetFeedBackCategoryList(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"

	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	result, err := h.service.FeedBackCategoryList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	var dtoData []FeedBackCategoryDTO
	helpers.JsonMarshaller(result, &dtoData)
	return res.RespSuccessInfo(c, " List Retrieved successfully", dtoData, p)
}

func (h *handler) CreateFeedbackEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("rating"),
	}
	eda.Produce(eda.CREATE_FEEDBACK, request_payload)
	return res.RespSuccess(c, "Rating Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateFeedback(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(model.Feedback)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.SaveFeedBack(edaMetaData, createPayload)

	var responseMessage eda.ConsumerResponse
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_FEEDBACK_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_FEEDBACK_ACK, responseMessage)
}

func (h *handler) UpdateFeedbackEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	feedbackId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         feedbackId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("rating"),
	}
	eda.Produce(eda.UPDATE_FEEDBACK, request_payload)
	return res.RespSuccess(c, "Rating Update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateFeedback(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(model.Feedback)
	helpers.JsonMarshaller(data, updatePayLoad)
	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdateFeedback(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_FEEDBACK_ACK, responsemessage)
		return
	}
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_FEEDBACK_ACK, responsemessage)
}

func (h *handler) GetFeedback(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	FeedbackId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         FeedbackId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetFeedback(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Feedback Details Retrieved successfully", result)
}

func (h *handler) ListFeedBack(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	result, err := h.service.GetFeedbackList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Feedback list Retrieved Successfully", result, p)
}

func (h *handler) DeleteFeedback(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	FeedbackId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         FeedbackId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
		//"deleted_by": metaData.TokenUserId,
	}

	err = h.service.DeleteFeedback(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", FeedbackId)
}

func (h *handler) CreateReportEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)

	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("rating"),
	}
	eda.Produce(eda.CREATE_REPORT, request_payload)
	return res.RespSuccess(c, "Report Creation Inprogress", edaMetaData.RequestId)
}

func (h *handler) CreateReport(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)

	data := request["data"].(map[string]interface{})

	createPayload := new(model.Report)
	helpers.JsonMarshaller(data, createPayload)

	err := h.service.CreateReport(edaMetaData, createPayload)
	createPayload.CreatedByID = &edaMetaData.TokenUserId
	var responseMessage eda.ConsumerResponse
	responseMessage.MetaData = edaMetaData
	if err != nil {
		responseMessage.ErrorMessage = err
		// eda.Produce(eda.CREATE_REPORT_ACK, responseMessage)
		return
	}
	responseMessage.Response = map[string]interface{}{
		"created_id": createPayload.ID,
	}
	// eda.Produce(eda.CREATE_REPORT_ACK, responseMessage)

}

func (h *handler) UpdateReportEvent(c echo.Context) (err error) {
	edaMetaData := c.Get("MetaData").(core.MetaData)
	reportId := c.Param("id")
	edaMetaData.Query = map[string]interface{}{
		"id":         reportId,
		"company_id": edaMetaData.CompanyId,
	}
	request_payload := map[string]interface{}{
		"meta_data": edaMetaData,
		"data":      c.Get("rating"),
	}
	eda.Produce(eda.UPDATE_REPORT, request_payload)
	return res.RespSuccess(c, "Report Update inprogress", edaMetaData.RequestId)
}

func (h *handler) UpdateReport(request map[string]interface{}) {
	var edaMetaData core.MetaData
	helpers.JsonMarshaller(request["meta_data"], &edaMetaData)
	data := request["data"].(map[string]interface{})
	updatePayLoad := new(model.Report)
	helpers.JsonMarshaller(data, updatePayLoad)
	updatePayLoad.UpdatedByID = &edaMetaData.TokenUserId
	err := h.service.UpdateReport(edaMetaData, updatePayLoad)
	responsemessage := new(eda.ConsumerResponse)
	responsemessage.MetaData = edaMetaData
	if err != nil {
		responsemessage.ErrorMessage = err
		// eda.Produce(eda.UPDATE_REPORT_ACK, responsemessage)
		return
	}
	responsemessage.Response = map[string]interface{}{
		"updated_id": edaMetaData.Query["id"],
	}
	// eda.Produce(eda.UPDATE_REPORT_ACK, responsemessage)
}

func (h *handler) GetReport(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	reportId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         reportId,
		"company_id": metaData.CompanyId,
	}
	result, err := h.service.GetReport(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Report Details Retrieved successfully", result)
}

func (h *handler) ListReport(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)

	result, err := h.service.GetReportList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "Report list Retrieved Successfully", result, p)
}

func (h *handler) DeleteReport(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	reportId := c.Param("id")
	metaData.Query = map[string]interface{}{
		"id":         reportId,
		"user_id":    metaData.TokenUserId,
		"company_id": metaData.CompanyId,
		//"deleted_by": metaData.TokenUserId,
	}

	err = h.service.DeleteReport(metaData)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "Record deleted successfully", reportId)
}

func (h *handler) ProductsRatingSync(c echo.Context) (err error) {

	var schedulerRanTime time.Time

	p := new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = "[[\"name\",\"=\",\"ProductRatingScheduler\"],[\"state\",\"=\",true]]"
	schedulerJob, _ := h.schedulerservice.GetAllSchedulerJob(p)
	if len(schedulerJob) > 0 {
		schedulerJobId := fmt.Sprintf("%v", schedulerJob[0].ID)
		fmt.Println("-------------------------------schedulerJobId----------------------------", schedulerJobId)
		var log_data []scheduler.SchedulerLogResponseDTO
		p := new(pagination.Paginatevalue)
		p.Filters = "[[\"scheduler_job_id\",\"=\"," + schedulerJobId + "],[\"state\",\"=\",\"success\"],[\"notes\",\"@>\",\"{\\\"meta\\\":{\\\"success\\\":true}}\"]]"
		p.Sort = "[[\"id\",\"desc\"]]"
		log_data, _ = h.schedulerservice.ListSchedulerLogs(p)
		if len(log_data) > 0 {
			schedulerRanTime = log_data[0].StartTime
		}
	}

	fmt.Println("-------------------------------Last ProductRatingScheduler Ran Time----------------------------", schedulerRanTime)

	var metaData core.MetaData

	p = new(pagination.Paginatevalue)
	p.Filters = "[[\"category_name\",\"=\",\"ITEM\"]]"
	ratingCategoryInterface, err := h.service.RatingCategoriesList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	ratingCategory := ratingCategoryInterface.([]model.RatingCategory)
	if len(ratingCategory) == 0 {
		return res.RespErr(c, errors.New("rating category not found"))
	}

	p = new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = fmt.Sprintf("[[\"rating_category_id\",\"=\",\"%v\"],[\"updated_date\",\">\",\"%v\"]]", ratingCategory[0].ID, schedulerRanTime.Format("2006-01-02 15:04:05-07:00"))

	ratingsInterface, err := h.service.GetRatingList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	ratings := ratingsInterface.([]model.Rating)

	lenData := len(ratings)

	fmt.Println("ratings --------------> ", lenData)

	if lenData == 0 {
		return res.RespSuccess(c, "success", nil)
	}

	productRatings := map[uint]map[string]int64{}

	for _, rating := range ratings {
		if productRatings[rating.RelatedId] != nil {
			productRatings[rating.RelatedId]["rating"] += rating.RatingValue
			productRatings[rating.RelatedId]["count"] += 1
			continue
		}
		productRatings[rating.RelatedId] = map[string]int64{
			"rating": rating.RatingValue,
			"count":  1,
		}
	}

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	for product_id, ratings := range productRatings {

		query := map[string]interface{}{
			"id": product_id,
		}

		productInterface, err := h.productservice.GetVariantView(query, token_id, access_template_id)
		if err != nil {
			return res.RespErr(c, err)
		}

		product := productInterface.(products.VariantResponseDTO)

		var updatePayload products.CreateProductVariantDTO

		fmt.Println("\n\nProductId", product_id)
		fmt.Println("current", *product.ID, product.RatingAverage, product.RatingCount)
		fmt.Println("new", ratings["rating"], ratings["count"])

		updatePayload.RatingAverage = ((product.RatingAverage * float64(product.RatingCount)) + float64(ratings["rating"])) / (float64(product.RatingCount) + float64(ratings["count"]))
		updatePayload.RatingCount = product.RatingCount + uint(ratings["count"])

		fmt.Println("payload", updatePayload.RatingAverage, updatePayload.RatingCount)

		err = h.productservice.UpdateVariant(updatePayload, query, token_id, access_template_id)
		if err != nil {
			return res.RespErr(c, err)
		}
	}

	return res.RespSuccess(c, "success", nil)
}

func (h *handler) CompanyRatingSync(c echo.Context) (err error) {

	var schedulerRanTime time.Time

	p := new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = "[[\"name\",\"=\",\"CompanyRatingScheduler\"],[\"state\",\"=\",true]]"
	schedulerJob, _ := h.schedulerservice.GetAllSchedulerJob(p)
	if len(schedulerJob) > 0 {
		schedulerJobId := fmt.Sprintf("%v", schedulerJob[0].ID)
		fmt.Println("-------------------------------schedulerJobId----------------------------", schedulerJobId)
		var log_data []scheduler.SchedulerLogResponseDTO
		p := new(pagination.Paginatevalue)
		p.Filters = "[[\"scheduler_job_id\",\"=\"," + schedulerJobId + "],[\"state\",\"=\",\"success\"],[\"notes\",\"@>\",\"{\\\"meta\\\":{\\\"success\\\":true}}\"]]"
		p.Sort = "[[\"id\",\"desc\"]]"
		log_data, _ = h.schedulerservice.ListSchedulerLogs(p)
		if len(log_data) > 0 {
			schedulerRanTime = log_data[0].StartTime
		}
	}

	fmt.Println("-------------------------------Last ProductRatingScheduler Ran Time----------------------------", schedulerRanTime)

	var metaData core.MetaData

	p = new(pagination.Paginatevalue)
	p.Filters = "[[\"category_name\",\"=\",\"PROVIDER\"]]"
	ratingCategoryInterface, err := h.service.RatingCategoriesList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	ratingCategory := ratingCategoryInterface.([]model.RatingCategory)
	if len(ratingCategory) == 0 {
		return res.RespErr(c, errors.New("rating category not found"))
	}

	p = new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = fmt.Sprintf("[[\"rating_category_id\",\"=\",\"%v\"],[\"updated_date\",\">\",\"%v\"]]", ratingCategory[0].ID, schedulerRanTime.Format("2006-01-02 15:04:05-07:00"))

	ratingsInterface, err := h.service.GetRatingList(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	ratings := ratingsInterface.([]model.Rating)

	lenData := len(ratings)

	fmt.Println("ratings --------------> ", lenData)

	if lenData == 0 {
		return res.RespSuccess(c, "success", nil)
	}

	companyRatings := map[uint]map[string]int64{}

	for _, rating := range ratings {
		if companyRatings[rating.RelatedId] != nil {
			companyRatings[rating.RelatedId]["rating"] += rating.RatingValue
			companyRatings[rating.RelatedId]["count"] += 1
			continue
		}
		companyRatings[rating.RelatedId] = map[string]int64{
			"rating": rating.RatingValue,
			"count":  1,
		}
	}

	for company_id, ratings := range companyRatings {

		query := map[string]interface{}{
			"id": company_id,
		}

		company, err := h.coreService.GetCompanyPreferences(query)
		if err != nil {
			return res.RespErr(c, err)
		}

		var updatePayload cores.Company

		updatePayload.RatingAverage = ((company.RatingAverage * float64(company.RatingCount)) + float64(ratings["rating"])) / (float64(company.RatingCount) + float64(ratings["count"]))
		updatePayload.RatingCount = company.RatingCount + uint(ratings["count"])

		err = h.coreService.UpdateCompany(query, updatePayload, "", "")
		if err != nil {
			return res.RespErr(c, err)
		}
	}

	return res.RespSuccess(c, "success", nil)
}

func (h *handler) ListFeedBackForm(c echo.Context) (err error) {
	metaData := c.Get("MetaData").(core.MetaData)
	metaData.ModuleAccessAction = "LIST"
	rating_value := c.QueryParam("rating_value")
	rating_category_name := c.QueryParam("rating_category_name")
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	helpers.AddMandatoryFilters(p, "company_id", "=", metaData.CompanyId)
	query := map[string]interface{}{
		"rating_category_name": rating_category_name,
	}
	metaData.Query = query
	data_interface, _ := h.service.GetFeedbackCategory(metaData)
	feedback_category_id := data_interface.(model.FeedbackCategory).ID

	query = map[string]interface{}{
		"rating_value":         rating_value,
		"feedback_category_id": feedback_category_id,
	}
	metaData.Query = query
	result, err := h.service.ListFeedBackForm(metaData, p)
	if err != nil {
		return res.RespErr(c, err)
	}
	var dtoData FeedBackFormDto
	helpers.JsonMarshaller(result, &dtoData)
	return res.RespSuccessInfo(c, "Feedback form list Retrieved Successfully", dtoData, p)
}
