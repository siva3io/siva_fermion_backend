package ondc

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"fermion/backend_core/controllers/cores"
	"fermion/backend_core/controllers/eda"
	"fermion/backend_core/controllers/mdm/products"
	scheduler "fermion/backend_core/controllers/scheduler/app"
	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/cache"
	"fermion/backend_core/pkg/util/helpers"
	res "fermion/backend_core/pkg/util/response"

	"github.com/google/uuid"
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
	productservice   products.Service
	coreService      cores.Service
	schedulerservice scheduler.Service
}

var OndcHandler *handler //singleton object

// singleton function
func NewHandler() *handler {
	if OndcHandler != nil {
		return OndcHandler
	}
	service := NewService()
	product_service := products.NewService()
	core_service := cores.NewService()
	scheduler_service := scheduler.NewService()
	OndcHandler = &handler{service, product_service, core_service, scheduler_service}
	return OndcHandler
}

// ============================================ BAP Worflow APIs ======================================================== //

// func (h *handler) BAPRequestHandler(c echo.Context) (err error) {
// 	company_id := c.Get("CompanyId").(string)
// 	data := make(map[string]interface{})
// 	c.Bind(&data)
// 	var result_response interface{}
// 	action := c.Param("action")
// 	company, err := h.coreService.GetCompanyPreferences(map[string]interface{}{"id": company_id})
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	ondc_details := h.service.BuildOndcPayload(*company, action)

// 	if action == "search" {
// 		result_response = h.service.Search(ondc_details, data)

// 	} else if action == "select" {
// 		result_response = h.service.Select(ondc_details, data)

// 	} else if action == "init" {
// 		result_response = h.service.Init(ondc_details, data)

// 	} else if action == "confirm" {
// 		result_response = h.service.Confirm(ondc_details, data)

// 	} else if action == "status" {
// 		result_response = h.service.Status(ondc_details, data)

// 	} else if action == "track" {
// 		result_response = h.service.Track(ondc_details, data)

// 	} else if action == "update" {
// 		result_response = h.service.Update(ondc_details, data)

// 	} else if action == "cancel" {
// 		result_response = h.service.Cancel(ondc_details, data)

// 	}
// 	return res.RespSuccess(c, "Successfull Request", result_response)
// }

// ============================================ BPP Worflow API ======================================================== //
// var mapper map[string]interface{}

func (h *handler) BPPRequestHandler(c echo.Context) (err error) {
	var data BPPRequest
	c.Bind(&data)
	data.Context["bpp_id"] = "EUNIMART_BPP"
	data.Context["bpp_uri"] = "EUNIMART_BPP_URL"
	var result_response interface{}

	action := c.Param("action")
	if action == "search" {
		result_response = h.service.OnSearch(data)

	} else if action == "select" {
		result_response = h.service.OnSelect(data)

	} else if action == "init" {
		result_response = h.service.OnInit(data)

	} else if action == "confirm" {
		result_response = h.service.OnConfirm(data)

	} else if action == "status" {
		result_response = h.service.OnStatus(data)

	} else if action == "track" {
		result_response = h.service.OnTrack(data)

	} else if action == "update" {
		result_response = h.service.OnUpdate(data)

	} else if action == "cancel" {
		result_response = h.service.OnCancel(data)

	}
	response := make(map[string]interface{}, 0)
	response["data"] = result_response
	response["ondc_details"] = data.Context

	return res.RespSuccess(c, "successfull", response)

}

func (h *handler) ListBPP(c echo.Context) (err error) {
	response := []string{
		"Eunimart BPP",
		"Amazon",
		"Flipkart",
		"Myntra",
	}
	return res.RespSuccess(c, "successfull", response)
}

func (h *handler) BecknProductsSync(c echo.Context) (err error) {

	var schedulerRanTime time.Time

	p := new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = "[[\"name\",\"=\",\"ONDC-BECKN-product_sync\"],[\"state\",\"=\",true]]"
	schedulerJob, _ := h.schedulerservice.GetAllSchedulerJob(p)
	if len(schedulerJob) > 0 {
		schedulerJobId := fmt.Sprintf("%v", schedulerJob[0].ID)
		fmt.Println("-------------------------------schedulerJobId----------------------------", schedulerJobId)
		var log_data []scheduler.SchedulerLogResponseDTO
		p := new(pagination.Paginatevalue)
		p.Filters = "[[\"scheduler_job_id\",\"=\"," + schedulerJobId + "],[\"state\",\"=\",\"success\"]]"
		p.Sort = "[[\"id\",\"desc\"]]"
		log_data, _ = h.schedulerservice.ListSchedulerLogs(p)
		// helpers.PrettyPrint("logData", log_data)
		for _, log := range log_data {
			var log_response map[string]interface{}
			json.Unmarshal(log.Notes, &log_response)
			if meta := log_response["meta"]; meta != nil {
				if success := meta.(map[string]interface{})["success"]; success != nil {
					if success.(bool) {
						schedulerRanTime = log.StartTime
						break
					}
				}
			}
		}
	}

	fmt.Println("-------------------------------Last ONDC-BECKN-product_sync scheduler Ran Time----------------------------", schedulerRanTime)
	p = new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = fmt.Sprintf("[[\"channel\",\"=\",\"ONDC\"],[\"updated_date\",\">\",\"%v\"]]", schedulerRanTime.Format("2006-01-02 15:04:05-07:00"))

	token_id := c.Get("TokenUserID").(string)
	access_template_id := c.Get("AccessTemplateId").(string)

	variant_data_arr_interface, err := h.productservice.GetAllVariant(p, token_id, access_template_id, "LIST")
	if err != nil {
		return err
	}
	variant_data_arr := variant_data_arr_interface.([]products.VariantResponseDTO)

	lenData := len(variant_data_arr)

	fmt.Println("prodcuts --------------> ", lenData)

	if lenData == 0 {
		return res.RespSuccess(c, "success", nil)
	}

	mapper, _ := helpers.ReadFile("./backend_core/controllers/ondc/mapper_template.json")

	batchSize := 10

	responses := make([]interface{}, 0)

	for i := 0; i < lenData; i += batchSize {
		low, up := i, i+batchSize
		if up > lenData {
			up = low + (lenData % batchSize)
		}

		batch := variant_data_arr[low:up]

		response, err := helpers.MakeRequest(helpers.Request{
			Method: "POST",
			Scheme: "http",
			Host:   "localhost:3031",
			Path:   "ipaas/boson_convertor",
			Header: map[string]string{
				"Authorization": c.Request().Header.Get("Authorization"),
				"Content-Type":  "application/json",
			},
			Body: map[string]interface{}{
				"data": map[string]interface{}{
					"input_data":      batch,
					"mapper_template": mapper,
				},
			},
		})

		if err != nil {
			return res.RespErr(c, errors.New("error while mapping"))
		}

		data := response.(map[string]interface{})["data"]
		if data == nil {
			return res.RespErr(c, errors.New("error while mapping"))
		}
		if data.(map[string]interface{})["error_message"] != nil {
			return res.RespErr(c, errors.New("error while mapping"))
		}

		var metaData core.MetaData
		metaData.RequestId = fmt.Sprint(uuid.New())
		produce_payload := map[string]interface{}{
			"data":      data.(map[string]interface{})["mapped_response"],
			"meta_data": metaData,
		}

		eda.Produce(eda.BECKN_PRODUCTS_SYNC, produce_payload, false)

		var ackResponse interface{}
		start := time.Now()
		for {
			ackResponse = cache.GetCacheVariable(metaData.RequestId)
			if ackResponse != nil {
				break
			}
			if time.Since(start) > (60 * time.Second) {
				return res.RespErr(c, errors.New("beckn ack timeout"))
			}
		}

		data = ackResponse.(map[string]interface{})["data"]
		if data == nil {
			return res.RespErr(c, errors.New(fmt.Sprint(ackResponse)))
		}
		if data.(map[string]interface{})["success"] == nil || !data.(map[string]interface{})["success"].(bool) {
			return res.RespErr(c, errors.New(fmt.Sprint(ackResponse)))
		}

		responses = append(responses, ackResponse)

		fmt.Println(low, up)
	}

	return res.RespSuccess(c, "success", responses)
}

func (h *handler) BecknCompanySync(c echo.Context) (err error) {

	var schedulerRanTime time.Time

	p := new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = "[[\"name\",\"=\",\"ONDC-BECKN-company_sync\"],[\"state\",\"=\",true]]"
	schedulerJob, _ := h.schedulerservice.GetAllSchedulerJob(p)
	if len(schedulerJob) > 0 {
		schedulerJobId := fmt.Sprintf("%v", schedulerJob[0].ID)
		fmt.Println("-------------------------------schedulerJobId----------------------------", schedulerJobId)
		var log_data []scheduler.SchedulerLogResponseDTO
		p.Filters = "[[\"scheduler_job_id\",\"=\"," + schedulerJobId + "],[\"state\",\"=\",\"success\"]]"
		log_data, _ = h.schedulerservice.ListSchedulerLogs(p)
		// helpers.PrettyPrint("logData", log_data)
		for _, log := range log_data {
			var log_response map[string]interface{}
			json.Unmarshal(log.Notes, &log_response)
			if meta := log_response["meta"]; meta != nil {
				if success := meta.(map[string]interface{})["success"]; success != nil {
					if success.(bool) {
						schedulerRanTime = log.StartTime
						break
					}
				}
			}
		}
	}

	fmt.Println("-------------------------------Last ONDC-BECKN-Company_sync scheduler Ran Time----------------------------", schedulerRanTime)
	p = new(pagination.Paginatevalue)
	p.Per_page = -1
	p.Filters = fmt.Sprintf("[[\"updated_date\",\">\",\"%v\"]]", schedulerRanTime.Format("2006-01-02 15:04:05-07:00"))

	company_data_arr, err := h.coreService.ListCompanies(p)
	if err != nil {
		return err
	}

	lenData := len(company_data_arr)

	fmt.Println("company --------------> ", lenData)

	if lenData == 0 {
		return res.RespSuccess(c, "success", nil)
	}

	temp := make([]map[string]interface{}, lenData)

	for index, company := range company_data_arr {
		temp[index] = map[string]interface{}{
			"id":              company.Name,
			"company_details": company,
		}
	}

	responses := make([]interface{}, 0)

	batchSize := 10

	for i := 0; i < lenData; i += batchSize {

		low, up := i, i+batchSize
		if up > lenData {
			up = low + (lenData % batchSize)
		}

		batch := temp[low:up]

		var metaData core.MetaData
		metaData.RequestId = fmt.Sprint(uuid.New())
		produce_payload := map[string]interface{}{
			"data":      batch,
			"meta_data": metaData,
		}

		eda.Produce(eda.BECKN_PROVIDER_SYNC, produce_payload, false)

		var ackResponse interface{}
		start := time.Now()
		for {
			ackResponse = cache.GetCacheVariable(metaData.RequestId)
			if ackResponse != nil {
				break
			}
			if time.Since(start) > (60 * time.Second) {
				return res.RespErr(c, errors.New("beckn ack timeout"))
			}
		}

		data := ackResponse.(map[string]interface{})["data"]
		if data == nil {
			return res.RespErr(c, errors.New(fmt.Sprint(ackResponse)))
		}
		if data.(map[string]interface{})["success"] == nil || !data.(map[string]interface{})["success"].(bool) {
			return res.RespErr(c, errors.New(fmt.Sprint(ackResponse)))
		}

		responses = append(responses, ackResponse)

		fmt.Println(low, up)
	}

	return res.RespSuccess(c, "success", responses)
}
