package catalogue

import (
	"encoding/json"

	"fermion/backend_core/internal/model/omnichannel"
	omnichannel_repo "fermion/backend_core/internal/repository/omnichannel"
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
type Service interface {

	//catalogue
	GetCatalogue(query map[string]interface{}) (interface{}, error)

	//catalogue template
	GetCatalogueTemplate(query map[string]interface{}) (omnichannel.CatalogueTemplate, error)
	UpsertCatalogueTemplate(data *omnichannel.CatalogueTemplate) error

	//catalogue template data
	GetCatalogueData(query map[string]interface{}) (omnichannel.CatalogueTemplateData, error)
	UpsertCatalogueTemplateData(data *omnichannel.CatalogueTemplateData) error
}

type service struct {
	catalogueRepository omnichannel_repo.Catalogue
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	catalogeRepository := omnichannel_repo.NewCatalogue()
	newServiceObj = &service{catalogeRepository}
	return newServiceObj
}

func (s *service) UpsertCatalogueTemplate(data *omnichannel.CatalogueTemplate) error {
	query := map[string]interface{}{
		"category_id": data.CategoryId,
		"channel_id":  data.ChannelId,
	}
	res, _ := s.catalogueRepository.GetCatalogue(query)
	if res.CategoryId == 0 && res.ChannelId == 0 {
		err := s.catalogueRepository.CreateCatalogue(data)
		if err != nil {
			return err
		}
		return nil
	}
	var resDto CatalogueDTO
	mdata, _ := json.Marshal(data)
	json.Unmarshal(mdata, &resDto)
	query1 := map[string]interface{}{
		"id":   res.ID,
		"data": resDto,
	}
	resp := s.catalogueRepository.UpdateCatalogue(query1)
	return resp
}

func (s *service) GetCatalogueTemplate(query map[string]interface{}) (omnichannel.CatalogueTemplate, error) {
	res, err := s.catalogueRepository.GetCatalogue(query)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *service) GetCatalogueData(query map[string]interface{}) (omnichannel.CatalogueTemplateData, error) {
	res, err := s.catalogueRepository.GetCatalogueData(query)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *service) GetCatalogue(query map[string]interface{}) (interface{}, error) {
	var response map[string]interface{}
	template_data, er := s.GetCatalogueData(query)
	delete(query, "variant_id")
	delete(query, "company_id")
	if er != nil {
		template, err := s.GetCatalogueTemplate(query)
		byte_template, _ := json.Marshal(template)
		json.Unmarshal(byte_template, &response)
		return response, err
	}
	byte_template, _ := json.Marshal(template_data)
	json.Unmarshal(byte_template, &response)
	return response, er
}

func (s *service) UpsertCatalogueTemplateData(data *omnichannel.CatalogueTemplateData) error {
	query := map[string]interface{}{
		"category_id": data.CategoryId,
		"channel_id":  data.ChannelId,
		"variant_id":  data.VariantId,
		"company_id":  data.CompanyId,
	}
	res, _ := s.catalogueRepository.GetCatalogueData(query)
	if res.CategoryId == 0 && res.ChannelId == 0 && res.CompanyId == 0 {
		err := s.catalogueRepository.CreateCatalogueData(data)
		if err != nil {
			return err
		}
		return nil
	}
	var resDto CatalogueDTO
	mdata, _ := json.Marshal(data)
	json.Unmarshal(mdata, &resDto)
	query1 := map[string]interface{}{
		"id":   res.ID,
		"data": resDto,
	}
	resp := s.catalogueRepository.UpdateCatalogueData(query1)
	return resp
}
