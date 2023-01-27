<!--
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
-->
# euni_plt_go_gorm

Eunimart Platform build using Go Echo framework.

- PORT : 3031

## Build & Run

```bash
# If u wanna run in development
$ ENV=DEV go run main.go

#if u wanna run in development with local db cleard
$ ENV=DEV go run main.go -clearDB=true

#if u wanna run in development with seed master data
$ ENV=DEV go run main.go -seedMD=true

#if u wanna run in development with migrations
$ ENV=DEV go run main.go -migrateDB=true

#if u wanna run particular changes in db with migrations
Follow The below given steps
1. Change the .env.development file 
Example:
VERSION=1.0.0.0 to VERSION=1.0.0.1

2. In migrations folder need to add another file for version 
Example:
File Name -- 1  and index.json, migration_version_control.sql to
File Name -- 2  and add index.json, migration_version_control.sql

3. In migration_version_controls.sql file need to add version 
Example:
INSERT INTO public.migration_version_controls (version) VALUES ("1")

4. In migration_version_control.sql file need to add the change of sql query
Example:
If you want to add a column in db
ALTER TABLE table_name ADD field_name type_of_field;
ALTER TABLE public.sales_orders ADD new_filed text;

After completing all the above steps need to run 

$ ENV=DEV go  run main.go -migrateDBVersion=true



# If u wanna build
$ go build

#If u wanna run in docker kafka
$ docker-compose -f kafka-docker-compose/docker-compose.yml up -d
```

## Generate swagger documentation

install go swagger

```bash
go get github.com/swaggo/swag

cd $GOPATH/src/github.com/swaggo/swag

go install cmd/swag

```

to generate API documentation

```bash

export PATH=$(go env GOPATH)/bin:$PATH

swag init --parseDependency

```

to see the results, run app and access {{base_url}}/swagger/index.html

## [Manifest](manifest.md)

## [Ipaas](ipaas_core.md)

## How to Install Redis Cache 
In order to install Redis, you first need to update the APT repository cache of your Ubuntu. You can do that with the following command:

```bash

sudo apt update

```

Now, you can enter the command to install Redis. Enter the following string into the command line:

```bash

sudo apt install redis

```

In order to check if Redis is installed properly and working correctly, you can enter the command:

```bash

redis-cli --version

```

Once you complete the installation, you can check if the Redis is running. You can do this with the following command:

```bash

sudo systemctl status redis

```

In the output, locate Active: active (running).
If Redis hasn’t been started, you can launch it by entering the command:

```bash

sudo systemctl start redis-server

```
Note: After successful installation by default redis will run on 6379 port and default password is empty string or blank string 

# Author

Eunimart
