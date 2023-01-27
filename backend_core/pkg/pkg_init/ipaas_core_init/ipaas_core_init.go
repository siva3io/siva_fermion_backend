package ipaas_core_init

import (
	"fmt"
	"time"

	pagination "fermion/backend_core/internal/model/pagination"
	core_repo "fermion/backend_core/internal/repository"
	ipaas "fermion/backend_core/ipaas_core/app"
	cache "fermion/backend_core/pkg/util/cache"
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

func Init() {

	KafkaConsumerTimeout := 30 * time.Minute

	go ipaas.InitConsumers(KafkaConsumerTimeout)
	go CacheSync()
}

// ===========================Cache Sync=========================================================================================================================================================================
func CacheSync() {
	p := new(pagination.Paginatevalue)
	p.Per_page = -1
	installed_apps, _ := core_repo.NewCore().ListInstalledApps(p)
	applist := make([]string, len(installed_apps))
	for index, app := range installed_apps {
		cache.SetCacheKey(app.Code, app)
		applist[index] = app.Code
	}
	fmt.Println("Cache sync is successfull for the following installed applications")
	fmt.Println(applist)
}
