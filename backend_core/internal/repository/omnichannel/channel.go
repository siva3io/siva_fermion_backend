package omnichannel

import (
	"fermion/backend_core/db"
	channel "fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
	"fermion/backend_core/pkg/util/helpers"

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
type ChannelCrud interface {
	CreateChannel(data *channel.Channel) (interface{}, error)
	FindOneChannel(query map[string]interface{}) (channel.Channel, error)
	UpdateChannel(id uint, data channel.Channel) (channel.Channel, error)
	DeleteChannel(query map[string]interface{}, user_id uint) error
	FindAllChannels(p *pagination.Paginatevalue) ([]channel.Channel, error)
}
type database struct {
	db *gorm.DB
}

var channelRepository *database //singleton object

// singleton function
func NewChannel() *database {
	if channelRepository != nil {
		return channelRepository
	}
	db := db.DbManager()
	channelRepository = &database{db}
	return channelRepository
}

func (r *database) CreateChannel(data *channel.Channel) (interface{}, error) {
	err := r.db.Model(&channel.Channel{}).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil

}
func (r *database) FindOneChannel(query map[string]interface{}) (channel.Channel, error) {
	var data channel.Channel
	err := r.db.Where(query).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}

func (r *database) UpdateChannel(id uint, data channel.Channel) (channel.Channel, error) {
	err := r.db.Model(&channel.Channel{}).Where("id", id).Updates(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
func (r *database) DeleteChannel(query map[string]interface{}, user_id uint) error {
	var data channel.Channel
	err := r.db.Model(&channel.Channel{}).Where(query).Update("deleted_by", user_id).Error
	if err != nil {
		return err
	}
	err = r.db.Model(&channel.Channel{}).Where((query)).Delete(&data).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *database) FindAllChannels(p *pagination.Paginatevalue) ([]channel.Channel, error) {
	var data []channel.Channel
	err := r.db.Model(&channel.Channel{}).Order("sequence").Scopes(helpers.Paginate(&channel.Channel{}, p, r.db)).Find(&data).Error
	if err != nil {
		return data, err
	}
	return data, nil
}
