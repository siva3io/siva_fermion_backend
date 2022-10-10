

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
