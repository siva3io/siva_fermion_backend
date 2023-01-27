package cache

import (
	"encoding/json"
	"fmt"

	"fermion/backend_core/db"

	"github.com/patrickmn/go-cache"
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

func GetCacheVariable(key string) interface{} {
	value, _ := db.InMemoryManager().Get(key)
	return value
}

func SetCacheVariable(key string, value interface{}) {
	db.InMemoryManager().Set(key, value, cache.NoExpiration)
}

func SetCacheKey(key string, value interface{}) {
	RedisClient := db.RedisManager()
	bytes, _ := json.Marshal(value)
	result := RedisClient.Set(key, string(bytes), 0)
	fmt.Println("Cache Result", result.Val())
}

func GetCacheKey(key string) interface{} {
	RedisClient := db.RedisManager()
	key_status := RedisClient.Exists(key).Val()
	if key_status > 0 {
		var cache_result interface{}
		result := RedisClient.Get(key).Val()
		json.Unmarshal([]byte(result), &cache_result)
		return cache_result
	}
	return nil
}
