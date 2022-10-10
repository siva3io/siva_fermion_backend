package pos

import (
	accounting_repo "fermion/backend_core/internal/repository/accounting"

	"fermion/backend_core/internal/model/accounting"
)

/*
Copyright (C) 2022 Eunimart Omnichannel Pvt Ltd. (www.eunimart.com)
All rights reserved.
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/
type Service interface {
	CreatePos(data *accounting.UserPosLink) error
	UpdatePos(query map[string]interface{}, data *accounting.UserPosLink) error
	GetAuthKeys(query map[string]interface{}) (interface{}, error)
}

type service struct {
	posRepository accounting_repo.Pos
}

func NewService() *service {
	return &service{accounting_repo.NewPosRepo()}
}

func (s *service) CreatePos(data *accounting.UserPosLink) error {

	err := s.posRepository.CreatePos(data)
	return err
}

func (s *service) UpdatePos(query map[string]interface{}, data *accounting.UserPosLink) error {

	_, er := s.posRepository.FindOne(query)
	if er != nil {
		return er
	}
	err := s.posRepository.UpdatePos(query, data)
	return err
}

func (s *service) GetAuthKeys(query map[string]interface{}) (interface{}, error) {
	result, err := s.posRepository.FindOne(query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
