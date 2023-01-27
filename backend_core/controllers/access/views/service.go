package views

import (
	// "context"

	"fmt"

	model_core "fermion/backend_core/internal/model/core"
	"fermion/backend_core/internal/model/pagination"
	access_repo "fermion/backend_core/internal/repository/access"
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

type Services interface {
	RegisterViews(data *model_core.ViewsLevelAccessItems) error
	UpdateViewDetails(id uint, data model_core.ViewsLevelAccessItems) error
	FindAllView(p *pagination.Paginatevalue) ([]model_core.ViewsLevelAccessItems, error)
	GetViewDetail(id uint) (model_core.ViewsLevelAccessItems, error)
	DeleteViewDetail(id uint, user_id uint) error
}
type service struct {
	viewRepository access_repo.View
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	che := access_repo.NewView()
	newServiceObj = &service{che}
	return newServiceObj
}

func (s *service) RegisterViews(data *model_core.ViewsLevelAccessItems) error {

	err := s.viewRepository.CreateView(data)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UpdateViewDetails(id uint, data model_core.ViewsLevelAccessItems) error {
	_, er := s.viewRepository.FindByidView(id)

	if er != nil {
		return er
	}

	count, err := s.viewRepository.UpdateView(id, data)
	if err != nil {
		return err
	} else if count == 0 {
		return fmt.Errorf("no data gets updated")
	}

	return nil
}

func (s *service) FindAllView(p *pagination.Paginatevalue) ([]model_core.ViewsLevelAccessItems, error) {
	result, err := s.viewRepository.FindAllView(p)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetViewDetail(id uint) (model_core.ViewsLevelAccessItems, error) {

	cout, err := s.viewRepository.FindByidView(id)
	if err != nil {
		return cout, err
	}

	return cout, nil
}

func (s *service) DeleteViewDetail(id uint, user_id uint) error {

	_, er := s.viewRepository.FindByidView(id)
	if er != nil {
		return er
	}

	err := s.viewRepository.DeleteView(id, user_id)
	if err != nil {
		return err
	}
	return nil
}
