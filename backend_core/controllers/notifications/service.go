package notifications

import (
	"encoding/json"

	core_service "fermion/backend_core/controllers/cores"

	"fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/internal/repository"
	notification_repo "fermion/backend_core/internal/repository/notification"
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

	//---------------------Product Template-------------------------------------------------------------------------------------
	CreateNotification(data map[string]interface{}) (interface{}, error)
	UpdateNotification(data map[string]interface{}, query map[string]interface{}) error
	DeleteNotification(query map[string]interface{}) error
	GetNotificationView(query map[string]interface{}) (interface{}, error)
	GetAllNotifications(p *pagination.Paginatevalue) (interface{}, error)
}

type service struct {
	notificationRepository notification_repo.Notifications
	userRepository         repository.User
	coreService            core_service.Service
}

func NewService() *service {
	notificationRepository := notification_repo.NewNotifications()
	userRepository := repository.NewUser()
	coreService := core_service.NewService()
	return &service{notificationRepository, userRepository, coreService}
}

//====================================TEMPLATE & VARIANT=====================================================================================================

// ---------------------Product Template-------------------------------------------------------------------------------------------------------------
func (s *service) CreateNotification(data map[string]interface{}) (interface{}, error) {

	var Notifications core.Notifications
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Notifications)
	if err != nil {
		return nil, err
	}
	err = s.notificationRepository.CreateNotification(&Notifications)
	if err != nil {
		return nil, err
	}
	return Notifications, nil
}
func (s *service) UpdateNotification(data map[string]interface{}, query map[string]interface{}) error {
	var Notifications core.Notifications

	dto, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dto, &Notifications)
	if err != nil {
		return err
	}

	err = s.notificationRepository.UpdateNotification(&Notifications, query)
	if err != nil {
		return err
	}
	return nil
}
func (s *service) DeleteNotification(query map[string]interface{}) error {

	err := s.notificationRepository.DeleteNotification(query)
	if err != nil {
		return err
	}

	return nil
}
func (s *service) GetNotificationView(query map[string]interface{}) (interface{}, error) {
	res, err := s.notificationRepository.FindNotification(query)
	if err != nil {
		return nil, err
	}

	var response TemplateReponseDTO
	dto, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &response)
	if err != nil {
		return nil, err
	}

	return response, err
}
func (s *service) GetAllNotifications(p *pagination.Paginatevalue) (interface{}, error) {
	result, err := s.notificationRepository.GetAllNotification(p)
	if err != nil {
		return nil, err
	}

	var response []interface{}
	mdata, _ := json.Marshal(result)
	json.Unmarshal(mdata, &response)

	return response, nil
}
