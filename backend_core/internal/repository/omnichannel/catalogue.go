package omnichannel

import (
	"encoding/json"

	"fermion/backend_core/db"
	"fermion/backend_core/internal/model/omnichannel"

	"gorm.io/gorm"
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
type Catalogue interface {
	//Catalogue template
	CreateCatalogue(data *omnichannel.CatalogueTemplate) error
	GetCatalogue(category int, marketplace int) (omnichannel.CatalogueTemplate, error)
	UpdateCatalogue(query map[string]interface{}) error
	//Cataloge data
	CreateCatalogueData(data *omnichannel.CatalogueTemplateData) error
	GetCatalogueData(category int, marketplace int, variant_id int) (omnichannel.CatalogueTemplateData, error)
	UpdateCatalogueData(query map[string]interface{}) error
}

type catalogue struct {
	db *gorm.DB
}

func NewCatalogue() *catalogue {
	db := db.DbManager()
	return &catalogue{db}

}

func (r *catalogue) CreateCatalogue(data *omnichannel.CatalogueTemplate) error {
	err := r.db.Model(&omnichannel.CatalogueTemplate{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *catalogue) UpdateCatalogue(query map[string]interface{}) error {
	dto_data := query["data"]
	delete(query, "data")
	var data omnichannel.CatalogueTemplate
	byte_data, _ := json.Marshal(dto_data)
	json.Unmarshal(byte_data, &data)
	err := r.db.Model(&omnichannel.CatalogueTemplate{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *catalogue) GetCatalogue(category int, marketplace int) (omnichannel.CatalogueTemplate, error) {
	var data omnichannel.CatalogueTemplate
	res := r.db.Model(&omnichannel.CatalogueTemplate{}).Where("category_id", category).Where("channel_id", marketplace).First(&data)
	if res.Error != nil {
		return data, res.Error
	}
	return data, nil
}

func (r *catalogue) GetCatalogueData(category int, marketplace int, variant_id int) (omnichannel.CatalogueTemplateData, error) {
	var data omnichannel.CatalogueTemplateData
	res := r.db.Model(&omnichannel.CatalogueTemplateData{}).Where("category_id", category).Where("channel_id", marketplace).Where("variant_id", variant_id).First(&data)
	if res.Error != nil {
		return data, res.Error
	}
	return data, nil
}

func (r *catalogue) CreateCatalogueData(data *omnichannel.CatalogueTemplateData) error {
	err := r.db.Model(&omnichannel.CatalogueTemplateData{}).Create(data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *catalogue) UpdateCatalogueData(query map[string]interface{}) error {
	dto_data := query["data"]
	delete(query, "data")
	var data omnichannel.CatalogueTemplateData
	byte_data, _ := json.Marshal(dto_data)
	json.Unmarshal(byte_data, &data)
	err := r.db.Model(&omnichannel.CatalogueTemplateData{}).Where(query).Updates(data).Error
	if err != nil {
		return err
	}
	return nil
}
