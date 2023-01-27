package ipaas_core

import (
	"errors"

	ipaas_repo "fermion/backend_core/ipaas_core/repository"
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
	GetFeature(collection, id string) (interface{}, error)
	CreateFeature(collection string, data map[string]interface{}) (interface{}, error)
	UpdateFeature(collection, id string, data map[string]interface{}) (interface{}, error)
	DeleteFeature(collection, id string) (interface{}, error)
	GetAppFeatures(query map[string]interface{}) (interface{}, error)
}

type service struct {
	ipaasFeaturesRepository ipaas_repo.IpaasFeatures
}

func NewService() *service {
	return &service{ipaas_repo.NewRepo()}
}

var Collections = map[string]string{
	"features": "ipaas_features",
	"mappers":  "ipaas_mappers",
	"apis":     "ipaas_apis",
}

func (s *service) GetFeature(collection, id string) (interface{}, error) {

	if Collections[collection] != "" {
		feature, err := s.ipaasFeaturesRepository.GetOne(Collections[collection], id)
		if err != nil {
			return nil, err
		}
		return feature, nil
	}

	return nil, errors.New("collection not found")
}

func (s *service) CreateFeature(collection string, data map[string]interface{}) (interface{}, error) {

	if Collections[collection] != "" {
		response, err := s.ipaasFeaturesRepository.CreateOne(Collections[collection], data)
		if err != nil {
			return nil, err
		}
		return response, nil

	}

	return nil, errors.New("collection not found")
}

func (s *service) UpdateFeature(collection, id string, data map[string]interface{}) (interface{}, error) {

	if Collections[collection] != "" {
		response, err := s.ipaasFeaturesRepository.UpdateOne(Collections[collection], id, data)
		if err != nil {
			return nil, err
		}
		return response, nil

	}

	return nil, errors.New("collection not found")
}

func (s *service) DeleteFeature(collection, id string) (interface{}, error) {

	if Collections[collection] != "" {
		data, err := s.ipaasFeaturesRepository.DeleteOne(Collections[collection], id)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	return nil, errors.New("collection not found")
}

func (s *service) GetAppFeatures(query map[string]interface{}) (interface{}, error) {

	collectionName := "ipaas_features"

	features, err := s.ipaasFeaturesRepository.GetMany(collectionName, query)
	if err != nil {
		return nil, err
	}

	return features, nil
}
