<!--/*
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
*/-->

## Swagger Docs Go
Swagger is a set of open-source tools built around the OpenAPI Specification that can help you design, build, document and consume REST APIs. The major Swagger tools include: Swagger Editor – browser-based editor where you can write OpenAPI definitions
## If you want to check your api's through swagger this docs.go file will help you.

## Step-1
install go swagger.

go get github.com/swaggo/swag

## Step-2
Add go path in your local machine.

cd $GOPATH/src/github.com/swaggo/swag

go install cmd/swag

## Step-3
To generate API documentation use below command.

swag init --parseDependency

## Step-4
To see the results, run app and access {{base_url}}/swagger/index.html



Eunimart
