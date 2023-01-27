package concurrency_management

import (
	"fmt"
	"time"

	"fermion/backend_core/db"
	pkg_helpers "fermion/backend_core/pkg/util/helpers"

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
type Service interface {
	CheckConcurrencyStatus(id uint, column_name string, model_name string) (ConcurrencyResponseDTO, uint, error)
	SetConcurrencyLock(id uint, StatusId uint, ConcurrencyLookupCodeId uint, column_name string, model_name string) error
	ReleaseConcurrencyLock(id uint, status_id uint, column_name string, model_name string) error
	ReleaseDeadLock(id uint, status_id uint, column_name string, model_name string)
}

const DeadLockOccurTime = 300 // 5 minutes

type service struct {
	db *gorm.DB
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	db := db.DbManager()
	newServiceObj = &service{db}
	return newServiceObj
}

func (s *service) CheckConcurrencyStatus(id uint, column_name string, model_name string) (ConcurrencyResponseDTO, uint, error) {
	var response ConcurrencyResponseDTO
	response.Id = id
	response.Type = model_name
	response.Block = false
	response.Message = "no locks"

	var StatusId uint
	query := fmt.Sprintf("SELECT %v FROM %v WHERE id = %v", column_name, model_name, id)
	err := s.db.Raw(query).Scan(&StatusId).Error
	if err != nil {
		return response, 0, err
	}

	ConcurrencyLookupCodeId, err := pkg_helpers.GetLookupcodeId("CONCURRENCY_TYPE", "LOCK")
	if err != nil {
		return response, 0, err
	}

	if ConcurrencyLookupCodeId == StatusId {
		response.Block = true
		response.Message = fmt.Sprintf("%s is in lock. please try again after some time.", model_name)
		pkg_helpers.PrettyPrint("concurrency_status", response)
		return response, StatusId, nil
	}

	if ConcurrencyLookupCodeId != StatusId {

		err := s.SetConcurrencyLock(id, StatusId, ConcurrencyLookupCodeId, column_name, model_name)
		if err != nil {
			return response, StatusId, err
		}
		pkg_helpers.PrettyPrint("concurrency_status", response)
		return response, StatusId, nil
	}

	return response, StatusId, nil
}
func (s *service) SetConcurrencyLock(id uint, StatusId uint, ConcurrencyLookupCodeId uint, column_name string, model_name string) error {

	LockStartTime := time.Now().Format("2006-01-02 15:04:05.000000-07:00")

	query := fmt.Sprintf("UPDATE %v SET %v = %v, updated_date = '%v' WHERE id = %v", model_name, column_name, ConcurrencyLookupCodeId, LockStartTime, id)
	err := s.db.Exec(query).Error
	if err != nil {
		return err
	}
	fmt.Println("-------------->>>> Concurrency Set Lock <<<<-----------------")
	//------------go routine----------------------------------
	go s.ReleaseDeadLock(id, StatusId, column_name, model_name)
	return nil
}
func (s *service) ReleaseConcurrencyLock(id uint, status_id uint, column_name string, model_name string) error {

	LockReleaseTime := time.Now().Format("2006-01-02 15:04:05.000000-07:00")

	ConcurrencyLookupCodeId, err := pkg_helpers.GetLookupcodeId("CONCURRENCY_TYPE", "LOCK")
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %v SET %v = %v, updated_date = '%v' WHERE id = %v AND %v = %v", model_name, column_name, status_id, LockReleaseTime, id, column_name, ConcurrencyLookupCodeId)
	err = s.db.Exec(query).Error
	if err != nil {
		return err
	}
	fmt.Println("-------------->>>> Concurrency Release Lock <<<<-----------------")
	return nil
}
func (s *service) ReleaseDeadLock(id uint, status_id uint, column_name string, model_name string) {

	time.Sleep(time.Second * DeadLockOccurTime)
	var LastUpdatedTime time.Time

	query := fmt.Sprintf("SELECT %v FROM %v WHERE id = %v", "updated_date", model_name, id)
	s.db.Raw(query).Scan(&LastUpdatedTime)

	LastUpdatedTime = LastUpdatedTime.Add(time.Second * DeadLockOccurTime)

	//-----------if the lock exceeds the maximum time, then forcefully release the lock------------------
	if LastUpdatedTime.Before(time.Now()) {
		fmt.Println("-------------->>>>  Release DeadLock <<<<-----------------")
		s.ReleaseConcurrencyLock(id, status_id, column_name, model_name)
	}
}
