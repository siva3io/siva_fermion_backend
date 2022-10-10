package ipaas_core

import (
	"errors"

	ipaas_repo "fermion/backend_core/ipaas_core/repository"
)

type Service interface {
	GetFeature(collection, id string) (interface{}, error)
	CreateFeature(collection string, data map[string]interface{}) (interface{}, error)
	UpdateFeature(collection, id string, data map[string]interface{}) (interface{}, error)
	DeleteFeature(collection, id string) (interface{}, error)
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
