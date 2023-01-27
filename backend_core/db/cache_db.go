package db

import (
	"fmt"
	"os"
	"time"

	// "fermion/backend_core/pkg/util/helpers"
	"github.com/go-redis/redis"
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
var redisClient *redis.Client
var CacheClient *cache.Cache

func CacheInit() {

	// Redis connection
	var (
		REDIS_HOST = os.Getenv("REDIS_HOST")
		REDIS_PORT = os.Getenv("REDIS_PORT")
		// REDIS_DB   = os.Getenv("REDIS_DB")
		REDIS_PASS = os.Getenv("REDIS_PASS")
	)

	redisClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: REDIS_PASS,
		DB:       0, // use default DB
	})
	status := redisClient.Ping().Val()

	if status == "" {
		fmt.Println("Redis cache No Status returned")
		// panic("failed to connect cache database")
	}

	// In-memory cache
	CacheClient = cache.New(5*time.Minute, 10*time.Minute)
	if CacheClient == nil {
		// panic("failed to create in-memory cache")
	}

}

func RedisManager() *redis.Client {
	return redisClient
}

func InMemoryManager() *cache.Cache {
	return CacheClient
}
