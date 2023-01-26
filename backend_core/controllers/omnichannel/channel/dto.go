package channel

import "gorm.io/datatypes"

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
type ChannelDTO struct {
	ID               uint           `json:"id"`
	Name             string         `json:"name"`
	CreatedByID      *uint          `json:"created_by"`
	UpdatedByID      *uint          `json:"updated_by"`
	DeletedByID      *uint          `json:"deleted_by"`
	Sequence         *uint          `json:"sequence"`
	ExternalId       *uint          `json:"external_id"`
	RelatedChannelId *uint          `json:"related_channel_id"`
	ChannelTypeId    *uint          `json:"channel_type_id"`
	ChannelDetails   datatypes.JSON `json:"channel_details" gorm:"type:json; default:'[]'"`
}
