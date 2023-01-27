package channel

import (
	"encoding/json"
	"fmt"

	channel "fermion/backend_core/internal/model/omnichannel"
	"fermion/backend_core/internal/model/pagination"
	repo "fermion/backend_core/internal/repository/omnichannel"
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
	CreateChannel(data ChannelDTO) (interface{}, error)
	CreateMarketplace(data ChannelDTO) (interface{}, error)
	GetChannel(data ChannelDTO, lookupcode string) (ChannelDTO, error)
	UpdateChannel(channel_data ChannelDTO, channel_values ChannelDTO, lookupcode string) (ChannelDTO, error)
	CreateWebstore(data ChannelDTO) (interface{}, error)
	DeleteChannel(data ChannelDTO, user_id uint, lookupcode string) (ChannelDTO, error)
	FindAllChannels(query map[string]interface{}, p *pagination.Paginatevalue) ([]ChannelDTO, error)
	GetChannelService(query map[string]interface{}) (ChannelDTO, error)
}
type service struct {
	repo             repo.ChannelCrud
	marketplace_repo repo.OmnichannelBase
	webstore_repo    repo.Webstore
	// marketplace marketplace.Marketplace
}

var newServiceObj *service //singleton object

// singleton function
func NewService() *service {
	if newServiceObj != nil {
		return newServiceObj
	}
	ChannelRepository := repo.NewChannel()
	marketplaceRepository := repo.NewmarketPlaceBase()
	webstoreRepository := repo.NewWebstore()
	newServiceObj = &service{ChannelRepository, marketplaceRepository, webstoreRepository}
	return newServiceObj
}

func (s *service) CreateChannel(data ChannelDTO) (interface{}, error) {
	var Channel channel.Channel
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Channel)
	if err != nil {
		return nil, err
	}
	fmt.Println(Channel)
	response, err := s.repo.CreateChannel(&Channel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) CreateMarketplace(data ChannelDTO) (interface{}, error) {

	var Channel channel.Channel
	var Marketplace channel.Marketplace
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Channel)
	if err != nil {
		return nil, err
	}
	dto, err = json.Marshal(data.ChannelDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Marketplace)
	if err != nil {
		return nil, err
	}
	err = s.marketplace_repo.CreateMarketplace(&Marketplace)
	if err != nil {
		return nil, err
	}
	Channel.RelatedChannelId = &Marketplace.ID
	response, err := s.repo.CreateChannel(&Channel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) CreateWebstore(data ChannelDTO) (interface{}, error) {

	var Channel channel.Channel
	var WebStore channel.Webstore
	dto, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &Channel)
	if err != nil {
		return nil, err
	}
	dto, err = json.Marshal(data.ChannelDetails)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dto, &WebStore)
	if err != nil {
		return nil, err
	}
	err = s.webstore_repo.CreateWebstoreData(&WebStore)
	if err != nil {
		return nil, err
	}
	Channel.RelatedChannelId = &WebStore.ID
	response, err := s.repo.CreateChannel(&Channel)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *service) GetChannelService(query map[string]interface{}) (ChannelDTO, error) {
	data, err := s.repo.FindOneChannel(query)
	var channel_dto ChannelDTO
	if err != nil {
		return channel_dto, err
	}
	dto, err := json.Marshal(data)
	if err != nil {
		return channel_dto, err
	}
	err = json.Unmarshal(dto, &channel_dto)
	if err != nil {
		return channel_dto, err
	}
	return channel_dto, err

}
func (s *service) GetChannel(data ChannelDTO, lookupcode string) (ChannelDTO, error) {

	var err error
	var result interface{}
	if lookupcode == "MARKETPLACE" {
		result, err = s.marketplace_repo.FindMarketplace(*data.RelatedChannelId)
	}
	if lookupcode == "WEBSTORE" {
		result, err = s.webstore_repo.FindWebstore(*data.RelatedChannelId)
	}

	dto, err := json.Marshal(result)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(dto, &data.ChannelDetails)
	if err != nil {
		return data, err
	}
	if err != nil {
		return data, err
	}
	return data, nil
}

func (s *service) UpdateChannel(channel_data ChannelDTO, channel_values ChannelDTO, lookupcode string) (ChannelDTO, error) {
	var marketPlace channel.Marketplace
	var webStore channel.Webstore
	var channel channel.Channel
	jsonData, err := json.Marshal(channel_data.ChannelDetails)
	if err != nil {
		return channel_data, err
	}
	if lookupcode == "MARKETPLACE" {
		err = json.Unmarshal(jsonData, &marketPlace)
		if err != nil {
			return channel_data, err
		}
		err := s.marketplace_repo.UpdateMarketplace(*channel_values.RelatedChannelId, marketPlace)
		if err != nil {
			return channel_data, err
		}
	}
	if lookupcode == "WEBSTORE" {
		err = json.Unmarshal(jsonData, &webStore)
		if err != nil {
			return channel_data, err
		}
		err := s.webstore_repo.UpdateWebstoreData(*channel_values.RelatedChannelId, webStore)

		if err != nil {
			return channel_data, err
		}

	}
	jsonData, err = json.Marshal(channel_data)
	if err != nil {
		return channel_data, err
	}
	err = json.Unmarshal(jsonData, &channel)
	if err != nil {
		return channel_data, err
	}
	_, err = s.repo.UpdateChannel(channel_values.ID, channel)
	if err != nil {
		return channel_data, err
	}
	return channel_data, nil

}

func (s *service) DeleteChannel(data ChannelDTO, user_id uint, lookupcode string) (ChannelDTO, error) {

	var channel_dto ChannelDTO
	var err error
	ID := data.ID
	var query = make(map[string]interface{}, 0)
	query["id"] = ID
	if lookupcode == "MARKETPLACE" {
		err = s.marketplace_repo.DeleteMarketplace(*data.RelatedChannelId, user_id)
	}
	if lookupcode == "WEBSTORE" {
		err = s.webstore_repo.DeleteWebstoreData(*data.RelatedChannelId, user_id)
	}

	if err != nil {
		return channel_dto, err
	}
	err = s.repo.DeleteChannel(query, user_id)
	if err != nil {
		return channel_dto, err
	}
	return channel_dto, nil
}

func (s *service) FindAllChannels(query map[string]interface{}, p *pagination.Paginatevalue) ([]ChannelDTO, error) {
	var data []channel.Channel
	var result interface{}
	data, err := s.repo.FindAllChannels(p)
	if err != nil {
		return nil, err
	}
	var dto = make([]ChannelDTO, len(data))
	for i, val := range data {
		external_id := val.ExternalId
		external_value := *external_id
		json_data, err := json.Marshal(&val)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(json_data, &dto[i])
		if external_value != uint(0) {
			if *dto[i].ChannelTypeId == query["marketplace"] {
				result, err = s.marketplace_repo.FindMarketplace(*dto[i].RelatedChannelId)
			}
			if *dto[i].ChannelTypeId == query["webstore"] {
				result, err = s.webstore_repo.FindWebstore(*dto[i].RelatedChannelId)
			}
			result_data, err := json.Marshal(&result)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(result_data, &dto[i].ChannelDetails)
			if err != nil {
				return dto, err
			}
		}
	}
	return dto, nil
}
