package core

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
type ViewsLevelAccessItems struct {
	Model
	TemplateLineItemId uint       `json:"template_line_item_id"`
	ViewsLineItemsId   uint       `json:"views_line_item_id"`
	ModuleTypeId       *uint      `json:"module_type_id" gorm:"type:int"`
	ModuleType         Lookupcode `json:"module_type" gorm:"foreignkey:ModuleTypeId; references:ID"`
	ParentViewID       *uint      `json:"parent_view_id" gorm:"column:parent_view"`
	ParentView         *ViewsLevelAccessItems
	SectionID          string `json:"section_id"`
}
