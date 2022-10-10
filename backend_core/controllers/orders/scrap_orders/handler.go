package scrap_orders

import (
	"fmt"
	"strconv"

	orders_base "fermion/backend_core/controllers/orders/base"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service      Service
	base_service orders_base.ServiceBase
}

func NewHandler() *handler {
	service := NewService()
	base_service := orders_base.NewServiceBase()
	return &handler{service, base_service}
}

// CreateScrapOrder godoc
// @Summary      do a CreateScrapOrder
// @Description  Create a ScrapOrder
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param CreateScrapOrderRequest body  ScrapOrders true "create a ScrapOrder"
// @Success      200  {object}  ScrapOrderCreateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/create [post]
func (h *handler) CreateScrapOrder(c echo.Context) (err error) {
	data := c.Get("scrap_orders").(*ScrapOrders)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.CreatedByID = helpers.ConvertStringToUint(token_id)
	seq_scrap := helpers.GenerateSequence("scrap", token_id, "scrap_orders")
	seq_ref := helpers.GenerateSequence("Ref", token_id, "scrap_orders")

	if data.AutoCreateScrapNumber {
		data.Scrap_order_number = seq_scrap
	}
	if data.AutoGenerateReferenceNumber {
		data.Reference_id = seq_ref
	}

	id, err := h.service.CreateScrapOrder(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrders order added successfully", map[string]interface{}{"created_id": id})
}

// UpdateScrap godoc
// @Summary      do a UpdateScrap
// @Description  Update a ScrapOrder
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id                      path  int               true  "id"
// @Param UpdateScrapOrderRequest body  ScrapOrders true "update a ScrapOrder"
// @Success      200  {object}  ScrapOrderUpdateResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/{id}/update [post]
func (h *handler) UpdateScrapOrder(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	data := c.Get("scrap_orders").(*ScrapOrders)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	data.UpdatedByID = helpers.ConvertStringToUint(token_id)
	err = h.service.UpdateScrapOrder(uint(id), data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrders details updated succesfully", map[string]interface{}{"updated_id": id})
}

// FindAllScrapOrdersDropDown godoc
// @Summary      Find All ScrapOrders and filter it by search query
// @Description  Find All ScrapOrders and filter it by search query
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters 	query    string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  ScrapOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/dropdown [get]
func (h *handler) FindAllScrapOrdersDropDown(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AllScrapOrders(p, token_id, access_template_id, "DROPDOWN_LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "dropdown data Retrieved successdully", result, p)
}

// FindAllScrapOrders godoc
// @Summary      Find All ScrapOrders and filter it by search query
// @Description  Find All ScrapOrders and filter it by search query
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param        filters 	query    string false "filters"
// @Param   	per_page 	  query int false "per_page"
// @Param  		page_no       query int false "page_no"
// @Param   	sort		  query string false "sort"
// @Success      200  {object}  ScrapOrderGetAllResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders [get]
func (h *handler) FindAllScrapOrders(c echo.Context) (err error) {
	p := new(pagination.Paginatevalue)
	c.Bind(p)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.AllScrapOrders(p, token_id, access_template_id, "LIST")
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccessInfo(c, "data Retrieved successdully", result, p)
}

// ScrapOrdersByid godoc
// @Summary      Find a ScrapOrder By id
// @Description  Find a ScrapOrder By id
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id     path  int  true  "id"
// @Success      200  {object}  ScrapOrderGetResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/{id} [get]
func (h *handler) ScrapOrdersByid(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	result, err := h.service.FindScrapOrder(uint(id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "scrap details", result)
}

// DeleteScrap godoc
// @Summary      Delete a ScrapOrder
// @Description  Delete a ScrapOrder
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param id     path  int  true  "id"
// @Success      200  {object}  ScrapOrderDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse
// @Router       /api/v1/scrap_orders/{id}/delete [delete]
func (h *handler) DeleteScrapOrder(c echo.Context) (err error) {
	ID := c.Param("id")
	id, _ := strconv.Atoi(ID)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	err = h.service.DeleteScrapOrder(uint(id), uint(user_id), token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrders order deleted successfully", map[string]int{"deleted_id": id})
}

// DeleteScrapOrders godoc
// @Summary      Delete a ScrapOrderLine and provide the product_id as queryParam
// @Description  Delete a ScrapOrderLine and provide the product_id as queryParam
// @Tags         ScrapOrders
// @Accept       json
// @Produce      json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"

// @Param        line_id             path      int    true  "line_id"
// @Param        product_id     query     string true "product_id"
// @Success      200  {object}  ScrapOrderLinesDeleteResponse
// @Failure      400  {object}  res.ErrorResponse
// @Failure      404  {object}  res.ErrorResponse
// @Failure      500  {object}  res.ErrorResponse

// @Router     api/v1/scrap_orders/{line_id}/delete_product [delete]
func (h *handler) DeleteScrapOrderLines(c echo.Context) (err error) {
	product_id := c.QueryParam("product_id")
	prod_id, _ := strconv.Atoi(product_id)
	Id := c.Param("id")
	id, _ := strconv.Atoi(Id)
	query := map[string]interface{}{
		"scrap_id":   uint(id),
		"product_id": uint(prod_id),
	}
	fmt.Println(query)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	err = h.service.DeleteScrapOrderLines(query, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "order line deleted successfully", query)
}

// BulkCreateScrapOrder godoc
// @Summary 	do a BulkCreateScrapOrder
// @Description do a BulkCreateScrapOrder
// @Tags 		ScrapOrders
// @Accept  	json
// @Produce  	json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Param		BulkScrapOrderCreateRequest body ScrapOrders true "Create a Bulk ScrapOrder"
// @Success 	200 {object} BulkScrapOrderCreateResponse
// @Failure 	400 {object} res.ErrorResponse
// @Failure 	404 {object} res.ErrorResponse
// @Failure 	500 {object} res.ErrorResponse
// @Router 		/api/v1/scrap_orders/bulkcreate [post]
func (h *handler) BulkCreateScrapOrder(c echo.Context) (err error) {
	data := new([]ScrapOrders)
	c.Bind(data)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	t := helpers.ConvertStringToUint(token_id)
	for i := 0; i < len(*data); i++ {
		(*data)[i].CreatedByID = t
		fmt.Println((*data)[i].CreatedByID)
	}
	err = h.service.BulkCreateScrapOrder(data, token_id, access_template_id)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrders order added successfully", map[string]interface{}{"created_data": data})
}

// SendEmailScrapOrders godoc
// @Summary Email ScrapOrders
// @Description Email ScrapOrders
// @Tags ScrapOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/send_email [post]
func (h *handler) SendEmailScrapOrders(c echo.Context) error {
	return res.RespSuccess(c, "Email sent successfully", nil)
}

// DownloadScrapOrdersPDF godoc
// @Summary Download ScrapOrders
// @Description Download ScrapOrders
// @Tags ScrapOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/download_pdf [post]
func (h *handler) DownloadScrapOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "File downloaded successfully", resp_data)
}

// GenerateScrapOrdersPDF godoc
// @Summary Generate ScrapOrders PDF
// @Description Generate ScrapOrders PDF
// @Tags ScrapOrders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/generate_pdf [post]
func (h *handler) GenerateScrapOrdersPDF(c echo.Context) error {
	resp_data := map[string]interface{}{
		"file_url": "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	return res.RespSuccess(c, "PDF Generated", resp_data)
}

// FavouriteScrapOrders godoc
// @Summary FavouriteScrapOrders
// @Description FavouriteScrapOrders
// @Tags ScrapOrders
// @Accept  json
// @Produce  json
// @param id path string true "ScrapOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/favourite [post]
func (h *handler) FavouriteScrapOrder(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	_, er := h.service.FindScrapOrder(uint(ID), token_id, access_template_id)
	if er != nil {
		return res.RespSuccess(c, "Specified record not found", er)
	}
	err = h.base_service.FavouriteScrapOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrders is Marked as Favourite", map[string]string{"id": id})
}

// UnFavouriteScrapOrders godoc
// @Summary UnFavouriteScrapOrders
// @Description UnFavouriteScrapOrders
// @Tags ScrapOrders
// @Accept  json
// @Produce  json
// @param id path string true "ScrapOrders ID"
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/unfavourite [post]
func (h *handler) UnFavouriteScrapOrder(c echo.Context) (err error) {
	id := c.Param("id")
	ID, _ := strconv.Atoi(id)
	token_id := c.Get("TokenUserID").(string)
	user_id, _ := strconv.Atoi(token_id)
	query := map[string]interface{}{
		"ID":      ID,
		"user_id": user_id,
	}
	err = h.base_service.UnFavouriteScrapOrders(query)
	if err != nil {
		return res.RespErr(c, err)
	}
	return res.RespSuccess(c, "ScrapOrders is Unmarked as Favourite", map[string]string{"id": id})
}

// GetScrapOrderTab godoc
// @Summary GetScrapOrderTab
// @Summary GetScrapOrderTab
// @Description GetScrapOrderTab
// @Tags ScrapOrders
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @param id path string true "id"
// @param tab path string true "tab"
// @Success 200 {object} res.SuccessResponse
// @Failure 400 {object} res.ErrorResponse
// @Failure 404 {object} res.ErrorResponse
// @Failure 500 {object} res.ErrorResponse
// @Router /api/v1/scrap_orders/{id}/filter_module/{tab} [get]
func (h *handler) GetScrapOrderTab(c echo.Context) (err error) {

	page := new(pagination.Paginatevalue)
	err = c.Bind(page)
	if err != nil {
		return res.RespErr(c, err)
	}

	id := c.Param("id")
	tab := c.Param("tab")

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	data, err := h.service.GetScrapOrderTab(id, tab, page, access_template_id, token_id)
	if err != nil {
		return res.RespErr(c, err)
	}

	return res.RespSuccessInfo(c, "success", data, page)
}
