package helpers

import (
	"fmt"
	"math"
	"strings"

	Pagination "fermion/backend_core/internal/model/pagination"

	"gorm.io/gorm"
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
func Paginate(value interface{}, pagination *Pagination.Paginatevalue, db *gorm.DB) func(db *gorm.DB) *gorm.DB {

	fields, values, err := ApplyFilter(pagination.Filters)
	if err != nil {
		panic("filter error")
	}

	sortfields, err := ApplySort(pagination.Sort)
	if err != nil {
		fmt.Println(err)
		panic("Sort error")
	}

	pagination.Filters = ""
	pagination.Sort = ""

	db.Model(value).Where(strings.Join(fields, " AND "), values...).Count(&pagination.TotalRows)

	per_page := pagination.GetLimit()

	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(per_page)))

	return func(db *gorm.DB) *gorm.DB {
		return db.Model(value).Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(sortfields).Where(strings.Join(fields, " AND "), values...)
	}
}
