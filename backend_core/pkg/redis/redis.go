package rds

import (
	"os"

	"github.com/go-redis/redis/v8"
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
var redisClient *redis.Client

func Init() {
	var (
		REDIS_HOST = os.Getenv("REDIS_HOST")
		REDIS_PORT = os.Getenv("REDIS_PORT")
		REDIS_PASS = os.Getenv("REDIS_PASS")
	)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: REDIS_PASS,
		DB:       0, // use default DB
	})
}

func RedisManager() *redis.Client {
	return redisClient
}
