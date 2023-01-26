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
	GetCatalogue(category int, marketplace int, variant_id int) (interface{}, error)

	//catalogue template
	GetCatalogueTemplate(category int, marketplace int) (omnichannel.CatalogueTemplate, error)
	UpsertCatalogueTemplate(data *omnichannel.CatalogueTemplate) error

	//catalogue template data
	GetCatalogueData(category int, marketplace int, variant_id int) (omnichannel.CatalogueTemplateData, error)
	UpsertCatalogueTemplateData(data *omnichannel.CatalogueTemplateData) error
}

type service struct {
	catalogueRepository omnichannel_repo.Catalogue
}

func NewService() *service {
	catalogeRepository := omnichannel_repo.NewCatalogue()
	return &service{catalogeRepository}
}

func (s *service) UpsertCatalogueTemplate(data *omnichannel.CatalogueTemplate) error {
	res, _ := s.catalogueRepository.GetCatalogue(int(data.CategoryId), int(data.ChannelId))
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
	query := map[string]interface{}{
		"id":   res.ID,
		"data": resDto,
	}
	resp := s.catalogueRepository.UpdateCatalogue(query)
	return resp
}

func (s *service) GetCatalogueTemplate(category int, marketplace int) (omnichannel.CatalogueTemplate, error) {
	res, err := s.catalogueRepository.GetCatalogue(category, marketplace)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *service) GetCatalogueData(category int, marketplace int, variant_id int) (omnichannel.CatalogueTemplateData, error) {
	res, err := s.catalogueRepository.GetCatalogueData(category, marketplace, variant_id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *service) GetCatalogue(category int, marketplace int, variant_id int) (interface{}, error) {
	var response map[string]interface{}
	template_data, er := s.GetCatalogueData(category, marketplace, variant_id)
	if er != nil {
		template, err := s.GetCatalogueTemplate(category, marketplace)
		byte_template, _ := json.Marshal(template)
		json.Unmarshal(byte_template, &response)
		return response, err
	}
	byte_template, _ := json.Marshal(template_data)
	json.Unmarshal(byte_template, &response)
	return response, er
}

func (s *service) UpsertCatalogueTemplateData(data *omnichannel.CatalogueTemplateData) error {
	res, _ := s.catalogueRepository.GetCatalogueData(int(data.CategoryId), int(data.ChannelId), int(data.VariantId))
	if res.CategoryId == 0 && res.ChannelId == 0 {
		err := s.catalogueRepository.CreateCatalogueData(data)
		if err != nil {
			return err
		}
		return nil
	}
	var resDto CatalogueDTO
	mdata, _ := json.Marshal(data)
	json.Unmarshal(mdata, &resDto)
	query := map[string]interface{}{
		"id":   res.ID,
		"data": resDto,
	}
	resp := s.catalogueRepository.UpdateCatalogueData(query)
	return resp
}
