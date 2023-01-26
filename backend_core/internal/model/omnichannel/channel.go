package omnichannel

import model_core "fermion/backend_core/internal/model/core"

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
type Channel struct {
	model_core.Model
	Name             string                `json:"name"`
	Sequence         *uint                 `json:"sequence" gorm:""`
	ExternalId       *uint                 `json:"external_id"`
	RelatedChannelId *uint                 `json:"related_channel_id"`
	ChannelTypeId    *uint                 `json:"channel_type_id"`
	Channel          model_core.Lookupcode `json:"lookup_data" gorm:"foreignKey:ChannelTypeId;references:ID"`
}
